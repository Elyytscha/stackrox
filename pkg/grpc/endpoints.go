package grpc

import (
	"crypto/tls"
	"fmt"
	golog "log"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/grpc/alpn"
	"github.com/stackrox/rox/pkg/mtls/verifier"
	"github.com/stackrox/rox/pkg/sliceutils"
	"github.com/stackrox/rox/pkg/tlsutils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	downgradingServer "golang.stackrox.io/grpc-http1/server"
	"google.golang.org/grpc"
)

// EndpointConfig configures an endpoint through which the server is exposed.
type EndpointConfig struct {
	ListenEndpoint string
	Optional       bool
	TLS            verifier.TLSConfigurer

	ServeGRPC, ServeHTTP bool

	NoHTTP2 bool
}

// Kind returns a human-readable description of this endpoint.
func (c *EndpointConfig) Kind() string {
	var sb strings.Builder
	if c.TLS == nil {
		sb.WriteString("Plaintext")
	} else {
		sb.WriteString("TLS-enabled")
	}
	sb.WriteRune(' ')
	if c.ServeHTTP && c.ServeGRPC {
		sb.WriteString("multiplexed HTTP/gRPC")
	} else if c.ServeHTTP {
		sb.WriteString("HTTP")
	} else if c.ServeGRPC {
		sb.WriteString("gRPC")
	} else {
		sb.WriteString("dummy")
	}
	return sb.String()
}

func (c *EndpointConfig) instantiate(httpHandler http.Handler, grpcSrv *grpc.Server) (net.Addr, []serverAndListener, error) {
	lis, err := net.Listen("tcp", asEndpoint(c.ListenEndpoint))
	if err != nil {
		return nil, nil, err
	}

	var httpLis, grpcLis net.Listener

	var result []serverAndListener

	var tlsConf *tls.Config
	if c.TLS != nil {
		tlsConf, err = c.TLS.TLSConfig()
		if err != nil {
			return nil, nil, errors.Wrap(err, "configuring TLS")
		}
	}

	if c.NoHTTP2 {
		tlsConf = tlsConf.Clone()
		tlsConf.NextProtos = sliceutils.StringDifference(tlsConf.NextProtos, []string{"h2", alpn.PureGRPCALPNString})
	}

	if tlsConf != nil {
		if c.ServeGRPC && !c.NoHTTP2 {
			tlsConf = alpn.ApplyPureGRPCALPNConfig(tlsConf)
		}
		lis = tls.NewListener(lis, tlsConf)

		if c.ServeGRPC && c.ServeHTTP {
			protoMap := map[string]*net.Listener{
				alpn.PureGRPCALPNString: &grpcLis,
				"":                      &httpLis,
			}
			tlsutils.ALPNDemux(lis, protoMap, tlsutils.ALPNDemuxConfig{})
		}
	}

	// Default to listen on the main listener, HTTP first
	if c.ServeHTTP && httpLis == nil {
		httpLis = lis
	} else if c.ServeGRPC && grpcLis == nil {
		grpcLis = lis
	}

	if httpLis != nil {
		actualHTTPHandler := httpHandler
		if c.ServeGRPC {
			actualHTTPHandler = downgradingServer.CreateDowngradingHandler(grpcSrv, actualHTTPHandler, downgradingServer.PreferGRPCWeb(true))
		}

		httpSrv := &http.Server{
			Handler:   actualHTTPHandler,
			TLSConfig: tlsConf,
			ErrorLog:  golog.New(httpErrorLogger{}, "", golog.LstdFlags),
		}
		if !c.NoHTTP2 {
			var h2Srv http2.Server
			if err := http2.ConfigureServer(httpSrv, &h2Srv); err != nil {
				log.Warnf("Failed to instantiate endpoint listening at %q for HTTP/2", c.ListenEndpoint)
			} else {
				httpSrv.Handler = h2c.NewHandler(actualHTTPHandler, &h2Srv)
			}
		}
		result = append(result, serverAndListener{
			srv:      httpSrv,
			listener: httpLis,
			endpoint: c,
		})
	}
	if grpcLis != nil {
		result = append(result, serverAndListener{
			srv:      grpcSrv,
			listener: grpcLis,
			endpoint: c,
		})
	}

	return lis.Addr(), result, nil
}

// asEndpoint returns an all-interface endpoint of form `:<port>` if `portOrEndpoint` is a port only (does not contain
// a ':'). Otherwise, `portOrEndpoint` is returned as-is.
func asEndpoint(portOrEndpoint string) string {
	if !strings.ContainsRune(portOrEndpoint, ':') {
		return fmt.Sprintf(":%s", portOrEndpoint)
	}
	return portOrEndpoint
}
