package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/httputil/proxy"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/utils"
)

var (
	log = logging.LoggerForModule()
)

const (
	connectTimeout = 5 * time.Second
)

// email notifier plugin
type email struct {
	config     *storage.Email
	smtpServer smtpServer

	notifier *storage.Notifier
}

type plainAuthUnencrypted struct {
	identity, username, password string
	host                         string
}

func unencryptedPlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuthUnencrypted{
		identity: identity,
		username: username,
		password: password,
		host:     host,
	}
}

func (a *plainAuthUnencrypted) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// This is modified from smtp.plainAuth.Start()
	// to remove the check that passwords can only be sent unencrypted
	// to localhost.
	// As long as we claim to support unencrypted SMTP we need to
	// override this check, since the user is explicitly telling us
	// to do the bad idea.
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuthUnencrypted) Next(fromServer []byte, more bool) ([]byte, error) {
	// This is copied from smtp.plainAuth.Next().
	// See Start() for reasons why we have copied this type.
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}

type loginAuth struct {
	username, password string
}

func loginAuthMethod(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		serverStr := string(fromServer)
		switch serverStr {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown value request %q from server", serverStr)
		}
	}
	return nil, nil
}

type smtpServer struct {
	host string
	port int
}

func (s *smtpServer) endpoint() string {
	return fmt.Sprintf("%v:%v", s.host, s.port)
}

func validate(email *storage.Email) error {
	errorList := errorhelpers.NewErrorList("Email validation")
	if email.GetServer() == "" {
		errorList.AddString("SMTP Server must be specified")
	}
	if email.GetSender() == "" {
		errorList.AddString("Sender must be specified")
	}
	if email.GetUsername() == "" {
		errorList.AddString("Username must be specified")
	}
	if email.GetPassword() == "" {
		errorList.AddString("Password must be specified")
	}
	if !email.GetDisableTLS() && email.GetStartTLSAuthMethod() != storage.Email_DISABLED {
		errorList.AddString("TLS must be disabled to use a StartTLS Auth Method")
	}
	return errorList.ToError()
}

func newEmail(notifier *storage.Notifier) (*email, error) {
	emailConfig, ok := notifier.GetConfig().(*storage.Notifier_Email)
	if !ok {
		return nil, errors.New("Email configuration required")
	}
	conf := emailConfig.Email
	if err := validate(conf); err != nil {
		return nil, err
	}

	port := 465 // default TLS SMTP Port
	server := conf.GetServer()
	host := conf.GetServer()
	idx := strings.Index(server, ":")
	if idx != -1 && idx != len(server)-1 {
		parsedPort, err := strconv.Atoi(server[idx+1:])
		if err != nil || parsedPort < 0 || parsedPort > 65535 {
			return nil, fmt.Errorf("Port number cannot be '%v' and must be valid port between 0-65535", server[idx+1:])
		}
		port = parsedPort
		host = server[:idx]
	}
	return &email{
		config: conf,
		smtpServer: smtpServer{
			host: host,
			port: port,
		},
		notifier: notifier,
	}, nil
}

type message struct {
	To      string
	From    string
	Subject string
	Body    string
}

func (m message) Bytes() []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", m.From)
	fmt.Fprintf(&buf, "To: %s\r\n", m.To)
	fmt.Fprintf(&buf, "Subject: %s\r\n", m.Subject)
	fmt.Fprint(&buf, "Content-Type: text/plain; charset=utf-8\r\n\r\n")
	fmt.Fprintf(&buf, "%s\r\n", m.Body)
	return buf.Bytes()
}

