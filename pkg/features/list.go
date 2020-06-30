package features

var (
	// AnalystNotesUI enables the Analyst Notes UI.
	// NB: When removing this feature flag, remove references in ui/src/utils/featureFlags.js
	AnalystNotesUI = registerFeature("Enable Analyst Notes UI", "ROX_ANALYST_NOTES_UI", true)

	// EventTimelineClusteredEventsUI enables the Event Timeline UI for Clustered Events.
	// NB: When removing this feature flag, remove references in ui/src/utils/featureFlags.js
	EventTimelineClusteredEventsUI = registerFeature("Enable Event Timeline Clustered Events UI", "ROX_EVENT_TIMELINE_CLUSTERED_EVENTS_UI", false)

	// ImageLabelPolicy enables the Required Image Label policy type
	ImageLabelPolicy = registerFeature("Enable the Required Image Label Policy", "ROX_REQUIRED_IMAGE_LABEL_POLICY", true)

	// AdmissionControlService enables running admission control as a separate microservice.
	AdmissionControlService = registerFeature("Separate admission control microservice", "ROX_ADMISSION_CONTROL_SERVICE", true)

	// AdmissionControlEnforceOnUpdate enables support for having the admission controller enforce on updates.
	AdmissionControlEnforceOnUpdate = registerFeature("Allow admission controller to enforce on update", "ROX_ADMISSION_CONTROL_ENFORCE_ON_UPDATE", true)

	// DryRunPolicyJobMechanism enables submitting dry run of a policy as a job, and querying the status using job id.
	DryRunPolicyJobMechanism = registerFeature("Dry run policy job mechanism", "ROX_DRY_RUN_JOB", true)

	// BooleanPolicyLogic enables support for an extended policy logic
	BooleanPolicyLogic = registerFeature("Enable Boolean Policy Logic", "ROX_BOOLEAN_POLICY_LOGIC", true)

	// PolicyImportExport feature flag enables policy import and export
	PolicyImportExport = registerFeature("Enable Import/Export for Analyst Workflow", "ROX_POLICY_IMPORT_EXPORT", true)

	// AuthTestMode feature flag allows test mode flow for new auth provider in UI
	AuthTestMode = registerFeature("Enable Auth Test Mode UI", "ROX_AUTH_TEST_MODE_UI", true)

	// CurrentUserInfo enables showing information about the current user in UI
	CurrentUserInfo = registerFeature("Enable Current User Info UI", "ROX_CURRENT_USER_INFO", true)

	// ComplianceInNodes enables running of node-related Compliance checks in the compliance pods
	ComplianceInNodes = registerFeature("Enable compliance checks in nodes", "ROX_COMPLIANCE_IN_NODES", false)
)
