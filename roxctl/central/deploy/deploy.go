package deploy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/renderer"
	"github.com/stackrox/rox/pkg/roxctl"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/zip"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/common/mode"
)

func generateJWTSigningKey(fileMap map[string][]byte) error {
	// Generate the private key that we will use to sign JWTs for API keys.
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return errors.Wrap(err, "couldn't generate private key")
	}
	fileMap["jwt-key.der"] = x509.MarshalPKCS1PrivateKey(privateKey)
	return nil
}

func generateMTLSFiles(fileMap map[string][]byte) (caCert, caKey []byte, err error) {
	// Add MTLS files
	serial, err := mtls.RandomSerial()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not generate a serial number")
	}
	req := csr.CertificateRequest{
		CN:           mtls.ServiceCACommonName,
		KeyRequest:   csr.NewBasicKeyRequest(),
		SerialNumber: serial.String(),
	}
	caCert, _, caKey, err = initca.New(&req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not generate keypair")
	}
	fileMap["ca.pem"] = caCert
	fileMap["ca-key.pem"] = caKey

	cert, err := mtls.IssueNewCertFromCA(mtls.CentralSubject, caCert, caKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not issue central cert")
	}
	fileMap["cert.pem"] = cert.CertPEM
	fileMap["key.pem"] = cert.KeyPEM

	scannerCert, err := mtls.IssueNewCertFromCA(mtls.ScannerSubject, caCert, caKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not issue scanner cert")
	}
	fileMap["scanner-cert.pem"] = scannerCert.CertPEM
	fileMap["scanner-key.pem"] = scannerCert.KeyPEM

	scannerDBCert, err := mtls.IssueNewCertFromCA(mtls.ScannerDBSubject, caCert, caKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not issue scanner db cert")
	}
	fileMap["scanner-db-cert.pem"] = scannerDBCert.CertPEM
	fileMap["scanner-db-key.pem"] = scannerDBCert.KeyPEM

	fileMap["scanner-db-password"] = []byte(renderer.CreatePassword())

	return caCert, caKey, nil
}

func generateMonitoringFiles(fileMap map[string][]byte, caCert, caKey []byte) error {
	monitoringCert, err := mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_DB_SERVICE, Identifier: "Monitoring DB"},
		caCert, caKey)
	if err != nil {
		return err
	}
	fileMap["monitoring-db-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-db-key.pem"] = monitoringCert.KeyPEM

	monitoringCert, err = mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_UI_SERVICE, Identifier: "Monitoring UI"},
		caCert, caKey)
	if err != nil {
		return err
	}

	fileMap["monitoring-ui-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-ui-key.pem"] = monitoringCert.KeyPEM

	monitoringCert, err = mtls.IssueNewCertFromCA(mtls.Subject{ServiceType: storage.ServiceType_MONITORING_CLIENT_SERVICE, Identifier: "Monitoring Client"},
		caCert, caKey)
	if err != nil {
		return err
	}

	fileMap["monitoring-client-cert.pem"] = monitoringCert.CertPEM
	fileMap["monitoring-client-key.pem"] = monitoringCert.KeyPEM

	return nil
}

