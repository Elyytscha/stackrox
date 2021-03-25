package renderer

import (
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	imageUtils "github.com/stackrox/rox/pkg/images/utils"
	kubernetesPkg "github.com/stackrox/rox/pkg/kubernetes"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/pkg/zip"
)

var (
	assetFileNameMap = FileNameMap{}
)

func init() {
	assetFileNameMap.AddWithName("assets/docker-auth.sh", "docker-auth.sh")
}

// mode is the mode we want the renderer to function in.
//go:generate stringer -type=mode
type mode int

const (
	// renderAll renders all objects (central+scanner).
	renderAll mode = iota
	// scannerOnly renders only the scanner.
	scannerOnly
	// centralTLSOnly renders only the central tls secret.
	centralTLSOnly
	// scannerTLSOnly renders only the scanner tls secret
	scannerTLSOnly
)

func postProcessConfig(c *Config, mode mode) error {
	// Make all items in SecretsByteMap base64 encoded
	c.SecretsBase64Map = make(map[string]string)
	for k, v := range c.SecretsByteMap {
		c.SecretsBase64Map[k] = base64.StdEncoding.EncodeToString(v)
	}
	if mode == centralTLSOnly || mode == scannerTLSOnly {
		return nil
	}
	if c.ClusterType == storage.ClusterType_KUBERNETES_CLUSTER {
		c.K8sConfig.Command = "kubectl"
	} else {
		c.K8sConfig.Command = "oc"
	}

	configureImageOverrides(c)

	var err error
	if mode == renderAll {
		c.K8sConfig.Registry, err = kubernetesPkg.GetResolvedRegistry(c.K8sConfig.MainImage)
		if err != nil {
			return err
		}
	}

	c.K8sConfig.ScannerRegistry, err = kubernetesPkg.GetResolvedRegistry(c.K8sConfig.ScannerImage)
	if err != nil {
		return err
	}
	if c.K8sConfig.Registry != c.K8sConfig.ScannerRegistry {
		c.K8sConfig.ScannerSecretName = "stackrox-scanner"
	} else {
		c.K8sConfig.ScannerSecretName = "stackrox"
	}

	if mode == renderAll {
		if err := injectImageTags(c); err != nil {
			return err
		}
	}

	return nil
}

// Render renders a bunch of zip files based on the given config.
func Render(c Config) ([]*zip.File, error) {
	return render(c, renderAll, nil)
}

// RenderWithOverrides renders a bunch of zip files based on the given config, allowing to selectively override some of
// the bundled files.
func RenderWithOverrides(c Config, centralOverrides map[string]func() io.ReadCloser) ([]*zip.File, error) {
	return render(c, renderAll, centralOverrides)
}

// RenderScannerOnly renders the zip files for the scanner based on the given config.
func RenderScannerOnly(c Config) ([]*zip.File, error) {
	return render(c, scannerOnly, nil)
}

func renderAndExtractSingleFileContents(c Config, mode mode) ([]byte, error) {
	files, err := render(c, mode, nil)
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, utils.Should(errors.Errorf("got unexpected number of files when rendering in mode %s: %d", mode, len(files)))
	}
	return files[0].Content, nil
}

// RenderCentralTLSSecretOnly renders just the file that contains the central-tls secret.
func RenderCentralTLSSecretOnly(c Config) ([]byte, error) {
	return renderAndExtractSingleFileContents(c, centralTLSOnly)
}

// RenderScannerTLSSecretOnly renders just the file that contains the scanner-tls secret.
func RenderScannerTLSSecretOnly(c Config) ([]byte, error) {
	return renderAndExtractSingleFileContents(c, scannerTLSOnly)
}

func render(c Config, mode mode, centralOverrides map[string]func() io.ReadCloser) ([]*zip.File, error) {
	err := postProcessConfig(&c, mode)
	if err != nil {
		return nil, err
	}

	return renderNew(c, mode)
}

func getTag(imageStr string) (string, error) {
	imageName, err := imageUtils.GenerateImageFromString(imageStr)
	if err != nil {
		return "", err
	}
	return imageName.GetName().GetTag(), nil
}

func injectImageTags(c *Config) error {
	var err error
	c.K8sConfig.MainImageTag, err = getTag(c.K8sConfig.MainImage)
	if err != nil {
		return err
	}
	return nil
}
