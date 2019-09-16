package kubernetes

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/checks/msgfmt"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/generated/internalapi/compliance"
	"gopkg.in/yaml.v2"
	"k8s.io/apiserver/pkg/server/options/encryptionconfig"
)

const kubeAPIProcessName = "kube-apiserver"

const tlsCiphers = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256," +
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256," +
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305," +
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384," +
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305," +
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384," +
	"TLS_RSA_WITH_AES_256_GCM_SHA384," +
	"TLS_RSA_WITH_AES_128_GCM_SHA256"

func init() {
	framework.MustRegisterChecks(
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_1", "anonymous-auth", "false", "true", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_2", "basic-auth-file", "", "", common.Unset),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_3", "insecure-allow-any-token", "", "", common.Unset),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_4", "kubelet-https", "true", "true", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_5", "insecure-bind-address", "127.0.0.1", "127.0.0.1", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_6", "insecure-port", "0", "8080", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_7", "secure-port", "0", "6443", common.NotMatches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_8", "profiling", "false", "true", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_9", "repair-malformed-updates", "false", "true", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_10", "enable-admission-plugins", "AlwaysAdmit", "AlwaysAdmit", common.NotContains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_11", "enable-admission-plugins", "AlwaysPullImages", "AlwaysAdmit", common.Contains),
		common.PerNodeDeprecatedCheck("CIS_Kubernetes_v1_4_1:1_1_12", "The 'DenyEscalatingExec' admission control policy has been deprecated in Kubernetes 1.13"),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_13", "enable-admission-plugins", "SecurityContextDeny", "AlwaysAdmit", common.Contains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_14", "enable-admission-plugins", "NamespaceLifecycle", "AlwaysAdmit", common.Contains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_15", "audit-log-path", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_16", "audit-log-maxage", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_17", "audit-log-maxbackup", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_18", "audit-log-maxsize", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_19", "authorization-mode", "AlwaysAllow", "AlwaysAllow", common.NotContains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_20", "token-auth-file", "", "", common.Unset),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_21", "kubelet-certificate-authority", "", "", common.Set),
		multipleFlagsSetCheck("CIS_Kubernetes_v1_4_1:1_1_22", "kube-apiserver", "kubelet-client-certificate", "kubelet-client-key"),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_23", "service-account-lookup", "true", "false", common.Matches),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_24", "enable-admission-plugins", "PodSecurityPolicy", "AlwaysAdmit", common.Contains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_25", "service-account-key-file", "", "", common.Set),
		multipleFlagsSetCheck("CIS_Kubernetes_v1_4_1:1_1_26", "kube-apiserver", "etcd-certfile", "etcd-keyfile"),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_27", "enable-admission-plugins", "ServiceAccount", "AlwaysAdmit", common.Contains),
		multipleFlagsSetCheck("CIS_Kubernetes_v1_4_1:1_1_28", "kube-apiserver", "tls-cert-file", "tls-private-key-file"),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_29", "client-ca-file", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_30", "etcd-cafile", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_31", "tls-cipher-suites", tlsCiphers, "", common.OnlyContains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_32", "authorization-mode", "Node", "AlwaysAllow", common.Contains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_33", "enable-admission-plugins", "NodeRestriction", "AlwaysAllow", common.Contains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_34", "encryption-provider-config", "", "", common.Set),
		encryptionProvider(),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_36", "enable-admission-plugins", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_37", "feature-gates", "AdvancedAuditing=false", "", common.NotContains),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_38", "request-timeout", "", "", common.Set),
		masterAPIServerCommandLine("CIS_Kubernetes_v1_4_1:1_1_39", "authorization-mode", "RBAC", "AlwaysAllow", common.Set),
	)
}

func masterAPIServerCommandLine(name string, key, target, defaultVal string, evalFunc common.CommandEvaluationFunc) framework.Check {
	return genericKubernetesCommandlineCheck(name, kubeAPIProcessName, key, target, defaultVal, evalFunc)
}

func encryptionProvider() framework.Check {
	md := framework.CheckMetadata{
		ID:                 "CIS_Kubernetes_v1_4_1:1_1_35",
		Scope:              framework.NodeKind,
		InterpretationText: "StackRox checks that the Kubernetes API server uses the `aescbc` encryption provider",
		DataDependencies:   []string{"HostScraped"},
	}
	return framework.NewCheckFromFunc(md, common.PerNodeCheck(
		func(ctx framework.ComplianceContext, ret *compliance.ComplianceReturn) {
			process, exists := common.GetProcess(ret, kubeAPIProcessName)
			if !exists {
				framework.NoteNowf(ctx, "Process %q not found on host therefore check is not applicable", kubeAPIProcessName)
			}
			arg := common.GetArgForFlag(process.Args, "experimental-encryption-provider-config")
			if arg == nil {
				framework.FailNowf(ctx, "experimental-encryption-provider-config is not set, which means that aescbc is not in use")
			} else if arg.GetFile() == nil {
				framework.FailNowf(ctx, "No file was found experimental-encryption-provider-config value of %q", msgfmt.FormatStrings(arg.GetValues()...))
			}

			var config encryptionconfig.EncryptionConfig
			if err := yaml.Unmarshal(arg.GetFile().GetContent(), &config); err != nil {
				framework.FailNowf(ctx, "Could not parse file %q to check for aescbc specification due to %v. Please manually check", arg.GetFile().GetPath(), err)
			}
			if config.Kind != "EncryptionConfig" {
				framework.FailNowf(ctx, "Incorrect configuration kind %q in file %q", config.Kind, arg.GetFile().GetPath())
				return
			}
			for _, resource := range config.Resources {
				for _, provider := range resource.Providers {

					if provider.AESCBC != nil {
						framework.PassNow(ctx, "Provider is set as aescbc")
					}
				}
			}
			framework.Fail(ctx, "Provider is not set as aescbc")
		}))
}
