import static Services.waitForViolation
import static Services.waitForSuspiciousProcessInRiskIndicators

import io.stackrox.proto.storage.RiskOuterClass
import io.stackrox.proto.api.v1.AlertServiceOuterClass
import io.stackrox.proto.storage.AlertOuterClass
import services.AlertService
import services.ClusterService

import groups.BAT

import io.stackrox.proto.storage.ProcessWhitelistOuterClass
import objects.Deployment

import org.apache.commons.lang.StringUtils

import org.junit.experimental.categories.Category

import services.ProcessWhitelistService
import spock.lang.Shared
import spock.lang.Unroll

class ProcessWhiteListsTest extends BaseSpecification {
    @Shared
    private String clusterId

    static final private String DEPLOYMENTNGINX = "pw-deploymentnginx"
    static final private String DEPLOYMENTNGINX_RESOLVE_VIOLATION = "pw-deploymentnginx-violation-resolve"
    static final private String DEPLOYMENTNGINX_RESOLVE_AND_WHITELIST_VIOLATION =
            "pw-deploymentnginx-violation-resolve-whitelist"
    static final private String DEPLOYMENTNGINX_SOFTLOCK = "pw-deploymentnginx-softlock"
    static final private String DEPLOYMENTNGINX_DELETE = "pw-deploymentnginx-delete"

    static final private String DEPLOYMENTNGINX_REMOVEPROCESS = "pw-deploymentnginx-removeprocess"
    static final private List<Deployment> DEPLOYMENTS =
            [
                    new Deployment()
                     .setName(DEPLOYMENTNGINX)
                     .setImage("nginx:1.7.9")
                     .addPort(22, "TCP")
                     .addAnnotation("test", "annotation")
                     .setEnv(["CLUSTER_NAME": "main"])
                     .addLabel("app", "test"),
             new Deployment()
                     .setName(DEPLOYMENTNGINX_RESOLVE_VIOLATION)
                     .setImage("nginx:1.7.9")
                     .addPort(22, "TCP")
                     .addAnnotation("test", "annotation")
                     .setEnv(["CLUSTER_NAME": "main"])
                     .addLabel("app", "test"),
             new Deployment()
                     .setName(DEPLOYMENTNGINX_RESOLVE_AND_WHITELIST_VIOLATION)
                     .setImage("nginx:1.7.9")
                     .addPort(22, "TCP")
                     .addAnnotation("test", "annotation")
                     .setEnv(["CLUSTER_NAME": "main"])
                     .addLabel("app", "test"),
             new Deployment()
                     .setName(DEPLOYMENTNGINX_SOFTLOCK)
                     .setImage("nginx:1.7.9")
                     .addPort(22, "TCP")
                     .addAnnotation("test", "annotation")
                     .setEnv(["CLUSTER_NAME": "main"])
                     .addLabel("app", "test"),
             new Deployment()
                     .setName(DEPLOYMENTNGINX_DELETE)
                     .setImage("nginx:1.7.9")
                     .addPort(22, "TCP")
                     .addAnnotation("test", "annotation")
                     .setEnv(["CLUSTER_NAME": "main"])
                     .addLabel("app", "test"),
                    new Deployment()
                          .setName(DEPLOYMENTNGINX_REMOVEPROCESS)
                          .setImage("nginx:1.7.9")
                          .addPort(22, "TCP")
                          .addAnnotation("test", "annotation")
                          .setEnv(["CLUSTER_NAME": "main"])
                          .addLabel("app", "test"),
            ]

    def setupSpec() {
        clusterId = ClusterService.getClusterId()

        orchestrator.batchCreateDeployments(DEPLOYMENTS)
        for (Deployment deployment : DEPLOYMENTS) {
            assert Services.waitForDeployment(deployment)
        }
    }

    def cleanupSpec() {
        for (Deployment deployment : DEPLOYMENTS) {
            orchestrator.deleteDeployment(deployment)
        }

        //need to  delete whitelists for the container deployed after each test
    }

    @Unroll
    @Category(BAT)
    def "Verify processes risk indicators for the given key after soft-lock on #deploymentName"() {
        when:
        "exec into the container and run a process and wait for soft lock to kick in"
        def deployment = DEPLOYMENTS.find { it.name == deploymentName }
        assert deployment != null
        String deploymentId = deployment.getDeploymentUid()
        String containerName = deployment.getName()
        ProcessWhitelistOuterClass.ProcessWhitelist whitelist = ProcessWhitelistService.
                    getProcessWhitelist(clusterId, deployment, containerName)
        assert (whitelist != null)
        assert ((whitelist.key.deploymentId.equalsIgnoreCase(deploymentId)) &&
                    (whitelist.key.containerName.equalsIgnoreCase(containerName)))
        assert whitelist.elementsList.find { it.element.processName == processName } != null
        // Check that startup processes are not impacted
        Thread.sleep(10000)
        orchestrator.execInContainer(deployment, "ls")
        Thread.sleep(50000)
        orchestrator.execInContainer(deployment, "pwd")

        then:
        "verify for suspicious process in risk indicator"
        RiskOuterClass.Risk.Result result = waitForSuspiciousProcessInRiskIndicators(deploymentId, 60)
        assert (result != null)
        // Check that ls doesn't exist as a risky process
        RiskOuterClass.Risk.Result.Factor lsFactor =  result.factorsList.find { it.message.contains("ls") }
        assert lsFactor == null
        // Check that pwd is a risky process
        RiskOuterClass.Risk.Result.Factor pwdFactor =  result.factorsList.find { it.message.contains("pwd") }
        assert pwdFactor != null
        where:
        "Data inputs are :"
        deploymentName                      |   processName
        DEPLOYMENTNGINX_SOFTLOCK            |   "/usr/sbin/nginx"
    }

