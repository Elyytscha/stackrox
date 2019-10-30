import groovy.json.JsonSlurper
import io.stackrox.proto.storage.PolicyOuterClass
import io.stackrox.proto.storage.NotifierOuterClass

import groups.BAT
import groups.Integration
import org.junit.experimental.categories.Category
import services.NotifierService
import spock.lang.Unroll
import objects.Deployment

class IntegrationsTest extends BaseSpecification {

    static final private String BUSYBOX = "genericbusybox"

    private static final CA_CERT = '''-----BEGIN CERTIFICATE-----
MIIDgDCCAmgCCQDYOU2KIlcBQjANBgkqhkiG9w0BAQsFADCBgTELMAkGA1UEBhMC
VVMxCzAJBgNVBAgMAkNBMQswCQYDVQQHDAJTRjERMA8GA1UECgwIc3RhY2tyb3gx
HzAdBgNVBAMMFndlYmhvb2tzZXJ2ZXIuc3RhY2tyb3gxJDAiBgkqhkiG9w0BCQEW
FXN0YWNrcm94QHN0YWNrcm94LmNvbTAeFw0xOTAzMjMxNTQzMjVaFw0yOTAzMjAx
NTQzMjVaMIGBMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExCzAJBgNVBAcMAlNG
MREwDwYDVQQKDAhzdGFja3JveDEfMB0GA1UEAwwWd2ViaG9va3NlcnZlci5zdGFj
a3JveDEkMCIGCSqGSIb3DQEJARYVc3RhY2tyb3hAc3RhY2tyb3guY29tMIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuPzgVGykTALNHDljDiCjwI4ZfF2r
lGKWdtvUhurh42Cl2Kfn0Vgy7mYRjdK/uOiSIl6LVXuNw7w4yg48dXm8By+I3+hs
vMH4ixykWxPn6Ez3Utuuwggn/yAs4kE2Wj0ztFMpRHBGL7Qi7oEv+Vo4349ZJg16
a55db45O3LgOED119F1hQvxblNZhcA2hnNOhveXsJLfdOQKz6UA4KtdBFXxEeZuB
fC45wCHw6kjRrBEPYKB4py4ywYMdUHqswBDn6B3LtwvrrJVPTySK4sgZmOTF2XGg
JRm52MS0rYEvBpEtgkPdknoIv0VnxihMUuRhMXHfGOTFyhWuf/nF2aihXwIDAQAB
MA0GCSqGSIb3DQEBCwUAA4IBAQCYT7jo6Durhx+liRgYNO3G3mRyc36syVNGsllU
Mf5wOUHjxWfppWHzldxMeZRKksrg7xfMXdcGaOOZgD8Ir/pPK2HP48g6KIDWCiVO
kh9AGCLY9osxkBqAihtvJWNkEda+wA9ggF/7wx+0Ci+b/1NvXHeNU3uO3rP7Npwc
rxhvyNqv7MwqpMN6V8hFxqM/3ny8aoUedFsYsEvm8Dm1VLyBiIqZk0CA2oj3NIjb
ObOdSTZUQI4TZOXOpJCpa97CnqroNi7RrT05JOfoe/DPmhoJmF4AUrnd/YUb8pgF
/jvC1xBvPVtJFbYeBVysQCrRk+f/NyyUejQv+OCJ+B1KtJh4
-----END CERTIFICATE-----'''

    @Unroll
    @Category([BAT])
    def "Verify Email Integration (port #port, disable TLS=#disableTLS, startTLS=#startTLS)"() {
        given:
        "a configuration that is expected to work"
        NotifierOuterClass.Notifier notifier = Services.addEmailNotifier(
                "mailgun",
                disableTLS,
                startTLS,
                port
        )

        when:
        "the integration is tested"
        Boolean response = Services.testNotifier(notifier)

        then:
        "the API should return an empty message or an error, depending on the config"
        assert response == shouldSucceed

        cleanup:
        "remove notifier"
        if (notifier != null) {
            Services.deleteNotifier(notifier.id)
        }

        where:
        "data"

        port | disableTLS | startTLS | shouldSucceed

        // Port 465 tests
        // This port speaks TLS from the start.
        // (Also test null, since 465 is the default.)
        /////////////////
        // Speaking TLS should work
        465  | false      | false    | true
        null | false      | false    | true
        // Sending STARTTLS is not expected to work when already using TLS
        465  | false      | true     | false
        null | false      | true     | false
        // Speaking non-TLS to a TLS port should fail and not time out, regardless of STARTTLS (see ROX-366)
        465  | true       | false    | false
        465  | true       | true     | false
        null | true       | false    | false
        null | true       | true     | false

        // Port 587 tests
        // At MailGun, this port begins unencrypted and supports STARTTLS.
        /////////////////
        // Starting unencrypted and _not_ using STARTTLS should work
        587  | true       | false    | true
        // Starting unencrypted and using STARTTLS should work
        587  | true       | true     | true
        // Speaking TLS to a non-TLS port should fail whether you use STARTTLS or not.
        587  | false      | false    | false
        587  | false      | true     | false

        // Cannot add port 25 tests since GCP blocks outgoing
        // connections to port 25
    }

