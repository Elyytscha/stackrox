package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	alertService "github.com/stackrox/rox/central/alert/service"
	apiTokenService "github.com/stackrox/rox/central/apitoken/service"
	"github.com/stackrox/rox/central/audit"
	authService "github.com/stackrox/rox/central/auth/service"
	"github.com/stackrox/rox/central/auth/userpass"
	authProviderDS "github.com/stackrox/rox/central/authprovider/datastore"
	authproviderService "github.com/stackrox/rox/central/authprovider/service"
	"github.com/stackrox/rox/central/certgen"
	"github.com/stackrox/rox/central/cli"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	clusterService "github.com/stackrox/rox/central/cluster/service"
	clusterInitService "github.com/stackrox/rox/central/clusterinit/service"
	clustersHelmConfig "github.com/stackrox/rox/central/clusters/helmconfig"
	clustersZip "github.com/stackrox/rox/central/clusters/zip"
	complianceDatastore "github.com/stackrox/rox/central/compliance/datastore"
	complianceHandlers "github.com/stackrox/rox/central/compliance/handlers"
	complianceManager "github.com/stackrox/rox/central/compliance/manager"
	complianceManagerService "github.com/stackrox/rox/central/compliance/manager/service"
	complianceService "github.com/stackrox/rox/central/compliance/service"
	configService "github.com/stackrox/rox/central/config/service"
	credentialExpiryService "github.com/stackrox/rox/central/credentialexpiry/service"
	"github.com/stackrox/rox/central/cve/csv"
	"github.com/stackrox/rox/central/cve/fetcher"
	cveService "github.com/stackrox/rox/central/cve/service"
	"github.com/stackrox/rox/central/cve/suppress"
	debugService "github.com/stackrox/rox/central/debug/service"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	deploymentService "github.com/stackrox/rox/central/deployment/service"
	detectionService "github.com/stackrox/rox/central/detection/service"
	developmentService "github.com/stackrox/rox/central/development/service"
	"github.com/stackrox/rox/central/docs"
	"github.com/stackrox/rox/central/ed"
	"github.com/stackrox/rox/central/endpoints"
	_ "github.com/stackrox/rox/central/externalbackups/plugins/all" // Import all of the external backup plugins
	backupService "github.com/stackrox/rox/central/externalbackups/service"
	featureFlagService "github.com/stackrox/rox/central/featureflags/service"
	"github.com/stackrox/rox/central/globaldb"
	dbAuthz "github.com/stackrox/rox/central/globaldb/authz"
	globaldbHandlers "github.com/stackrox/rox/central/globaldb/handlers"
	backupRestoreService "github.com/stackrox/rox/central/globaldb/v2backuprestore/service"
	graphqlHandler "github.com/stackrox/rox/central/graphql/handler"
	groupService "github.com/stackrox/rox/central/group/service"
	"github.com/stackrox/rox/central/grpc/metrics"
	helmHandler "github.com/stackrox/rox/central/helm/handler"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	imageService "github.com/stackrox/rox/central/image/service"
	"github.com/stackrox/rox/central/imageintegration"
	iiDatastore "github.com/stackrox/rox/central/imageintegration/datastore"
	iiService "github.com/stackrox/rox/central/imageintegration/service"
	iiStore "github.com/stackrox/rox/central/imageintegration/store"
	integrationHealthService "github.com/stackrox/rox/central/integrationhealth/service"
	"github.com/stackrox/rox/central/jwt"
	licenseEnforcer "github.com/stackrox/rox/central/license/enforcer"
	licenseService "github.com/stackrox/rox/central/license/service"
	licenseSingletons "github.com/stackrox/rox/central/license/singleton"
	logimbueHandler "github.com/stackrox/rox/central/logimbue/handler"
	metadataService "github.com/stackrox/rox/central/metadata/service"
	namespaceService "github.com/stackrox/rox/central/namespace/service"
	networkEntityDataStore "github.com/stackrox/rox/central/networkgraph/entity/datastore"
	"github.com/stackrox/rox/central/networkgraph/entity/gatherer"
	networkFlowService "github.com/stackrox/rox/central/networkgraph/service"
	networkPolicyService "github.com/stackrox/rox/central/networkpolicies/service"
	nodeService "github.com/stackrox/rox/central/node/service"
	"github.com/stackrox/rox/central/notifier/processor"
	notifierService "github.com/stackrox/rox/central/notifier/service"
	_ "github.com/stackrox/rox/central/notifiers/all" // These imports are required to register things from the respective packages.
	pingService "github.com/stackrox/rox/central/ping/service"
	podService "github.com/stackrox/rox/central/pod/service"
	policyDataStore "github.com/stackrox/rox/central/policy/datastore"
	policyService "github.com/stackrox/rox/central/policy/service"
	probeUploadService "github.com/stackrox/rox/central/probeupload/service"
	processIndicatorService "github.com/stackrox/rox/central/processindicator/service"
	processWhitelistDataStore "github.com/stackrox/rox/central/processwhitelist/datastore"
	processWhitelistService "github.com/stackrox/rox/central/processwhitelist/service"
	"github.com/stackrox/rox/central/pruning"
	rbacService "github.com/stackrox/rox/central/rbac/service"
	"github.com/stackrox/rox/central/reprocessor"
	"github.com/stackrox/rox/central/risk/handlers/timeline"
	"github.com/stackrox/rox/central/role"
	roleDataStore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/central/role/mapper"
	"github.com/stackrox/rox/central/role/resources"
	roleService "github.com/stackrox/rox/central/role/service"
	centralSAC "github.com/stackrox/rox/central/sac"
	sacService "github.com/stackrox/rox/central/sac/service"
	"github.com/stackrox/rox/central/sac/transitional"
	"github.com/stackrox/rox/central/scanner"
	scannerDefinitionsHandler "github.com/stackrox/rox/central/scannerdefinitions/handler"
	searchService "github.com/stackrox/rox/central/search/service"
	secretService "github.com/stackrox/rox/central/secret/service"
	sensorService "github.com/stackrox/rox/central/sensor/service"
	"github.com/stackrox/rox/central/sensor/service/connection"
	"github.com/stackrox/rox/central/sensor/service/pipeline/all"
	sensorUpgradeControlService "github.com/stackrox/rox/central/sensorupgrade/controlservice"
	sensorUpgradeService "github.com/stackrox/rox/central/sensorupgrade/service"
	sensorUpgradeConfigStore "github.com/stackrox/rox/central/sensorupgradeconfig/datastore"
	serviceAccountService "github.com/stackrox/rox/central/serviceaccount/service"
	siStore "github.com/stackrox/rox/central/serviceidentities/datastore"
	siService "github.com/stackrox/rox/central/serviceidentities/service"
	"github.com/stackrox/rox/central/splunk"
	summaryService "github.com/stackrox/rox/central/summary/service"
	telemetryService "github.com/stackrox/rox/central/telemetry/service"
	"github.com/stackrox/rox/central/tlsconfig"
	"github.com/stackrox/rox/central/ui"
	userService "github.com/stackrox/rox/central/user/service"
	"github.com/stackrox/rox/central/version"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/auth/authproviders/iap"
	"github.com/stackrox/rox/pkg/auth/authproviders/oidc"
	"github.com/stackrox/rox/pkg/auth/authproviders/saml"
	authProviderUserpki "github.com/stackrox/rox/pkg/auth/authproviders/userpki"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/config"
	"github.com/stackrox/rox/pkg/debughandler"
	"github.com/stackrox/rox/pkg/devbuild"
	"github.com/stackrox/rox/pkg/devmode"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/features"
	pkgGRPC "github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authn/service"
	"github.com/stackrox/rox/pkg/grpc/authn/servicecerttoken"
	"github.com/stackrox/rox/pkg/grpc/authn/tokenbased"
	authnUserpki "github.com/stackrox/rox/pkg/grpc/authn/userpki"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/allow"
	"github.com/stackrox/rox/pkg/grpc/authz/or"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/grpc/routes"
	"github.com/stackrox/rox/pkg/httputil/proxy"
	"github.com/stackrox/rox/pkg/logging"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/osutils"
	"github.com/stackrox/rox/pkg/premain"
	"github.com/stackrox/rox/pkg/sac"
	pkgVersion "github.com/stackrox/rox/pkg/version"
)