    /* TODO(ROX-3108)
    @Unroll
    @Category(BAT)
    def "Verify  whitelist processes for the given key before and after locking "() {
        when:
        def deployment = DEPLOYMENTS.find { it.name == deploymentName }
        assert deployment != null
        String deploymentId = deployment.getDeploymentUid()
        // Currently, we always create a deployment where the container name is the same
        // as the deployment name
        String containerName = deployment.getName()
        "get process whitelists is called for a key"
        ProcessWhitelistOuterClass.ProcessWhitelist whitelist = ProcessWhitelistService.
                getProcessWhitelist(clusterId, deployment, containerName)

        assert (whitelist != null)

        then:
        "Verify  whitelisted processes for a given key before and after calling lock whitelists"
        assert ((whitelist.key.deploymentId.equalsIgnoreCase(deploymentId)) &&
                    (whitelist.key.containerName.equalsIgnoreCase(containerName)))
        assert  whitelist.getElements(0).element.processName.contains(processName)

        //lock the whitelist with the key of the container just deployed
        List<ProcessWhitelistOuterClass.ProcessWhitelist> lockProcessWhitelists = ProcessWhitelistService.
                lockProcessWhitelists(clusterId, deployment, containerName, true)
        assert  lockProcessWhitelists.size() == 1
        assert  lockProcessWhitelists.get(0).getElementsList().
            find { it.element.processName.equalsIgnoreCase(processName) } != null

        where:
        "Data inputs are :"
        deploymentName     | processName

        DEPLOYMENTNGINX    | "/usr/sbin/nginx"
    }
    */

    @Unroll
    @Category(BAT)
    def "Verify whitelist process violation after resolve whitelist on #deploymentName"() {
               /*
                    a)Lock the whitelists for the key
                    b)exec into the container and run a process
                    c)verify violation alert for Unauthorized Process Execution
                    d)
                        test case :choose to only resolve violation
                            exec into the container and run the  process again and verify violation alert
                        test case : choose to both resolve and whitelist
                            exec into the container and run the  process again and verify no violation alert
               */
        when:
        "exec into the container after locking whitelists and create a whitelist violation"
        def deployment = DEPLOYMENTS.find { it.name == deploymentName }
        assert deployment != null
        String deploymentId = deployment.getDeploymentUid()
        String containerName = deployment.getName()
        ProcessWhitelistOuterClass.ProcessWhitelist whitelist = ProcessWhitelistService.
                 getProcessWhitelist(clusterId, deployment, containerName)
        assert (whitelist != null)
        assert ((whitelist.key.deploymentId.equalsIgnoreCase(deploymentId)) &&
                 (whitelist.key.containerName.equalsIgnoreCase(containerName)))
        assert whitelist.elementsList.find { it.element.processName == processName } != null

        List<ProcessWhitelistOuterClass.ProcessWhitelist> lockProcessWhitelists = ProcessWhitelistService.
                 lockProcessWhitelists(clusterId, deployment, containerName, true)
        assert (!StringUtils.isEmpty(lockProcessWhitelists.get(0).getElements(0).getElement().processName))
        orchestrator.execInContainer(deployment, "pwd")

        //check for whitelist  violation
        assert waitForViolation(containerName, "Unauthorized Process Execution", 120)
        List<AlertOuterClass.ListAlert> alertList = AlertService.getViolations(AlertServiceOuterClass.ListAlertsRequest
                 .newBuilder().build())
        String alertId
        for (AlertOuterClass.ListAlert alert : alertList) {
            if (alert.getPolicy().name.equalsIgnoreCase("Unauthorized Process Execution") &&
                     alert.deployment.id.equalsIgnoreCase(deploymentId)) {
                alertId = alert.id
                AlertService.resolveAlert(alertId, resolveWhitelist)
            }
         }
        orchestrator.execInContainer(deployment, "pwd")
        if (resolveWhitelist) {
            waitForViolation(containerName, "Unauthorized Process Execution", 90)
        }
        else {
            assert waitForViolation(containerName, "Unauthorized Process Execution", 90)
        }
        then:
        "Verify for violation or no violation after resolve/resolve and whitelist"
        List<AlertOuterClass.ListAlert> alertListAnother = AlertService
                 .getViolations(AlertServiceOuterClass.ListAlertsRequest
                 .newBuilder().build())
        int numAlertsAfterResolve
        for (AlertOuterClass.ListAlert alert : alertListAnother) {
            if (alert.getPolicy().name.equalsIgnoreCase("Unauthorized Process Execution")
                     && alert.deployment.id.equalsIgnoreCase(deploymentId)) {
                numAlertsAfterResolve++
                AlertService.resolveAlert(alert.id, false)
                break
             }
         }
        System.out.println("numAlertsAfterResolve .. " + numAlertsAfterResolve)
        assert (numAlertsAfterResolve  == expectedViolationsCount)

        where:
        "Data inputs are :"
        deploymentName                                   | processName  | resolveWhitelist | expectedViolationsCount

        DEPLOYMENTNGINX_RESOLVE_VIOLATION               | "/usr/sbin/nginx"      | false            | 1

        DEPLOYMENTNGINX_RESOLVE_AND_WHITELIST_VIOLATION | "/usr/sbin/nginx"      | true             | 0
     }

