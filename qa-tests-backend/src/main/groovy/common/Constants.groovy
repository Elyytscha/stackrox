package common

class Constants {
    static final ORCHESTRATOR_NAMESPACE = "qa"
    static final SCHEDULES_SUPPORTED = false
    static final CHECK_CVES_IN_COMPLIANCE = false
    static final RUN_FLAKEY_TESTS = false
    static final Map<String, String> CSV_COLUMN_MAPPING = [
            "Standard" : "standard",
            "Cluster" : "cluster",
            "Namespace" : "namespace",
            "Object Type" : "objectType",
            "Object Name" : "objectName",
            "Control" : "control",
            "Control Description" : "controlDescription",
            "State" : "state",
            "Evidence" : "evidence",
            "Assessment Time" : "timestamp",
    ]
    static final VIOLATIONS_WHITELIST = [
            "Monitoring" : ["CVSS >= 7"],
            "clairify" : ["Red Hat Package Manager Execution"],
            "authorization-plugin" : ["Latest tag", "90-Day Image Age"],
            "webhookserver" : ["90-Day Image Age"],
    ]

    /*
        StackRox Product Feature Flags

        We need to manually maintain this list here
     */
    static final String K8SRBAC_FEATURE_FLAG = "ROX_K8S_RBAC"
    static final String CLIENT_CA_AUTH_FEATURE_FLAG = "ROX_CLIENT_CA_AUTH"
}
