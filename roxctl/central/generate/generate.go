package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/certgen"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/ioutils"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/renderer"
	"github.com/stackrox/rox/pkg/roxctl"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/zip"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/common/mode"
	"github.com/stackrox/rox/roxctl/common/util"
)

func generateJWTSigningKey(fileMap map[string][]byte) error {
	// Generate the private key that we will use to sign JWTs for API keys.
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return errors.Wrap(err, "couldn't generate private key")
	}
	jwtKey := x509.MarshalPKCS1PrivateKey(privateKey)
	fileMap["jwt-key.der"] = jwtKey
	fileMap["jwt-key.pem"] = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: jwtKey,
	})
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

	if err := certgen.IssueCentralCert(fileMap); err != nil {
		return nil, nil, err
	}

	if err := certgen.IssueScannerCerts(fileMap); err != nil {
		return nil, nil, err
	}

	fileMap["scanner-db-password"] = []byte(renderer.CreatePassword())

	return caCert, caKey, nil
}

func openFileFunc(path string) func() io.ReadCloser {
	return func() io.ReadCloser {
		f, err := os.Open(path)
		if err != nil {
			return ioutil.NopCloser(ioutils.ErrorReader(err))
		}
		return f
	}
}

func buildConfigFileOverridesMap(filePathMap map[string]string) map[string]func() io.ReadCloser {
	if len(filePathMap) == 0 {
		return nil
	}

	fileOverridesMap := make(map[string]func() io.ReadCloser, len(filePathMap))
	for configFilePath, replacementFilePath := range filePathMap {
		fileOverridesMap[path.Join("config", configFilePath)] = openFileFunc(replacementFilePath)
	}
	return fileOverridesMap
}

func createBundle(config renderer.Config) (*zip.Wrapper, error) {
	wrapper := zip.NewWrapper()

	if config.ClusterType == storage.ClusterType_GENERIC_CLUSTER {
		return nil, errors.Errorf("invalid cluster type: %s", config.ClusterType)
	}

	config.SecretsByteMap = make(map[string][]byte)
	if err := generateJWTSigningKey(config.SecretsByteMap); err != nil {
		return nil, err
	}

	if len(config.LicenseData) > 0 {
		config.SecretsByteMap["central-license"] = config.LicenseData
	}

	if len(config.DefaultTLSCertPEM) > 0 {
		config.SecretsByteMap["default-tls.crt"] = config.DefaultTLSCertPEM
		config.SecretsByteMap["default-tls.key"] = config.DefaultTLSKeyPEM
	}

	config.Environment = make(map[string]string)
	// Feature flags can only be overridden on release builds.
	if !buildinfo.ReleaseBuild {
		for _, flag := range features.Flags {
			if value := os.Getenv(flag.EnvVar()); value != "" {
				config.Environment[flag.EnvVar()] = strconv.FormatBool(flag.Enabled())
			}
		}
	}

	htpasswd, err := renderer.GenerateHtpasswd(&config)
	if err != nil {
		return nil, err
	}

	for _, setting := range env.Settings {
		if _, ok := os.LookupEnv(setting.EnvVar()); ok {
			config.Environment[setting.EnvVar()] = setting.Setting()
		}
	}
	if config.K8sConfig != nil {
		config.Environment[env.OfflineModeEnv.EnvVar()] = strconv.FormatBool(config.K8sConfig.OfflineMode)

		if config.K8sConfig.EnableTelemetry {
			fmt.Fprintln(os.Stderr, "NOTE: Unless run in offline mode, StackRox Kubernetes Security Platform collects and transmits aggregated usage and system health information.  If you want to OPT OUT from this, re-generate the deployment bundle with the '--enable-telemetry=false' flag")
		}
		config.Environment[env.InitialTelemetryEnabledEnv.EnvVar()] = strconv.FormatBool(config.K8sConfig.EnableTelemetry)
	}

	config.SecretsByteMap["htpasswd"] = htpasswd
	wrapper.AddFiles(zip.NewFile("password", []byte(config.Password+"\n"), zip.Sensitive))

	if _, _, err := generateMTLSFiles(config.SecretsByteMap); err != nil {
		return nil, err
	}

	files, err := renderer.RenderWithOverrides(config, buildConfigFileOverridesMap(config.ConfigFileOverrides))
	if err != nil {
		return nil, errors.Wrap(err, "could not render files")
	}
	wrapper.AddFiles(files...)

	return wrapper, nil
}

// OutputZip renders a deployment bundle. The deployment bundle can either be
// written directly into a directory, or as a zipfile to STDOUT.
func OutputZip(config renderer.Config) error {
	fmt.Fprint(os.Stderr, "Generating deployment bundle... \n")

	wrapper, err := createBundle(config)
	if err != nil {
		return err
	}

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
		Use: "interactive",
		RunE: util.RunENoArgs(func(*cobra.Command) error {
			c := Command()
			c.SilenceUsage = true
			return runInteractive(c)
		}),
		SilenceUsage: true,
	}
}

// Command defines the generate command tree
func Command() *cobra.Command {
	c := &cobra.Command{
		Use: "generate",
	}
	c.PersistentFlags().StringVarP(&cfg.Password, "password", "p", "", "administrator password (default: autogenerated)")
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
		flags.ForSetting(env.PlaintextEndpoints), "plaintext-endpoints", "",
		"The ports or endpoints to use for plaintext (unencrypted) exposure; comma-separated list.")
	utils.Must(
		c.PersistentFlags().SetAnnotation("plaintext-endpoints", flags.NoInteractiveKey, []string{"true"}))

	c.PersistentFlags().Var(&flags.FileMapVar{
		FileMap: &cfg.ConfigFileOverrides,
	}, "with-config-file", "Use the given local file(s) to override default config files")
	utils.Must(
		c.PersistentFlags().MarkHidden("with-config-file"))

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
