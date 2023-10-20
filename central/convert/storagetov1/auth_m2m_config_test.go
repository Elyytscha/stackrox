package storagetov1

import (
	"testing"

	convertTestUtils "github.com/stackrox/rox/central/convert/testutils"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestAuthM2MConfig(t *testing.T) {
	config := &storage.AuthMachineToMachineConfig{}
	require.NoError(t, testutils.FullInit(config, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	// Currently required as issuer won't be filled by FullInit.
	config.IssuerConfig = &storage.AuthMachineToMachineConfig_GenericIssuerConfig{
		GenericIssuerConfig: &storage.AuthMachineToMachineConfig_GenericIssuer{
			Issuer: "https://stackrox.io",
		},
	}
	v1Config := AuthM2MConfig(config)

	convertTestUtils.AssertProtoMessageEqual(t, config, v1Config)
}