    @Category(BAT)
    def "Verify whitelists are deleted when their deployment is deleted"() {
        /*
                a)get all whitelists
                b)verify whitelists exist for a deployment
                c)delete the deployment
                d)get all whitelists
                e)verify all whitelists for the deployment have been deleted
        */
        when:
        "a deployment is deleted"
        //Get all whitelists for our deployment and assert they exist
        def deployment = DEPLOYMENTS.find { it.name == DEPLOYMENTNGINX_DELETE }
        assert deployment != null
        String containerName = deployment.getName()
        def whitelistsCreated = ProcessWhitelistService.
                waitForDeploymentWhitelistsCreated(clusterId, deployment, containerName)
        assert(whitelistsCreated)

        //Delete the deployment
        orchestrator.deleteDeployment(deployment)
        Services.waitForSRDeletion(deployment)

        then:
        "Verify that all whitelists with that deployment ID have been deleted"
        def whitelistsDeleted = ProcessWhitelistService.
                waitForDeploymentWhitelistsDeleted(clusterId, deployment, containerName)
        assert(whitelistsDeleted)
    }

    @Unroll
    @Category(BAT)
    def "Verify removed whitelist process not getting added back to whitelist after rerun on #deploymentName"() {
        /*
                1.run a process and verify if it exists in the whitelist
                2.remove the process
                3.rerun the process to verify it it does not get added to the whitelist
         */
        when:
        "an added process is removed and whitelist is locked and the process is run"
        def deployment = DEPLOYMENTS.find { it.name == deploymentName }
        assert deployment != null
        def deploymentId = deployment.deploymentUid
        def containerName = deploymentName
        def namespace = deployment.getNamespace()

        //Wait for whitelist to be created
        def initialWhitelist = ProcessWhitelistService.
                getProcessWhitelist(clusterId, deployment, containerName)
        assert (initialWhitelist != null)

        //Add the process to the whitelist
        ProcessWhitelistOuterClass.ProcessWhitelistKey [] keys = [
                new ProcessWhitelistOuterClass
                .ProcessWhitelistKey().newBuilderForType().setContainerName(containerName)
                .setDeploymentId(deploymentId).setClusterId(clusterId).setNamespace(namespace).build(),
        ]
        String [] toBeAddedProcesses = ["pwd"]
        String [] toBeRemovedProcesses = []
        List<ProcessWhitelistOuterClass.ProcessWhitelist> updatedList = ProcessWhitelistService
                .updateProcessWhitelists(keys, toBeAddedProcesses, toBeRemovedProcesses)
        assert ( updatedList!= null)
        ProcessWhitelistOuterClass.ProcessWhitelist whitelist = ProcessWhitelistService.
                getProcessWhitelist(clusterId, deployment, containerName)
        List<ProcessWhitelistOuterClass.WhitelistElement> elements = whitelist.elementsList
        ProcessWhitelistOuterClass.WhitelistElement element = elements.find { it.element.processName.contains("pwd") }
        assert ( element != null)

        //Remove the process from the whitelist
        toBeAddedProcesses = []
        toBeRemovedProcesses = ["pwd"]
        List<ProcessWhitelistOuterClass.ProcessWhitelist> updatedListAfterRemoveProcess = ProcessWhitelistService
                .updateProcessWhitelists(keys, toBeAddedProcesses, toBeRemovedProcesses)
        assert ( updatedListAfterRemoveProcess!= null)
        orchestrator.execInContainer(deployment, "pwd")
        then:
        "verify process is not added to the whitelist"
        ProcessWhitelistOuterClass.ProcessWhitelist whitelistAfterReRun = ProcessWhitelistService.
                getProcessWhitelist(clusterId, deployment, containerName)
        assert  ( whitelistAfterReRun.elementsList.find { it.element.processName.contains("pwd") } == null)
        where:
        deploymentName                                   | processName
        DEPLOYMENTNGINX_REMOVEPROCESS           |   "nginx"
    }

    }