    @Category(BAT)
    def "Verify Splunk Integration"() {
        when:
        "the integration is tested"

        Deployment  deployment =
            new Deployment()
                .setNamespace("stackrox")
                .setName("splunk")
                .setImage("store/splunk/enterprise:latest")
                .addPort (8000)
                .addPort (8088)
                .addAnnotation("test", "annotation")
                .setEnv([ "SPLUNK_START_ARGS": "--accept-license", "SPLUNK_USER": "root" ])
                .addLabel("app", "splunk")
                .setPrivilegedFlag(true)
                .addVolume("test", "/tmp")
                .setSkipReplicaWait(true)
                .addImagePullSecret("stackrox")

        orchestrator.createDeployment(deployment)

        Deployment serviceDeployment =
            new Deployment()
                .addLabel("app", "splunk")
                .setCreateLoadBalancer(true)
                .setNamespace("stackrox")
                .setName("splunk")
                .setTargetPort(8000)
                .addPort(8000, "TCP")
                .setServiceName("splunk-http")

        orchestrator.createService(serviceDeployment)

        Deployment serviceDeploymentHec =
            new Deployment()
                .addLabel("app" , "splunk")
                .setCreateLoadBalancer(true)
                .setNamespace("stackrox")
                .setName("splunk")
                .setTargetPort(8088)
                .addPort(8088, "TCP")
                .setServiceName("splunk-hec")

        orchestrator.createService(serviceDeploymentHec)

        then :
        "the API should return an empty message or an error, depending on the config"

        cleanup:
        "remove Deployment and services"
        if (deployment != null) {
            orchestrator.deleteDeployment(deployment)
        }
        orchestrator.deleteService( "splunk-hec", "stackrox")
        orchestrator.deleteService( "splunk-http", "stackrox")
    }

    @Category(Integration)
    def "Verify PagerDuty Integration"() {
        when:
        "Add PagerDuty integration and test it"
        NotifierOuterClass.Notifier notifier = NotifierService.addPagerDutyNotifier("pdTest")
        Boolean response = NotifierService.testNotifier(notifier)
        assert response

        and:
        "Get current the incidents' number from the PagerDuty"
        int preNum = NotifierService.getFirstPagerDutyIncident().incidents[0].incident_number

        and:
        "Binding the notifier with the policy"
        def policy = Services.getPolicyByName("Latest tag")
        def updatedPolicy = PolicyOuterClass.Policy.newBuilder(policy).addNotifiers(notifier.getId()).build()
        Services.updatePolicy(updatedPolicy)

        and:
        "Create a new deployment to trigger the policy"
        Deployment  deployment =
                new Deployment()
                        .setName ("pgtest")
                        .setImage ("nginx:latest")
                        .addPort (22)
                        .addLabel ("app", "test")
        orchestrator.createDeployment(deployment)
        assert Services.waitForViolation("pgtest", "Latest tag", 30)

        and:
        "Get current the first incident from the PagerDuty"
        def firIncident = NotifierService.waitForPagerDutyUpdate(preNum)

        then:
        "Verify a new incident is triggered and it contains the latest tag alert information"
        assert firIncident != null
        assert firIncident.incidents[0].description.contains("Alert on deployments with images using tag 'latest'")

        cleanup:
        "remove Deployment and service"
        if (deployment != null) {
            orchestrator.deleteDeployment(deployment)
        }
        if (notifier != null) {
            NotifierService.deleteNotifier(notifier.getId())
        }
    }

