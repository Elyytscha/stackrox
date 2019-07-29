package common

import (
	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/grpc/authn/basic"
	"github.com/stackrox/rox/pkg/netutil"
	"github.com/stackrox/rox/roxctl/common/flags"
	"google.golang.org/grpc"
)

// GetGRPCConnection gets a grpc connection to Central with the correct auth
func GetGRPCConnection() (*grpc.ClientConn, error) {
	endpoint := flags.Endpoint()
	serverName := flags.ServerName()
	if serverName == "" {
		var err error
		serverName, _, _, err = netutil.ParseEndpoint(endpoint)
		if err != nil {
			return nil, errors.Wrap(err, "parsing central endpoint")
		}
	}

	if token := env.TokenEnv.Setting(); token != "" {
		return clientconn.GRPCConnectionWithToken(endpoint, serverName, token)
	}
	return clientconn.GRPCConnectionWithBasicAuth(endpoint, serverName, basic.DefaultUsername, flags.Password())
}
