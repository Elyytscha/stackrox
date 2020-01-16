package env

var (
	// InitialTelemetryEnabledEnv indicates whether StackRox was installed with telemetry enabled.  This flag is
	// overridden by the telemetry configuration in the database  Defaults to false here and true in the install process
	// so that it will default to on for new installations and off for old installations
	InitialTelemetryEnabledEnv = registerBooleanSetting("ROX_INIT_TELEMETRY_ENABLED", false)

	// TelemetryEndpoint is the endpoint to which to send telemetry data.
	TelemetryEndpoint = RegisterSetting("ROX_TELEMETRY_ENDPOINT", AllowEmpty())
)
