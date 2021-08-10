package centralclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	cTLS "github.com/google/certificate-transparency-go/tls"
	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/cryptoutils"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/mtls/verifier"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/x509utils"
)

var (
	log = logging.LoggerForModule()
)

const (
	requestTimeout          = 10 * time.Second
	tlsChallengeRoute       = "/v1/tls-challenge"
	metadataRoute           = "/v1/metadata"
	challengeTokenParamName = "challengeToken"
)

// Client is a client which provides functions to call rest routes in central
type Client struct {
	endpoint   *url.URL
	httpClient *http.Client
}

// NewClient creates a new client
func NewClient(endpoint string) (*Client, error) {
	if endpoint == "" {
		return nil, errors.New("creating Client with empty endpoint is not allowed")
	}

	parts := strings.SplitN(endpoint, "://", 2)
	switch parts[0] {
	case "wss":
		endpoint = fmt.Sprintf("https://%s", parts[1])
	case "https":
		break
	default:
		if len(parts) == 1 {
			endpoint = fmt.Sprintf("https://%s", endpoint)
			break
		}
		return nil, errors.Errorf("creating client unsupported scheme %s", parts[0])
	}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "parsing endpoint url")
	}

	// Load the client certificate. Note that while all endpoints accessed by the client do not require
	// authentication, it is possible that a user has required client certificate authentication for the
	// endpoint sensor is connecting to. Since a client certificate can be used without harm even if the
	// remote is not trusted, make it available here to be on the safe side.
	clientCert, err := mtls.LeafCertificateFromFile()
	if err != nil {
		return nil, errors.Wrap(err, "obtaining client certificate")
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{
			clientCert,
		},
	}
	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConf},
		Timeout:   requestTimeout,
	}

	return &Client{
		httpClient: httpClient,
		endpoint:   endpointURL,
	}, nil
}

// GetMetadata returns centrals metadata
func (c *Client) GetMetadata(ctx context.Context) (*v1.Metadata, error) {
	resp, err := c.doHTTPRequest(ctx, http.MethodGet, metadataRoute, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "receiving central metadata")
	}
	defer utils.IgnoreError(resp.Body.Close)

	var metadata v1.Metadata
	err = jsonpb.Unmarshal(resp.Body, &metadata)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing central %s response with status code %d", metadataRoute, resp.StatusCode)
	}

	return &metadata, nil
}

// GetTLSTrustedCerts returns all certificates which are trusted by central and its leaf certificates.
// Sensor validates the identity of central by verifying the given signature against centrals public key presented by its leaf cert.
func (c *Client) GetTLSTrustedCerts(ctx context.Context) ([]*x509.Certificate, error) {
	token, err := c.generateChallengeToken()
	if err != nil {
		return nil, errors.Wrap(err, "creating challenge token")
	}

	resp, err := c.doTLSChallengeRequest(ctx, &v1.TLSChallengeRequest{ChallengeToken: token})
	if err != nil {
		return nil, err
	}

	trustInfo, err := c.parseTLSChallengeResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err, "verifying tls challenge")
	}

	var certs []*x509.Certificate
	for _, ca := range trustInfo.GetAdditionalCas() {
		cert, err := x509.ParseCertificate(ca)
		if err != nil {
			return nil, errors.Wrap(err, "parsing additional CA")
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

func (c *Client) parseTLSChallengeResponse(challenge *v1.TLSChallengeResponse) (*v1.TrustInfo, error) {
	var trustInfo v1.TrustInfo
	err := proto.Unmarshal(challenge.GetTrustInfoSerialized(), &trustInfo)
	if err != nil {
		return nil, errors.Wrap(err, "parsing TrustInfo proto")
	}

	if len(trustInfo.GetCertChain()) == 0 {
		return nil, errors.New("reading centrals leaf certificate from response")
	}

	rootCAs, err := verifier.TrustedCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "reading CA cert")
	}

	x509CertChain, err := c.verifyCertificateChain(trustInfo.GetCertChain(), rootCAs)
	if err != nil {
		return nil, err
	}
	if len(x509CertChain) == 0 {
		return nil, errors.New("parsing central chain was empty, expected certificate chain")
	}

	centralLeafCert := x509CertChain[0]
	err = cTLS.VerifySignature(centralLeafCert.PublicKey, challenge.TrustInfoSerialized, cTLS.DigitallySigned{
		Signature: challenge.Signature,
		Algorithm: cTLS.SignatureAndHashAlgorithm{
			Hash:      cTLS.SHA256,
			Signature: cTLS.SignatureAlgorithmFromPubKey(centralLeafCert.PublicKey),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "verifying central trust info signature")
	}

	return &trustInfo, nil
}

func (c *Client) verifyCertificateChain(certChain [][]byte, rootCAs *x509.CertPool) ([]*x509.Certificate, error) {
	x509CertChain, err := x509utils.ParseCertificateChain(certChain)
	if err != nil {
		return nil, errors.Wrap(err, "parsing central cert chain")
	}

	err = x509utils.VerifyCertificateChain(x509CertChain, x509.VerifyOptions{
		Roots:   rootCAs,
		DNSName: mtls.CentralSubject.Hostname(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "verifying central cert chain")
	}

	return x509CertChain, nil
}

// doTLSChallengeRequest send the HTTP request to central and receives the trust info.
func (c *Client) doTLSChallengeRequest(ctx context.Context, req *v1.TLSChallengeRequest) (*v1.TLSChallengeResponse, error) {
	params := url.Values{challengeTokenParamName: []string{req.GetChallengeToken()}}

	resp, err := c.doHTTPRequest(ctx, http.MethodGet, tlsChallengeRoute, params, nil)
	if err != nil {
		return nil, errors.Wrap(err, "receiving centrals trust info")
	}
	defer utils.IgnoreError(resp.Body.Close)

	tlsChallengeResp := &v1.TLSChallengeResponse{}
	err = jsonpb.Unmarshal(resp.Body, tlsChallengeResp)
	if err != nil {
		return nil, errors.Wrap(err, "parsing central response")
	}
	return tlsChallengeResp, nil
}

func (c *Client) doHTTPRequest(ctx context.Context, method, route string, params url.Values, body io.Reader) (*http.Response, error) {
	u := *c.endpoint
	u.Path = route
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "creating tls-challenge request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "calling %s", u.String())
	}
	if !httputil.Is2xxStatusCode(resp.StatusCode) {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, errors.Wrapf(err, "reading response body with HTTP status code '%s'", resp.Status)
		}
		return nil, errors.Errorf("HTTP request %s%s with code '%s', body: %s", c.endpoint, tlsChallengeRoute, resp.Status, body)
	}
	return resp, nil
}

func (c *Client) generateChallengeToken() (string, error) {
	nonceGenerator := cryptoutils.NewNonceGenerator(centralsensor.ChallengeTokenLength, nil)
	challenge, err := nonceGenerator.Nonce()
	if err != nil {
		return "", err
	}

	return challenge, nil
}