func (e *email) plainTextAlert(alert *storage.Alert) (string, error) {
	funcMap := template.FuncMap{
		"header": func(s string) string {
			return fmt.Sprintf("\r\n%s\r\n", s)
		},
		"subheader": func(s string) string {
			return fmt.Sprintf("\r\n\t%s\r\n", s)
		},
		"line": func(s string) string {
			return fmt.Sprintf("%s\r\n", s)
		},
		"list": func(s string) string {
			return fmt.Sprintf("\t - %s\r\n", s)
		},
		"nestedList": func(s string) string {
			return fmt.Sprintf("\t\t - %s\r\n", s)
		},
		"codeBlock": func(s string) string {
			return fmt.Sprintf("\n %s \n", s)
		},
		"section": func(s string) string {
			return fmt.Sprintf("\r\n\t\t%s\r\n", s)
		},
		"group": func(s string) string {
			return fmt.Sprintf("\r\n\t\t\t - %s", s)
		},
	}
	alertLink := notifiers.AlertLink(e.notifier.UiEndpoint, alert.GetId())
	return notifiers.FormatPolicy(alert, alertLink, funcMap)
}

func (*email) Close(context.Context) error {
	return nil
}

// AlertNotify takes in an alert and generates the email
func (e *email) AlertNotify(ctx context.Context, alert *storage.Alert) error {
	subject := fmt.Sprintf("Deployment %v (%v) violates '%v' Policy", alert.GetDeployment().GetName(),
		alert.GetDeployment().GetId(), alert.GetPolicy().GetName())
	body, err := e.plainTextAlert(alert)
	if err != nil {
		return err
	}

	recipient := notifiers.GetLabelValue(alert, e.notifier.GetLabelKey(), e.notifier.GetLabelDefault())
	return e.sendEmail(ctx, recipient, subject, body)
}

// YamlNotify takes in a yaml file and generates the email message
func (e *email) NetworkPolicyYAMLNotify(ctx context.Context, yaml string, clusterName string) error {
	subject := fmt.Sprintf("New network policy YAML for cluster '%s' needs to be applied", clusterName)

	body, err := notifiers.FormatNetworkPolicyYAML(yaml, clusterName, template.FuncMap{
		"codeBlock": func(s string) string {
			return s
		},
	})
	if err != nil {
		return err
	}
	return e.sendEmail(ctx, e.notifier.GetLabelDefault(), subject, body)
}

// Test sends a test notification
func (e *email) Test(ctx context.Context) error {
	subject := "StackRox Test Email"
	body := fmt.Sprintf("%v\r\n", "This is a test email created to test integration with StackRox.")
	err := e.sendEmail(ctx, e.notifier.GetLabelDefault(), subject, body)
	return err
}

func (e *email) sendEmail(ctx context.Context, recipient, subject, body string) error {
	var from string
	if e.config.GetFrom() != "" {
		from = fmt.Sprintf("%s <%s>", e.config.GetFrom(), e.config.GetSender())
	} else {
		from = e.config.GetSender()
	}

	msg := message{
		To:      recipient,
		From:    from,
		Subject: subject,
		Body:    body,
	}

	conn, auth, err := e.connection(ctx)
	if err != nil {
		log.Errorf("Connection failed: %v", err)
		return err
	}

	client, err := e.createClient(conn)
	if err != nil {
		log.Errorf("SMTP client creation failed: %v", err)
		return err
	}
	defer func() {
		if err := client.Quit(); err != nil {
			log.Errorf("Failed to quit client cleanly: %v", err)
		}
	}()

	if e.config.GetStartTLSAuthMethod() != storage.Email_DISABLED {
		if err = client.StartTLS(e.tlsConfig()); err != nil {
			log.Errorf("SMTP STARTTLS failed: %v", err)
			return err
		}
	}

	if err = client.Auth(auth); err != nil {
		log.Errorf("SMTP authentication failed: %v", err)
		return err
	}

	if err = client.Mail(e.config.GetSender()); err != nil {
		log.Errorf("SMTP MAIL command failed: %v", err)
		return err
	}
	if err = client.Rcpt(recipient); err != nil {
		log.Errorf("SMTP RCPT command failed: %v", err)
		return err
	}

	w, err := client.Data()
	if err != nil {
		log.Errorf("SMTP DATA command failed: %v", err)
		return err
	}
	defer utils.IgnoreError(w.Close)

	_, err = w.Write(msg.Bytes())
	if err != nil {
		log.Errorf("SMTP message writing failed: %v", err)
		return err
	}

	return nil
}