var (
	log = logging.LoggerForModule()

	authProviderBackendFactories = map[string]authproviders.BackendFactoryCreator{
		oidc.TypeName:                oidc.NewFactory,
		"auth0":                      oidc.NewFactory, // legacy
		saml.TypeName:                saml.NewFactory,
		authProviderUserpki.TypeName: authProviderUserpki.NewFactoryFactory(tlsconfig.ManagerInstance()),
		iap.TypeName:                 iap.NewFactory,
	}

	imageIntegrationContext = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.ImageIntegration),
		))
)

const (
	ssoURLPathPrefix     = "/sso/"
	tokenRedirectURLPath = "/auth/response/generic"

	grpcServerWatchdogTimeout = 20 * time.Second

	maxServiceCertTokenLeeway = 1 * time.Minute

	proxyConfigPath = "/run/secrets/stackrox.io/proxy-config"
	proxyConfigFile = "config.yaml"
)

func init() {
	if !proxy.UseWithDefaultTransport() {
		log.Warn("Failed to use proxy transport with default HTTP transport. Some proxy features may not work.")
	}
}

func runSafeMode() {
	log.Info("Started Central up in safe mode. Sleeping forever...")

	signalsC := make(chan os.Signal, 1)
	signal.Notify(signalsC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-signalsC
	log.Infof("Caught %s signal", sig)
	log.Info("Central terminated")
}

func main() {
	premain.StartMain()

	conf, err := config.ReadConfig()
	if err != nil || conf.Maintenance.SafeMode {
		if err != nil {
			log.Errorf("error reading config file: %v. Starting up in safe mode", err)
		}
		runSafeMode()
		return
	}

	proxy.WatchProxyConfig(context.Background(), proxyConfigPath, proxyConfigFile, true)

	if devbuild.IsEnabled() {
		debughandler.MustStartServerAsync("")

		devmode.StartBinaryWatchdog("central")

		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	}

	log.Infof("Running StackRox Version: %s", pkgVersion.GetMainVersion())
	ensureDB()

	// Now that we verified that the DB can be loaded, remove the .backup directory
	if err := os.RemoveAll(filepath.Join(migrations.DBMountPath, ".backup")); err != nil {
		log.Errorf("Failed to remove backup DB: %v", err)
	}

	// Start the prometheus metrics server
	pkgMetrics.NewDefaultHTTPServer().RunForever()
	pkgMetrics.GatherThrottleMetricsForever(pkgMetrics.CentralSubsystem.String())

	var restartingFlag concurrency.Flag

	licenseMgr := licenseSingletons.ManagerSingleton()
	initialLicense, err := licenseMgr.Initialize(licenseEnforcer.New(&restartingFlag))
	if err != nil {
		log.Fatalf("Could not initialize license manager: %v", err)
	}

	if initialLicense == nil {
		log.Error("*** No valid license found")
		log.Error("*** ")
		log.Error("*** Server starting in limited mode until license activated")
		go startLimitedModeServer(&restartingFlag)
		waitForTerminationSignal()
		return
	}

	log.Info("Extracting StackRox data ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := ed.ED(ctx); err != nil {
		log.Fatalf("Could not extract data: %v", err)
	}
	log.Info("Successfully extracted StackRox data")

	go startMainServer(&restartingFlag)

	waitForTerminationSignal()
}

func ensureDB() {
	err := version.Ensure(globaldb.GetGlobalDB(), globaldb.GetRocksDB())
	if err != nil {
		log.Panicf("DB version check failed. You may need to run migrations: %v", err)
	}
}

type invalidLicenseFactory struct {
	restartingFlag *concurrency.Flag
}

func (f invalidLicenseFactory) ServicesToRegister(authproviders.Registry) []pkgGRPC.APIService {
	return []pkgGRPC.APIService{
		licenseService.New(true, licenseSingletons.ManagerSingleton()),
		metadataService.New(f.restartingFlag, licenseSingletons.ManagerSingleton()),
		pingService.Singleton(), // required for dev scripts & health checking
	}
}

func (invalidLicenseFactory) StartServices() {
}

func (invalidLicenseFactory) CustomRoutes() (customRoutes []routes.CustomRoute) {
	return []routes.CustomRoute{uiRoute()}
}

func startLimitedModeServer(restartingFlag *concurrency.Flag) {
	startGRPCServer(invalidLicenseFactory{
		restartingFlag: restartingFlag,
	})
}

type serviceFactory interface {
	CustomRoutes() (customRoutes []routes.CustomRoute)
	ServicesToRegister(authproviders.Registry) []pkgGRPC.APIService
	StartServices()
}

type defaultFactory struct {
	restartingFlag *concurrency.Flag
}

func (defaultFactory) StartServices() {
	if err := complianceManager.Singleton().Start(); err != nil {
		log.Panicf("could not start compliance manager: %v", err)
	}
	reprocessor.Singleton().Start()
	suppress.Singleton().Start()
	pruning.Singleton().Start()

	if features.NetworkGraphExternalSrcs.Enabled() {
		gatherer.Singleton().Start()
	}

	go registerDelayedIntegrations(iiStore.DelayedIntegrations)
}

func (f defaultFactory) ServicesToRegister(registry authproviders.Registry) []pkgGRPC.APIService {
	servicesToRegister := []pkgGRPC.APIService{
		alertService.Singleton(),
		apiTokenService.Singleton(),
		authService.New(),
		authproviderService.New(registry),
		backupRestoreService.Singleton(),
		backupService.Singleton(),
		certgen.ServiceSingleton(),
		clusterService.Singleton(),
		complianceManagerService.Singleton(),
		complianceService.Singleton(),
		configService.Singleton(),
		credentialExpiryService.Singleton(),
		debugService.Singleton(),
		deploymentService.Singleton(),
		detectionService.Singleton(),
		featureFlagService.Singleton(),
		groupService.Singleton(),
		imageService.Singleton(),
		iiService.Singleton(),
		licenseService.New(false, licenseSingletons.ManagerSingleton()),
		metadataService.New(f.restartingFlag, licenseSingletons.ManagerSingleton()),
		namespaceService.Singleton(),
		networkFlowService.Singleton(),
		networkPolicyService.Singleton(),
		nodeService.Singleton(),
		notifierService.Singleton(),
		pingService.Singleton(),
		podService.Singleton(),
		policyService.Singleton(),
		probeUploadService.Singleton(),
		processIndicatorService.Singleton(),
		processWhitelistService.Singleton(),
		rbacService.Singleton(),
		roleService.Singleton(),
		sacService.Singleton(),
		searchService.Singleton(),
		secretService.Singleton(),
		sensorService.New(connection.ManagerSingleton(), all.Singleton(), clusterDataStore.Singleton()),
		sensorUpgradeControlService.Singleton(),
		sensorUpgradeService.Singleton(),
		serviceAccountService.Singleton(),
		siService.Singleton(),
		summaryService.Singleton(),
		telemetryService.Singleton(),
		userService.Singleton(),
		cveService.Singleton(),
		integrationHealthService.Singleton(),
	}

	if features.SensorInstallationExperience.Enabled() {
		servicesToRegister = append(servicesToRegister, clusterInitService.Singleton())
	}

	autoTriggerUpgrades := sensorUpgradeConfigStore.Singleton().AutoTriggerSetting()
	if err := connection.ManagerSingleton().Start(
		clusterDataStore.Singleton(),
		networkEntityDataStore.Singleton(),
		policyDataStore.Singleton(),
		processWhitelistDataStore.Singleton(),
		autoTriggerUpgrades,
	); err != nil {
		log.Panicf("Couldn't start sensor connection manager: %v", err)
	}

	m := fetcher.SingletonManager()
	if env.OfflineModeEnv.Setting() != "true" {
		go m.Start()
	}

	if devbuild.IsEnabled() {
		servicesToRegister = append(servicesToRegister, developmentService.Singleton())
	}

	return servicesToRegister
}

func startMainServer(restartingFlag *concurrency.Flag) {
	factory := defaultFactory{
		restartingFlag: restartingFlag,
	}
	startGRPCServer(factory)
}

func watchdog(signal *concurrency.Signal, timeout time.Duration) {
	if !concurrency.WaitWithTimeout(signal, timeout) {
		log.Errorf("API server failed to start within %v!", timeout)
		log.Error("This usually means something is *very* wrong. Terminating ...")
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGABRT); err != nil {
			panic(err)
		}
	}
}

