package renderer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/docker/distribution/reference"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/defaultimages"
	"github.com/stackrox/rox/pkg/grpc/authn/basic"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/zip"
)

var (
	log = logging.LoggerForModule()
)

// ExternalPersistence holds the data for a volume that is already created (e.g. docker volume, PV, etc)
type ExternalPersistence struct {
	Name         string `json:"name,omitempty"`
	StorageClass string `json:"storageClass,omitempty"`
	Size         uint32 `json:"size,omitempty"`
}

// HostPathPersistence describes the parameters for a bind mount
type HostPathPersistence struct {
	HostPath          string
	NodeSelectorKey   string
	NodeSelectorValue string
}

// WithNodeSelector is a helper function for the templater that returns if node selectors are used
func (h *HostPathPersistence) WithNodeSelector() bool {
	if h == nil {
		return false
	}
	return h.NodeSelectorKey != ""
}

// CommonConfig contains the common config between orchestrators that cannot be placed at the top level
// Image is an example as it can be parameterized per orchestrator with different defaults so it cannot be placed
// at the top level
type CommonConfig struct {
	MainImage        string
	ScannerImage     string
	ScannerV2Image   string
	ScannerV2DBImage string
	MonitoringImage  string
}

// MonitoringType is the enum for the place monitoring is hosted
type MonitoringType int

// Types of monitoring
const (
	None MonitoringType = iota
	OnPrem
	StackRoxHosted
)

// String returns the string form of the enum
func (m MonitoringType) String() string {
	switch m {
	case OnPrem:
		return "on-prem"
	case None:
		return "none"
	case StackRoxHosted:
		return "stackrox-hosted"
	}
	return "unknown"
}

// OnPrem is true if the monitoring is hosted on prem
func (m MonitoringType) OnPrem() bool {
	return m == OnPrem
}

// StackRoxHosted is true if the monitoring is hosted by StackRox
func (m MonitoringType) StackRoxHosted() bool {
	return m == StackRoxHosted
}

// None returns true if there is no monitoring solution
func (m MonitoringType) None() bool {
	return m == None
}

// PersistenceType describes the type of persistence
type PersistenceType string

// Types of persistence
var (
	PersistenceNone     = newPersistentType("none")
	PersistenceHostpath = newPersistentType("hostpath")
	PersistencePVC      = newPersistentType("pvc")
)

// StringToPersistentTypes is a map from the persistenttype string value to its object
var StringToPersistentTypes = make(map[string]PersistenceType)

func newPersistentType(t string) PersistenceType {
	pt := PersistenceType(t)
	StringToPersistentTypes[t] = pt
	return pt
}

// String returns the string form of the enum
func (m PersistenceType) String() string {
	return string(m)
}

// MonitoringConfig encapsulates the monitoring configuration
type MonitoringConfig struct {
	Endpoint         string
	Image            string
	Type             MonitoringType
	LoadBalancerType v1.LoadBalancerType
	Password         string
	PasswordAuto     bool

	PersistenceType PersistenceType
	External        *ExternalPersistence
	HostPath        *HostPathPersistence
}

// ScannerV2Config encapsulates the scanner v2 configuration.
type ScannerV2Config struct {
	Enable bool `json:"enable"`

	PersistenceType PersistenceType     `json:"persistenceType"`
	External        ExternalPersistence `json:"externalPersistence,omitempty"`
	HostPath        HostPathPersistence `json:"hostPathPersistence,omitempty"`
}

// K8sConfig contains k8s fields
type K8sConfig struct {
	CommonConfig
	ConfigType v1.DeploymentFormat

	// K8s Application name
	AppName string

	// k8s fields
	Registry string

	ScannerRegistry string
	// If the scanner registry is different from the central registry get a separate secret
	ScannerSecretName string

	ScannerV2Config ScannerV2Config

	// These variables are not prompted for by Cobra, but are set based on
	// provided inputs for use in templating.
	MainImageTag    string
	ScannerImageTag string

	DeploymentFormat v1.DeploymentFormat
	LoadBalancerType v1.LoadBalancerType

	// Command is either oc or kubectl depending on the value of cluster type
	Command string

	Monitoring MonitoringConfig

	OfflineMode bool
}

// Config configures the deployer for the central service.
type Config struct {
	ClusterType storage.ClusterType
	OutputDir   string

	K8sConfig *K8sConfig

	External *ExternalPersistence
	HostPath *HostPathPersistence

	Password     string
	PasswordAuto bool

	LicenseData []byte

	DefaultTLSCertPEM []byte
	DefaultTLSKeyPEM  []byte

	SecretsByteMap   map[string][]byte
	SecretsBase64Map map[string]string

	Environment map[string]string

	GCPMarketplace bool
	Version        string
}

func executeRawTemplate(raw string, c *Config) ([]byte, error) {
	t, err := template.New("temp").Parse(raw)
	if err != nil {
		return nil, err
	}
	return ExecuteTemplate(t, c)
}

// ExecuteTemplate renders a given template, injecting the given values.
func ExecuteTemplate(temp *template.Template, values interface{}) ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	err := temp.Execute(buf, values)
	if err != nil {
		log.Errorf("Template execution failed: %s", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateMonitoringImage(mainImage string, monitoringImage string) (string, error) {
	if monitoringImage != "" {
		_, err := reference.ParseAnyReference(monitoringImage)
		if err != nil {
			return "", err
		}
		return monitoringImage, nil
	}
	img, err := utils.GenerateImageFromString(mainImage)
	if err != nil {
		return "", err
	}
	imgName := img.GetName()
	monitoringImageName := defaultimages.GenerateNamedImageFromMainImage(imgName, imgName.GetTag(), defaultimages.Monitoring)
	return monitoringImageName.FullName, nil
}

func generateReadmeFile(c *Config, mode mode) (*zip.File, error) {
	instructions, err := generateReadme(c, mode)
	if err != nil {
		return nil, err
	}
	return zip.NewFile("README", []byte(instructions), 0), nil
}

// WriteInstructions writes the instructions for the configured cluster
// to the provided writer.
func (c Config) WriteInstructions(w io.Writer) error {
	instructions, err := generateReadme(&c, renderAll)
	if err != nil {
		return err
	}
	fmt.Fprint(w, standardizeWhitespace(instructions))

	if c.PasswordAuto {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Use username '%s' and the following auto-generated password for administrator login (also stored in the 'password' file):\n", basic.DefaultUsername)
		fmt.Fprintf(w, " %s\n", c.Password)
	}
	if c.K8sConfig != nil && c.K8sConfig.Monitoring.PasswordAuto {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Use the following auto-generated password for accessing monitoring (also stored in the 'monitoring/password' file):")
		fmt.Fprintf(w, " %s\n", c.K8sConfig.Monitoring.Password)
	}
	return nil
}

func standardizeWhitespace(instructions string) string {
	return strings.TrimSpace(instructions) + "\n"
}
