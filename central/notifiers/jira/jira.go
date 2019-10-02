package jira

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	jiraLib "github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/urlfmt"
)

const (
	timeout = 5 * time.Second
)

var (
	log = logging.LoggerForModule()

	defaultPriorities = map[storage.Severity]string{
		storage.Severity_CRITICAL_SEVERITY: "P0-Highest",
		storage.Severity_HIGH_SEVERITY:     "P1-High",
		storage.Severity_MEDIUM_SEVERITY:   "P2-Medium",
		storage.Severity_LOW_SEVERITY:      "P3-Low",
		storage.Severity_UNSET_SEVERITY:    "P4-Lowest",
	}
)

// Jira notifier plugin
type jira struct {
	client *jiraLib.Client

	conf *storage.Jira

	notifier *storage.Notifier

	priorities map[storage.Severity]string
}

func (j *jira) getAlertDescription(alert *storage.Alert) (string, error) {
	funcMap := template.FuncMap{
		"header": func(s string) string {
			return fmt.Sprintf("\r\n h4. %v\r\n", s)
		},
		"subheader": func(s string) string {
			return fmt.Sprintf("\r\n h5. %v\r\n", s)
		},
		"line": func(s string) string {
			return fmt.Sprintf("%v\r\n", s)
		},
		"list": func(s string) string {
			return fmt.Sprintf("* %v\r\n", s)
		},
		"nestedList": func(s string) string {
			return fmt.Sprintf("** %v\r\n", s)
		},
	}
	alertLink := notifiers.AlertLink(j.notifier.UiEndpoint, alert.GetId())
	return notifiers.FormatPolicy(alert, alertLink, funcMap)
}

// AlertNotify takes in an alert and generates the notification
func (j *jira) AlertNotify(alert *storage.Alert) error {
	description, err := j.getAlertDescription(alert)
	if err != nil {
		return err
	}

	project := notifiers.GetLabelValue(alert, j.notifier.GetLabelKey(), j.notifier.GetLabelDefault())
	i := &jiraLib.Issue{
		Fields: &jiraLib.IssueFields{
			Summary: fmt.Sprintf("Deployment %v (%v) violates '%v' Policy", alert.Deployment.Name, alert.Deployment.Id, alert.Policy.Name),
			Type: jiraLib.IssueType{
				Name: j.conf.GetIssueType(),
			},
			Project: jiraLib.Project{
				Key: project,
			},
			Description: description,
			Priority: &jiraLib.Priority{
				Name: j.severityToPriority(alert.GetPolicy().GetSeverity()),
			},
		},
	}
	return j.createIssue(i)
}

func (j *jira) NetworkPolicyYAMLNotify(yaml string, clusterName string) error {
	funcMap := template.FuncMap{
		"codeBlock": func(s string) string {
			return fmt.Sprintf("{code:title=Network Policy YAML|theme=FadeToGrey|language=yaml}%s{code}", s)
		},
	}

	description, err := notifiers.FormatNetworkPolicyYAML(yaml, clusterName, funcMap)
	if err != nil {
		return err
	}

	project := j.notifier.GetLabelDefault()
	i := &jiraLib.Issue{
		Fields: &jiraLib.IssueFields{
			Summary: fmt.Sprintf("Network policy yaml to apply on cluster %s", clusterName),
			Type: jiraLib.IssueType{
				Name: j.conf.GetIssueType(),
			},
			Project: jiraLib.Project{
				Key: project,
			},
			Description: description,
			Priority: &jiraLib.Priority{
				Name: j.severityToPriority(storage.Severity_MEDIUM_SEVERITY),
			},
		},
	}
	return j.createIssue(i)
}

func validate(jira *storage.Jira) error {
	errorList := errorhelpers.NewErrorList("Jira validation")
	if jira.GetIssueType() == "" {
		errorList.AddString("Issue Type must be specified")
	}
	if jira.GetUrl() == "" {
		errorList.AddString("URL must be specified")
	}
	if jira.GetUsername() == "" {
		errorList.AddString("Username must be specified")
	}
	if jira.GetPassword() == "" {
		errorList.AddString("Password must be specified")
	}
	return errorList.ToError()
}

func newJira(notifier *storage.Notifier) (*jira, error) {
	conf := notifier.GetJira()
	if conf == nil {
		return nil, errors.New("Jira configuration required")
	}
	if err := validate(conf); err != nil {
		return nil, err
	}

	url, err := urlfmt.FormatURL(conf.GetUrl(), urlfmt.HTTPS, urlfmt.TrailingSlash)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Timeout: timeout,
	}
	bat := &jiraLib.BasicAuthTransport{
		Username:  conf.GetUsername(),
		Password:  conf.GetPassword(),
		Transport: httpClient.Transport,
	}
	client, err := jiraLib.NewClient(bat.Client(), url)
	if err != nil {
		return nil, err
	}

	prios, _, err := client.Priority.GetList()
	if err != nil {
		return nil, err
	}
	log.Infof("Retrieved priority list: %+v", prios)

	return &jira{
		client:     client,
		conf:       notifier.GetJira(),
		notifier:   notifier,
		priorities: mapPriorities(prios),
	}, nil
}

func (j *jira) ProtoNotifier() *storage.Notifier {
	return j.notifier
}

func (j *jira) createIssue(i *jiraLib.Issue) error {
	_, resp, err := j.client.Issue.Create(i)
	if err != nil && resp == nil {
		return errors.Errorf("Error creating issue. Response: %v", err)
	}

	if err != nil {
		bytes, readErr := ioutil.ReadAll(resp.Body)
		if readErr == nil {
			log.Errorf("Error creating issue. Response: %v %v", err, string(bytes))
			return errors.Errorf("Error creating issue: received HTTP status code %d. Check central logs for full error.", resp.StatusCode)
		}
	}
	return err
}

func (j *jira) Test() error {
	i := &jiraLib.Issue{
		Fields: &jiraLib.IssueFields{
			Description: "StackRox Test Issue",
			Type: jiraLib.IssueType{
				Name: j.conf.GetIssueType(),
			},
			Project: jiraLib.Project{
				Key: j.notifier.GetLabelDefault(),
			},
			Summary: "This is a test issue created to test integration with StackRox.",
			Priority: &jiraLib.Priority{
				Name: j.severityToPriority(storage.Severity_LOW_SEVERITY),
			},
		},
	}
	return j.createIssue(i)
}

func mapPriorities(prios []jiraLib.Priority) map[storage.Severity]string {
	output := make(map[storage.Severity]string)
	for k, name := range defaultPriorities {
		for _, p := range prios {
			if len(p.Name) < 3 {
				continue
			}
			if name[:3] == p.Name[:3] {
				name = p.Name
			}
		}
		output[k] = name
	}
	return output
}

func (j *jira) severityToPriority(sev storage.Severity) string {
	name, ok := j.priorities[sev]
	if ok {
		return name
	}
	return j.priorities[storage.Severity_UNSET_SEVERITY]
}

func init() {
	notifiers.Add("jira", func(notifier *storage.Notifier) (notifiers.Notifier, error) {
		j, err := newJira(notifier)
		return j, err
	})
}