func startGRPCServer(factory serviceFactory) {
	// Temporarily elevate permissions to modify auth providers.
	authProviderRegisteringCtx := sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.AuthProvider)))

	// Create the registry of applied auth providers.
	registry, err := authproviders.NewStoreBackedRegistry(
		ssoURLPathPrefix, tokenRedirectURLPath,
		authProviderDS.Singleton(), jwt.IssuerFactorySingleton(),
		mapper.FactorySingleton())
	if err != nil {
		log.Panicf("Could not create auth provider registry: %v", err)
	}

	for typeName, factoryCreator := range authProviderBackendFactories {
		if err := registry.RegisterBackendFactory(authProviderRegisteringCtx, typeName, factoryCreator); err != nil {
			log.Panicf("Could not register %s auth provider factory: %v", typeName, err)
		}
	}
	if err := registry.Init(); err != nil {
		log.Panicf("Could not initialize auth provider registry: %v", err)
	}

	basicAuthMgr := userpass.CreateManager()

	basicAuthProvider := userpass.RegisterAuthProviderOrPanic(authProviderRegisteringCtx, basicAuthMgr, registry)

	serviceMTLSExtractor, err := service.NewExtractor()
	if err != nil {
		log.Panicf("Could not create mTLS-based service identity extractor: %v", err)
	}

	serviceTokenExtractor, err := servicecerttoken.NewExtractor(maxServiceCertTokenLeeway)
	if err != nil {
		log.Panicf("Could not create ServiceCert token-based identity extractor: %v", err)
	}

	idExtractors := []authn.IdentityExtractor{
		serviceMTLSExtractor, // internal services
		tokenbased.NewExtractor(roleDataStore.Singleton(), jwt.ValidatorSingleton()), // JWT tokens
		userpass.IdentityExtractorOrPanic(basicAuthMgr, basicAuthProvider),
		serviceTokenExtractor,
		authnUserpki.NewExtractor(tlsconfig.ManagerInstance()),
	}

	endpointCfgs, err := endpoints.InstantiateAll(tlsconfig.ManagerInstance())
	if err != nil {
		log.Panicf("Could not instantiate endpoint configs: %v", err)
	}

	config := pkgGRPC.Config{
		CustomRoutes:       factory.CustomRoutes(),
		IdentityExtractors: idExtractors,
		AuthProviders:      registry,
		Auditor:            audit.New(processor.Singleton()),
		GRPCMetrics:        metrics.GRPCSingleton(),
		HTTPMetrics:        metrics.HTTPSingleton(),
		Endpoints:          endpointCfgs,
	}

	// This helps validate that SAC is being used correctly.
	if devbuild.IsEnabled() {
		config.UnaryInterceptors = append(config.UnaryInterceptors, transitional.VerifySACScopeChecksInterceptor)
	}

	// The below enrichers handle SAC being off or on.
	// Before authorization is checked, we want to inject the sac client into the context.
	config.PreAuthContextEnrichers = append(config.PreAuthContextEnrichers,
		centralSAC.GetEnricher().PreAuthContextEnricher,
	)
	// After auth checks are run, we want to use the client (if available) to add scope checking.
	config.PostAuthContextEnrichers = append(config.PostAuthContextEnrichers,
		centralSAC.GetEnricher().PostAuthContextEnricher,
	)

	server := pkgGRPC.NewAPI(config)
	server.Register(factory.ServicesToRegister(registry)...)

	factory.StartServices()
	startedSig := server.Start()

	go watchdog(startedSig, grpcServerWatchdogTimeout)
}