// createClient creates an SMTP client but bails out in cases where
// smtp.NewClient would otherwise hang.
// The known case (ROX-366) is when dialing a TLS server with a non-TLS dialer;
// in this case the first dial will succeed, but then the net/textproto reader
// will hang for a few minutes.
func (e *email) createClient(conn net.Conn) (c *smtp.Client, err error) {
	var timedOut concurrency.Flag
	// If the timer expires before we return and thereby stop it,
	// we'll close the connection and thereby cause the Client creation
	// to abort immediately instead of waiting for minutes for an EOF.
	//
	// There's a possibility that we have _just_ succeeded returning from
	// NewClient when this timer fires; in this case the subsequent use of
	// the client will fail with an error about using a closed connection.
	// This particular failure mode seems sufficiently unlikely.
	// Importantly, a net.Conn can have multiple clients safely call methods
	// on it at the same time, including Close().
	t := time.AfterFunc(connectTimeout, func() {
		timedOut.Toggle()
		defer utils.IgnoreError(conn.Close)
	})
	defer func() {
		t.Stop()
		if timedOut.Get() {
			err = errors.New("timeout: possibly speaking unencrypted to a server running TLS")
		}
	}()

	return smtp.NewClient(conn, e.smtpServer.host)
}

func (e *email) connection(ctx context.Context) (conn net.Conn, auth smtp.Auth, err error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	if e.config.GetDisableTLS() {
		if e.config.GetStartTLSAuthMethod() != storage.Email_DISABLED {
			return e.startTLSConn(ctx)
		}
		return e.unencryptedConn(ctx)
	}
	return e.tlsConn(ctx)
}

func (e *email) tlsConn(dialCtx context.Context) (conn net.Conn, auth smtp.Auth, err error) {
	// With a connection that starts with TLS, we can simply use the standard
	// library to authenticate.
	auth = smtp.PlainAuth("", e.config.GetUsername(), e.config.GetPassword(), e.smtpServer.host)
	conn, err = proxy.AwareDialContextTLS(dialCtx, e.smtpServer.endpoint(), e.tlsConfig())
	return
}

func (e *email) unencryptedConn(dialCtx context.Context) (conn net.Conn, auth smtp.Auth, err error) {
	// With a completely unencrypted connection, we must override the
	// standard library's SMTP authenticator, since it blocks attempts
	// to send credentials over any non-TLS connection that isn't localhost.
	auth = unencryptedPlainAuth("", e.config.GetUsername(), e.config.GetPassword(), e.smtpServer.host)
	conn, err = proxy.AwareDialContext(dialCtx, e.smtpServer.endpoint())
	return
}

func (e *email) startTLSConn(dialCtx context.Context) (conn net.Conn, auth smtp.Auth, err error) {
	// With STARTTLS, we will first connect unencrypted and later
	// "upgrade" the connection to use TLS by the time we authenticate.
	// Hence, we can use the stdlib authenticator, which treats
	// STARTTLS as TLS for purposes of deciding whether it's safe to
	// transmit a password.
	switch e.notifier.GetEmail().GetStartTLSAuthMethod() {
	case storage.Email_PLAIN:
		auth = smtp.PlainAuth("", e.config.GetUsername(), e.config.GetPassword(), e.smtpServer.host)
	case storage.Email_LOGIN:
		auth = loginAuthMethod(e.config.GetUsername(), e.config.GetPassword())
	}
	conn, err = proxy.AwareDialContext(dialCtx, e.smtpServer.endpoint())
	return
}

func (e *email) tlsConfig() *tls.Config {
	return &tls.Config{
		ServerName: e.smtpServer.host,
	}
}

func (e *email) ProtoNotifier() *storage.Notifier {
	return e.notifier
}

func init() {
	notifiers.Add("email", func(notifier *storage.Notifier) (notifiers.Notifier, error) {
		e, err := newEmail(notifier)
		return e, err
	})
}
