package cluster

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/netutil"
	"github.com/stackrox/rox/pkg/stringutils"
	"github.com/stackrox/rox/pkg/urlfmt"
)

// Validate validates a cluster object
func Validate(cluster *storage.Cluster) *errorhelpers.ErrorList {
	errorList := errorhelpers.NewErrorList("Cluster Validation")
	if cluster.GetName() == "" {
		errorList.AddString("Cluster name is required")
	}
	if _, err := reference.ParseAnyReference(cluster.GetMainImage()); err != nil {
		errorList.AddError(errors.Wrapf(err, "invalid main image '%s'", cluster.GetMainImage()))
	}
	if cluster.GetCollectorImage() != "" {
		ref, err := reference.ParseAnyReference(cluster.GetCollectorImage())
		if err != nil {
			errorList.AddError(errors.Wrapf(err, "invalid collector image '%s'", cluster.GetCollectorImage()))
		}
		namedTagged, ok := ref.(reference.NamedTagged)
		if ok {
			errorList.AddError(fmt.Errorf("collector image may not specify a tag.  Please "+
				"remove tag '%s' to continue", namedTagged.Tag()))
		}
	}
	if cluster.GetCentralApiEndpoint() == "" {
		errorList.AddString("Central API Endpoint is required")
	} else if !strings.Contains(cluster.GetCentralApiEndpoint(), ":") {
		errorList.AddString("Central API Endpoint must have port specified")
	}

	if stringutils.ContainsWhitespace(cluster.GetCentralApiEndpoint()) {
		errorList.AddString("Central API endpoint cannot contain whitespace")
	}

	centralEndpoint := urlfmt.FormatURL(cluster.GetCentralApiEndpoint(), urlfmt.NONE, urlfmt.NoTrailingSlash)
	_, _, _, err := netutil.ParseEndpoint(centralEndpoint)
	if err != nil {
		errorList.AddString(fmt.Sprintf("Central API Endpoint must be a valid endpoint. Error: %s", err))
	}
	cluster.CentralApiEndpoint = centralEndpoint

	return errorList
}
