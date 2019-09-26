package testutils

import (
	"context"
	"os"
	"testing"

	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/netutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	usernameEnvVar = "ROX_USERNAME"
	passwordEnvVar = "ROX_PASSWORD"

	apiEndpointEnvVar = "API_ENDPOINT"
)

func mustGetEnvVar(t *testing.T, envVar string) string {
	value := os.Getenv(envVar)
	require.NotEmpty(t, value, "Please set %s", envVar)
	return value
}

// RoxUsername returns the rox basic auth username (or fails the test if not found).
func RoxUsername(t *testing.T) string {
	return mustGetEnvVar(t, usernameEnvVar)
}

// RoxPassword returns the rox basic auth password (or fails the test if not found).
func RoxPassword(t *testing.T) string {
	return mustGetEnvVar(t, passwordEnvVar)
}

// RoxAPIEndpoint returns the central API endpoint (or fails the test if not found).
func RoxAPIEndpoint(t *testing.T) string {
	return mustGetEnvVar(t, apiEndpointEnvVar)
}

// GRPCConnectionToCentral returns a GRPC connection to Central, which can be used in E2E tests.
// It fatals the test if there's an error.
func GRPCConnectionToCentral(t *testing.T) *grpc.ClientConn {
	endpoint := RoxAPIEndpoint(t)
	host, _, _, err := netutil.ParseEndpoint(endpoint)
	require.NoError(t, err)

	opts := clientconn.Options{
		TLS: clientconn.TLSConfigOptions{
			InsecureSkipVerify: true,
			ServerName:         host,
		},
	}
	opts.ConfigureBasicAuth(RoxUsername(t), RoxPassword(t))
	conn, err := clientconn.GRPCConnection(context.Background(), mtls.CentralSubject, endpoint, opts)
	require.NoError(t, err)
	return conn
}
