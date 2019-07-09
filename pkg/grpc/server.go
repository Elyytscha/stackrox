package grpc

import (
	"context"
	"crypto/tls"
	golog "log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/pkg/audit"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/contextutil"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authz/deny"
	"github.com/stackrox/rox/pkg/grpc/authz/interceptor"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/grpc/routes"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mtls/verifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxMsgSize = 8 * 1024 * 1024
)

func init() {
	grpc_prometheus.EnableHandlingTimeHistogram()
}

var (
	log = logging.LoggerForModule()
)

type server interface {
	Serve(l net.Listener) error
}

type serverAndListener struct {
	srv      server
	listener net.Listener
	kind     string
}

// APIService is the service interface
type APIService interface {
	RegisterServiceServer(server *grpc.Server)
	RegisterServiceHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

// API listens for new connections on port 443, and redirects them to the gRPC-Gateway
type API interface {
	// Start runs the API in a goroutine, and returns a signal that can be checked for when the API server is started.
	Start() *concurrency.Signal
	// Register adds a new APIService to the list of API services
	Register(services ...APIService)
}

type apiImpl struct {
	apiServices        []APIService
	config             Config
	requestInfoHandler *requestinfo.Handler
}

// A Config configures the server.
type Config struct {
	TLS                verifier.TLSConfigurer
	CustomRoutes       []routes.CustomRoute
	IdentityExtractors []authn.IdentityExtractor
	AuthProviders      authproviders.Registry
	Auditor            audit.Auditor

	PreAuthContextEnrichers  []contextutil.ContextUpdater
	PostAuthContextEnrichers []contextutil.ContextUpdater

	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor

	InsecureLocalEndpoint string
	PublicEndpoint        string

	PlaintextEndpoints EndpointsConfig
}

// NewAPI returns an API object.
func NewAPI(config Config) API {
	return &apiImpl{
		config:             config,
		requestInfoHandler: requestinfo.NewDefaultRequestInfoHandler(),
	}
}

func (a *apiImpl) Start() *concurrency.Signal {
	startedSig := concurrency.NewSignal()
	go a.run(&startedSig)
	return &startedSig
}

func (a *apiImpl) Register(services ...APIService) {
	a.apiServices = append(a.apiServices, services...)
}

func (a *apiImpl) unaryInterceptors() []grpc.UnaryServerInterceptor {
	u := []grpc.UnaryServerInterceptor{
		contextutil.UnaryServerInterceptor(a.requestInfoHandler.UpdateContextForGRPC),
		grpc_prometheus.UnaryServerInterceptor,
		contextutil.UnaryServerInterceptor(authn.ContextUpdater(a.config.IdentityExtractors...)),
	}

	if len(a.config.PreAuthContextEnrichers) > 0 {
		u = append(u, contextutil.UnaryServerInterceptor(a.config.PreAuthContextEnrichers...))
	}

	// Check auth and update the context with the error
	u = append(u, interceptor.AuthContextUpdaterInterceptor())

	if a.config.Auditor != nil {
		// Audit the request
		u = append(u, a.config.Auditor.UnaryServerInterceptor())
	}

	// Check if there was an auth failure and return error if so
	u = append(u, interceptor.AuthCheckerInterceptor())

	if len(a.config.PostAuthContextEnrichers) > 0 {
		u = append(u, contextutil.UnaryServerInterceptor(a.config.PostAuthContextEnrichers...))
	}

	u = append(u, a.config.UnaryInterceptors...)
	u = append(u, a.unaryRecovery())
	return u
}

func (a *apiImpl) streamInterceptors() []grpc.StreamServerInterceptor {
	s := []grpc.StreamServerInterceptor{
		contextutil.StreamServerInterceptor(a.requestInfoHandler.UpdateContextForGRPC),
		grpc_prometheus.StreamServerInterceptor,
		contextutil.StreamServerInterceptor(
			authn.ContextUpdater(a.config.IdentityExtractors...)),
	}
	if len(a.config.PreAuthContextEnrichers) > 0 {
		s = append(s, contextutil.StreamServerInterceptor(a.config.PreAuthContextEnrichers...))
	}

	// Default to deny all access. This forces services to properly override the AuthFunc.
	s = append(s, grpc_auth.StreamServerInterceptor(deny.AuthFunc))

	if len(a.config.PostAuthContextEnrichers) > 0 {
		s = append(s, contextutil.StreamServerInterceptor(a.config.PostAuthContextEnrichers...))
	}

	s = append(s, a.config.StreamInterceptors...)

	s = append(s, a.streamRecovery())
	return s
}

func (a *apiImpl) listenOnLocalEndpoint(server *grpc.Server) error {
	lis, err := net.Listen("tcp", a.config.InsecureLocalEndpoint)
	if err != nil {
		return err
	}

	log.Infof("Launching backend GRPC listener on %v", a.config.InsecureLocalEndpoint)
	// Launch the GRPC listener
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
		log.Fatal("The local API server should never terminate")
	}()
	return nil
}

func (a *apiImpl) connectToLocalEndpoint() (*grpc.ClientConn, error) {
	return grpc.Dial(a.config.InsecureLocalEndpoint, grpc.WithInsecure())
}