func registerDelayedIntegrations(integrationsInput []iiStore.DelayedIntegration) {
	integrations := make(map[int]iiStore.DelayedIntegration, len(integrationsInput))
	for k, v := range integrationsInput {
		integrations[k] = v
	}
	ds := iiDatastore.Singleton()
	for len(integrations) > 0 {
		for idx, integration := range integrations {
			_, exists, _ := ds.GetImageIntegration(imageIntegrationContext, integration.Integration.GetId())
			if exists {
				delete(integrations, idx)
				continue
			}
			ready := integration.Trigger()
			if !ready {
				continue
			}
			// add the integration first, which is more likely to fail. If it does, no big deal -- you can still try to
			// manually add it and get the error message.
			err := imageintegration.ToNotify().NotifyUpdated(integration.Integration)
			if err == nil {
				err = ds.UpdateImageIntegration(imageIntegrationContext, integration.Integration)
				if err != nil {
					// so, we added the integration to the set but we weren't able to save it.
					// This is ok -- the image scanner will "work" and after a restart we'll try to save it again.
					log.Errorf("We added the %q integration, but saving it failed with: %v. We'll try again next restart", integration.Integration.GetName(), err)
				} else {
					log.Infof("Registered integration %q", integration.Integration.GetName())
				}
				reprocessor.Singleton().ShortCircuit()
			} else {
				log.Errorf("Unable to register integration %q: %v", integration.Integration.GetName(), err)
			}
			// either way, time to stop watching this entry
			delete(integrations, idx)
		}
		time.Sleep(5 * time.Second)
	}
	log.Debug("All dynamic integrations registered, exiting")
}

