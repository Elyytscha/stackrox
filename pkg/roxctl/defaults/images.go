package defaults

import (
	"fmt"

	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/version"
)

// ScannerImage is the Docker image name for the scanner image. Image
// repo changes depending on whether or not this is a release build.
func ScannerImage() string {
	return fmt.Sprintf("%s/scanner:%s", MainImageRegistry(), version.GetScannerVersion())
}

// ScannerDBImage is the Docker image name for the scanner db image
func ScannerDBImage() string {
	return fmt.Sprintf("%s/scanner-db:%s", MainImageRegistry(), version.GetScannerVersion())
}

// MainImage is the Docker image name for the "main" image. Image repo
// changes depending on whether or not this is a release build.
func MainImage() string {
	return fmt.Sprintf("%s:%s", MainImageRepo(), version.GetMainVersion())
}

// MainImageRepo is the Docker image repo for the "main" image. It
// changes depending on whether or not this is a release build.
func MainImageRepo() string {
	return MainImageRegistry() + "/main"
}

// MainImageRegistry is the Docker image registry for the "main" image.
func MainImageRegistry() string {
	if buildinfo.ReleaseBuild {
		return "stackrox.io"
	}
	return "docker.io/stackrox"
}