    @Unroll
    @Category(BAT)
    def "Verify Generic Integration Test Endpoint (#tlsOptsDesc, audit=#audit)"() {
        when:
        "the integration is tested"

        NotifierOuterClass.Notifier notifier = Services.getWebhookIntegrationConfiguration(
                enableTLS, caCert, skipTLSVerification, auditLoggingEnabled)

        then :
        "the API should return an empty message or an error, depending on the config"
        assert shouldSucceed == Services.testNotifier(notifier)

        where:
        "data"

        enableTLS | caCert | skipTLSVerification | auditLoggingEnabled | shouldSucceed | tlsOptsDesc

        false | ""         | false               | false | true | "no TLS"
        true  | ""         | true                | false | true | "TLS, no verify"
        true  | CA_CERT    | false               | false | true | "TLS, verify custom CA"
        true  | ""         | false               | false | false | "TLS, verify system CA"
        false | ""         | false               | true | true | "no TLS"
        true  | ""         | true                | true | true | "TLS, no verify"
        true  | CA_CERT    | false               | true | true | "TLS, verify custom CA"
        true  | ""         | false               | true | false | "TLS, verify system CA"
    }

    @Category(BAT)
    def "Verify Generic Integration Values With Audit Off"() {
        when:
        "the integration is created"
        NotifierOuterClass.Notifier notifier = Services.getWebhookIntegrationConfiguration(
                false, "", false, false)
        String notifierId = Services.addNotifier(notifier)

        def policy = Services.getPolicyByName("Latest tag")
        def updatedPolicy = PolicyOuterClass.Policy.newBuilder(policy).addNotifiers(notifierId).build()
        Services.updatePolicy(updatedPolicy)

        Deployment  deployment =
                new Deployment()
                        .setName(BUSYBOX)
                        .setImage("busybox")
                        .setCommand(["sleep", "8000"])

        orchestrator.createDeployment(deployment)

        then:
        "We should check to make sure we got a value"
        assert Services.waitForViolation(BUSYBOX, "Latest tag", 30)

        def get = new URL("http://localhost:8080").openConnection()
        def jsonSlurper = new JsonSlurper()
        def object = jsonSlurper.parseText(get.getInputStream().getText())
        def generic = object[-1]

        assert generic["headers"]["Headerkey"] == ["headervalue"]
        assert generic["headers"]["Content-Type"] == ["application/json"]
        assert generic["headers"]["Authorization"] == ["Basic YWRtaW46YWRtaW4="]
        assert generic["data"]["fieldkey"] == "fieldvalue"
        assert generic["data"]["alert"]["policy"]["name"] == "Latest tag"
        assert generic["data"]["alert"]["deployment"]["name"] == BUSYBOX

        cleanup:
        if (notifier != null) {
            Services.deleteNotifier(notifierId)
        }
        if (deployment != null) {
            orchestrator.deleteDeployment(deployment)
        }
    }

    @Category(BAT)
    def "Verify Generic Integration Values With Audit On"() {
        when:
        "the integration is created"
        NotifierOuterClass.Notifier notifier = Services.getWebhookIntegrationConfiguration(
                false, "", false, true)
        String notifierId = Services.addNotifier(notifier)

        def policy = Services.getPolicyByName("Latest tag")
        def updatedPolicy = PolicyOuterClass.Policy.newBuilder(policy).addNotifiers(notifierId).build()
        Services.updatePolicy(updatedPolicy)

        Deployment  deployment =
                new Deployment()
                        .setName(BUSYBOX)
                        .setImage("busybox")
                        .setCommand(["sleep", "8000"])

        orchestrator.createDeployment(deployment)

        then:
        "We should check to make sure we got a value"
        assert Services.waitForViolation(BUSYBOX, "Latest tag", 30)

        def get = new URL("http://localhost:8080").openConnection()
        def jsonSlurper = new JsonSlurper()
        def object = jsonSlurper.parseText(get.getInputStream().getText())

        for (def generic : object) {
            if (generic["data"]["audit"] == null) {
                continue
            }
            if (generic["data"]["audit"]["policy"] == null) {
                continue
            }

            assert generic["headers"]["Headerkey"] == ["headervalue"]
            assert generic["headers"]["Content-Type"] == ["application/json"]
            assert generic["headers"]["Authorization"] == ["Basic YWRtaW46YWRtaW4="]
            assert generic["data"]["fieldkey"] == "fieldvalue"
            assert generic["data"]["audit"]["policy"]["name"] == "Latest tag"
            assert generic["data"]["audit"]["deployment"]["name"] == BUSYBOX
        }

        cleanup:
        if (notifier != null) {
            Services.deleteNotifier(notifierId)
        }
        if (deployment != null) {
            orchestrator.deleteDeployment(deployment)
        }
    }
}