// OutputZip renders a deployment bundle. The deployment bundle can either be
// written directly into a directory, or as a zipfile to STDOUT.
func OutputZip(config renderer.Config) error {
	fmt.Fprint(os.Stderr, "Generating deployment bundle... \n")

	wrapper := zip.NewWrapper()

	if config.ClusterType == storage.ClusterType_GENERIC_CLUSTER {
		return errors.Errorf("invalid cluster type: %s", config.ClusterType)
	}

	config.SecretsByteMap = make(map[string][]byte)
	if err := generateJWTSigningKey(config.SecretsByteMap); err != nil {
		return err
	}

	if len(config.LicenseData) > 0 {
		config.SecretsByteMap["central-license"] = config.LicenseData
	}

	if len(config.DefaultTLSCertPEM) > 0 {
		config.SecretsByteMap["default-tls.crt"] = config.DefaultTLSCertPEM
		config.SecretsByteMap["default-tls.key"] = config.DefaultTLSKeyPEM
	}

	config.Environment = make(map[string]string)
	for _, flag := range features.Flags {
		if value := os.Getenv(flag.EnvVar()); value != "" {
			config.Environment[flag.EnvVar()] = strconv.FormatBool(flag.Enabled())
		}
	}

	htpasswd, err := renderer.GenerateHtpasswd(&config)
	if err != nil {
		return err
	}

	for _, setting := range env.Settings {
		if _, ok := os.LookupEnv(setting.EnvVar()); ok {
			config.Environment[setting.EnvVar()] = setting.Setting()
		}
	}
	if config.K8sConfig != nil {
		config.Environment[env.OfflineModeEnv.EnvVar()] = strconv.FormatBool(config.K8sConfig.OfflineMode)

		if features.Telemetry.Enabled() {
			if config.K8sConfig.EnableTelemetry {
				fmt.Fprintln(os.Stderr, "NOTE: Unless run in offline mode, StackRox Kubernetes Security Platform collects and transmits aggregated usage and system health information.  If you want to OPT OUT from this, re-generate the deployment bundle with the '--enable-telemetry=false' flag")
			}
			config.Environment[env.InitialTelemetryEnabledEnv.EnvVar()] = strconv.FormatBool(config.K8sConfig.EnableTelemetry)
		}
	}

	config.SecretsByteMap["htpasswd"] = htpasswd
	wrapper.AddFiles(zip.NewFile("password", []byte(config.Password+"\n"), zip.Sensitive))

	cert, key, err := generateMTLSFiles(config.SecretsByteMap)
	if err != nil {
		return err
	}

	if config.K8sConfig != nil && config.K8sConfig.Monitoring.Type.OnPrem() {
		if err := generateMonitoringFiles(config.SecretsByteMap, cert, key); err != nil {
			return err
		}

		if config.K8sConfig.Monitoring.Password == "" {
			config.K8sConfig.Monitoring.Password = renderer.CreatePassword()
			config.K8sConfig.Monitoring.PasswordAuto = true
		}

		config.SecretsByteMap["monitoring-password"] = []byte(config.K8sConfig.Monitoring.Password)
		wrapper.AddFiles(zip.NewFile("monitoring/password", []byte(config.K8sConfig.Monitoring.Password+"\n"), zip.Sensitive))
	}

	files, err := renderer.Render(config)
	if err != nil {
		return errors.Wrap(err, "could not render files")
	}
	wrapper.AddFiles(files...)

	var outputPath string
	if roxctl.InMainImage() {
		bytes, err := wrapper.Zip()
		if err != nil {
			return errors.Wrap(err, "error generating zip file")
		}
		_, err = os.Stdout.Write(bytes)
		if err != nil {
			return errors.Wrap(err, "couldn't write zip file")
		}
	} else {
		var err error
		outputPath, err = wrapper.Directory(config.OutputDir)
		if err != nil {
			return errors.Wrap(err, "error generating directory for Central output")
		}
	}

	fmt.Fprintln(os.Stderr, "Done!")
	fmt.Fprintln(os.Stderr)

	if outputPath != "" {
		fmt.Fprintf(os.Stderr, "Wrote central bundle to %q\n", outputPath)
		fmt.Fprintln(os.Stderr)
	}

	if err := config.WriteInstructions(os.Stderr); err != nil {
		return err
	}
	return nil
}

func interactive() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Interactive runs the CLI in interactive mode with user prompts",
		RunE: func(_ *cobra.Command, args []string) error {
			c := Command()
			c.SilenceUsage = true
			return runInteractive(c)
		},
		SilenceUsage: true,
	}
}

// Command defines the deploy command tree
func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "generate",
		Short: "Generate creates the required YAML files to deploy StackRox Central.",
		Long:  "Generate creates the required YAML files to deploy StackRox Central.",
		Run: func(c *cobra.Command, _ []string) {
			if !mode.IsInInteractiveMode() {
				_ = c.Help()
			}
		},
	}
	c.PersistentFlags().StringVar(&cfg.Password, "password", "", "administrator password (default: autogenerated)")
	utils.Must(c.PersistentFlags().SetAnnotation("password", flags.PasswordKey, []string{"true"}))

	c.PersistentFlags().Var(&flags.LicenseVar{
		Data: &cfg.LicenseData,
	}, "license", flags.LicenseUsage)

	utils.Must(
		c.PersistentFlags().SetAnnotation("license", flags.InteractiveUsageKey, []string{flags.LicenseUsageInteractive}),
		c.PersistentFlags().SetAnnotation("license", flags.OptionalKey, []string{"true"}),
	)

	c.PersistentFlags().Var(&flags.FileContentsVar{
		Data: &cfg.DefaultTLSCertPEM,
	}, "default-tls-cert", "PEM cert bundle file")
	utils.Must(c.PersistentFlags().SetAnnotation("default-tls-cert", flags.OptionalKey, []string{"true"}))

	c.PersistentFlags().Var(&flags.FileContentsVar{
		Data: &cfg.DefaultTLSKeyPEM,
	}, "default-tls-key", "PEM private key file")
	utils.Must(
		c.PersistentFlags().SetAnnotation("default-tls-key", flags.DependenciesKey, []string{"default-tls-cert"}),
		c.PersistentFlags().SetAnnotation("default-tls-key", flags.MandatoryKey, []string{"true"}),
	)

	c.PersistentFlags().VarPF(
		flags.ForSettingWithOptions(env.PlaintextEndpoints, flags.SettingVarOpts{
			Validator: func(spec string) error {
				cfg := grpc.EndpointsConfig{}
				if err := cfg.AddFromParsedSpec(spec); err != nil {
					return err
				}
				return cfg.Validate()
			},
		}), "plaintext-endpoints", "",
		"The ports or endpoints to use for plaintext (unencrypted) exposure; comma-separated list.")
	utils.Must(
		c.PersistentFlags().SetAnnotation("plaintext-endpoints", flags.NoInteractiveKey, []string{"true"}))

	c.AddCommand(interactive())

	c.AddCommand(k8s())
	c.AddCommand(openshift())
	return c
}

func runInteractive(cmd *cobra.Command) error {
	mode.SetInteractiveMode()
	// Overwrite os.Args because cobra uses them
	os.Args = walkTree(cmd)
	return cmd.Execute()
}