func (a *apiImpl) muxer(localConn *grpc.ClientConn) http.Handler {
	contextUpdaters := []contextutil.ContextUpdater{authn.ContextUpdater(a.config.IdentityExtractors...)}
	contextUpdaters = append(contextUpdaters, a.config.PreAuthContextEnrichers...)

	// Interceptors for HTTP/1.1 requests (in order of processing):
	// - RequestInfo handler (consumed by other handlers)
	// - IdentityExtractor
	// - AuthConfigChecker
	httpInterceptors := httputil.ChainInterceptors(
		a.requestInfoHandler.HTTPIntercept,
		contextutil.HTTPInterceptor(contextUpdaters...),
	)

	postAuthHTTPInterceptor := contextutil.HTTPInterceptor(a.config.PostAuthContextEnrichers...)

	mux := http.NewServeMux()
	for _, route := range a.config.CustomRoutes {
		mux.Handle(route.Route, httpInterceptors(route.Handler(postAuthHTTPInterceptor)))
	}

	if a.config.AuthProviders != nil {
		mux.Handle(a.config.AuthProviders.URLPathPrefix(), httpInterceptors(a.config.AuthProviders))
	}

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}),
		runtime.WithMetadata(a.requestInfoHandler.AnnotateMD))
	if localConn != nil {
		for _, service := range a.apiServices {
			if err := service.RegisterServiceHandler(context.Background(), gwMux, localConn); err != nil {
				log.Panicf("failed to register API service: %v", err)
			}
		}
	}
	mux.Handle("/v1/", gziphandler.GzipHandler(gwMux))
	return mux
}

func (a *apiImpl) run(startedSig *concurrency.Signal) {
	tlsConf, err := a.config.TLS.TLSConfig()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(a.streamInterceptors()...),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(a.unaryInterceptors()...),
		),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 40 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	for _, service := range a.apiServices {
		service.RegisterServiceServer(grpcServer)
	}

	var localConn *grpc.ClientConn
	if a.config.InsecureLocalEndpoint != "" {
		if err := a.listenOnLocalEndpoint(grpcServer); err != nil {
			log.Panicf("Could not listen on local endpoint: %v", err)
		}
		localConn, err = a.connectToLocalEndpoint()
		if err != nil {
			log.Panicf("Could not connect to local endpoint: %v", err)
		}
	}

	listener, err := tls.Listen("tcp", a.config.PublicEndpoint, tlsConf)
	if err != nil {
		log.Panicf("Could not listen on public API endpoint: %v", err)
	}

	httpHandler := a.muxer(localConn)

	muxedSrv := &http.Server{
		Handler:   wireOrJSONMuxer(grpcServer, httpHandler),
		ErrorLog:  golog.New(httpErrorLogger{}, "", golog.LstdFlags),
		TLSConfig: tlsConf,
	}

	serverAndListeners := []serverAndListener{
		{
			srv:      muxedSrv,
			listener: listener,
			kind:     "TLS-enabled HTTP/gRPC",
		},
	}

	for _, plaintextEndpoint := range a.config.PlaintextEndpoints.MultiplexedEndpoints {
		plaintextListener, err := net.Listen("tcp", plaintextEndpoint)
		if err != nil {
			log.Panicf("Could not listen on plaintext API endpoint %q: %v", plaintextEndpoint, err)
		}
		serverAndListeners = append(serverAndListeners, serverAndListener{
			srv:      muxedSrv,
			listener: plaintextListener,
			kind:     "Plaintext multiplexed HTTP/gRPC",
		})
	}

	if len(a.config.PlaintextEndpoints.HTTPEndpoints) > 0 {
		httpSrv := &http.Server{
			Handler:  httpHandler,
			ErrorLog: golog.New(httpErrorLogger{}, "", golog.LstdFlags),
		}
		for _, plaintextEndpoint := range a.config.PlaintextEndpoints.HTTPEndpoints {
			plaintextListener, err := net.Listen("tcp", plaintextEndpoint)
			if err != nil {
				log.Panicf("Could not listen on plaintext HTTP API endpoint %q: %v", plaintextEndpoint, err)
			}
			serverAndListeners = append(serverAndListeners, serverAndListener{
				srv:      httpSrv,
				listener: plaintextListener,
				kind:     "Plaintext HTTP",
			})
		}
	}

	for _, plaintextEndpoint := range a.config.PlaintextEndpoints.GRPCEndpoints {
		plaintextListener, err := net.Listen("tcp", plaintextEndpoint)
		if err != nil {
			log.Panicf("Could not listen on plaintext GRPC endpoint %q: %v", plaintextEndpoint, err)
		}
		serverAndListeners = append(serverAndListeners, serverAndListener{
			srv:      grpcServer,
			listener: plaintextListener,
			kind:     "Plaintext gRPC",
		})
	}

	if startedSig != nil {
		startedSig.Signal()
	}

	errC := make(chan error, len(serverAndListeners))

	for _, srvAndListener := range serverAndListeners {
		log.Infof("%s server listening on %s", srvAndListener.kind, srvAndListener.listener.Addr())
		go serveAsync(srvAndListener.srv, srvAndListener.listener, errC)
	}

	if err := <-errC; err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func wireOrJSONMuxer(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func serveAsync(srv server, listener net.Listener, errC chan<- error) {
	if err := srv.Serve(listener); err != nil {
		errC <- err
	}
}