func uiRoute() routes.CustomRoute {
	return routes.CustomRoute{
		Route:         "/",
		Authorizer:    allow.Anonymous(),
		ServerHandler: ui.Mux(),
		Compression:   true,
	}
}

func (defaultFactory) CustomRoutes() (customRoutes []routes.CustomRoute) {
	customRoutes = []routes.CustomRoute{
		uiRoute(),
		{
			Route:         "/api/extensions/clusters/zip",
			Authorizer:    or.SensorOrAuthorizer(user.With(permissions.View(resources.Cluster), permissions.View(resources.ServiceIdentity))),
			ServerHandler: clustersZip.Handler(clusterDataStore.Singleton(), siStore.Singleton()),
			Compression:   false,
		},
		{
			Route:         "/api/extensions/scanner/zip",
			Authorizer:    user.With(permissions.View(resources.ScannerBundle)),
			ServerHandler: scanner.Handler(),
			Compression:   false,
		},
		{
			Route:         "/api/cli/download/",
			Authorizer:    user.With(),
			ServerHandler: cli.Handler(),
			Compression:   true,
		},
		{
			Route:         "/db/backup",
			Authorizer:    dbAuthz.DBReadAccessAuthorizer(),
			ServerHandler: globaldbHandlers.BackupDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB(), false),
			Compression:   true,
		},
		{
			Route:         "/db/backup/full",
			Authorizer:    user.WithRole(role.Admin),
			ServerHandler: globaldbHandlers.BackupDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB(), true),
			Compression:   true,
		},
		{
			Route:         "/db/restore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: globaldbHandlers.RestoreDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB()),
		},
		{
			Route:         "/api/docs/swagger",
			Authorizer:    user.With(permissions.View(resources.APIToken)),
			ServerHandler: docs.Swagger(),
			Compression:   true,
		},
		{
			Route:         "/api/graphql",
			Authorizer:    user.With(), // graphql enforces permissions internally
			ServerHandler: graphqlHandler.Handler(),
			Compression:   true,
		},
		{
			Route:         "/api/compliance/export/csv",
			Authorizer:    user.With(permissions.View(resources.Compliance)),
			ServerHandler: complianceHandlers.CSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/risk/timeline/export/csv",
			Authorizer:    user.With(permissions.View(resources.Deployment), permissions.View(resources.Indicator), permissions.View(resources.ProcessWhitelist)),
			ServerHandler: timeline.CSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/vm/export/csv",
			Authorizer:    user.With(permissions.View(resources.Image), permissions.View(resources.Deployment)),
			ServerHandler: csv.CVECSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/splunk/ta/vulnmgmt",
			Authorizer:    user.With(permissions.View(resources.Image), permissions.View(resources.Deployment)),
			ServerHandler: splunk.NewVulnMgmtHandler(deploymentDatastore.Singleton(), imageDatastore.Singleton()),
			Compression:   true,
		},
		{
			Route:         "/api/splunk/ta/compliance",
			Authorizer:    user.With(permissions.View(resources.Compliance)),
			ServerHandler: splunk.NewComplianceHandler(complianceDatastore.Singleton()),
			Compression:   true,
		},
		{
			Route:         "/db/v2/restore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: backupRestoreService.Singleton().RestoreHandler(),
		},
		{
			Route:         "/db/v2/resumerestore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: backupRestoreService.Singleton().ResumeRestoreHandler(),
		},
	}

	if features.SensorInstallationExperience.Enabled() {
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/api/extensions/clusters/helm-config.yaml",
			Authorizer:    or.SensorOrAuthorizer(user.With(permissions.View(resources.Cluster))),
			ServerHandler: clustersHelmConfig.Handler(clusterDataStore.Singleton()),
			Compression:   true,
		})
	}

	logImbueRoute := "/api/logimbue"
	customRoutes = append(customRoutes,
		routes.CustomRoute{
			Route: logImbueRoute,
			Authorizer: perrpc.FromMap(map[authz.Authorizer][]string{
				user.With(permissions.View(resources.ImbuedLogs)): {
					routes.RPCNameForHTTP(logImbueRoute, http.MethodGet),
				},
				user.With(permissions.Modify(resources.ImbuedLogs)): {
					routes.RPCNameForHTTP(logImbueRoute, http.MethodPost),
				},
			}),
			ServerHandler: logimbueHandler.Singleton(),
			Compression:   false,
		},
	)

	scannerDefinitionsRoute := "/api/extensions/scannerdefinitions"
	customRoutes = append(customRoutes,
		routes.CustomRoute{
			Route: scannerDefinitionsRoute,
			Authorizer: perrpc.FromMap(map[authz.Authorizer][]string{
				or.ScannerOr(user.With(permissions.View(resources.ScannerDefinitions))): {
					routes.RPCNameForHTTP(scannerDefinitionsRoute, http.MethodGet),
				},
				user.With(permissions.Modify(resources.ScannerDefinitions)): {
					routes.RPCNameForHTTP(scannerDefinitionsRoute, http.MethodPost),
				},
			}),
			ServerHandler: scannerDefinitionsHandler.Singleton(),
			Compression:   false,
		},
	)

	helmClusterAddRoute := "/api/helm/cluster/add"
	customRoutes = append(customRoutes,
		routes.CustomRoute{
			Route: helmClusterAddRoute,
			Authorizer: user.With(permissions.Modify(resources.Cluster),
				permissions.Modify(resources.ServiceIdentity)),
			ServerHandler: helmHandler.Handler(siStore.Singleton(), clusterService.Singleton()),
			Compression:   false,
		},
	)

	customRoutes = append(customRoutes, debugRoutes()...)
	return
}

func debugRoutes() []routes.CustomRoute {
	customRoutes := make([]routes.CustomRoute, 0, len(routes.DebugRoutes))

	for r, h := range routes.DebugRoutes {
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         r,
			Authorizer:    user.WithRole(role.Admin),
			ServerHandler: h,
			Compression:   true,
		})
	}
	return customRoutes
}

func waitForTerminationSignal() {
	signalsC := make(chan os.Signal, 1)
	signal.Notify(signalsC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-signalsC
	log.Infof("Caught %s signal", sig)
	reprocessor.Singleton().Stop()
	log.Info("Stopped reprocessor loop")
	suppress.Singleton().Stop()
	log.Info("Stopped cve unsuppress loop")
	pruning.Singleton().Stop()
	log.Info("Stopped garbage collector")
	if features.NetworkGraphExternalSrcs.Enabled() {
		gatherer.Singleton().Stop()
		log.Info("Stopped network graph default external sources gatherer")
	}

	globaldb.Close()

	if sig == syscall.SIGHUP {
		log.Info("Restarting central")
		osutils.Restart()
	}
	log.Info("Central terminated")
}
