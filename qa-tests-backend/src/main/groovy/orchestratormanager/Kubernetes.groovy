package orchestratormanager

import common.YamlGenerator
import io.fabric8.kubernetes.api.model.Capabilities
import io.fabric8.kubernetes.api.model.ConfigMap
import io.fabric8.kubernetes.api.model.ConfigMapEnvSource
import io.fabric8.kubernetes.api.model.ConfigMapKeySelectorBuilder
import io.fabric8.kubernetes.api.model.Container
import io.fabric8.kubernetes.api.model.ContainerPort
import io.fabric8.kubernetes.api.model.ContainerStatus
import io.fabric8.kubernetes.api.model.EnvFromSource
import io.fabric8.kubernetes.api.model.EnvVar
import io.fabric8.kubernetes.api.model.EnvVarBuilder
import io.fabric8.kubernetes.api.model.EnvVarSourceBuilder
import io.fabric8.kubernetes.api.model.HostPathVolumeSource
import io.fabric8.kubernetes.api.model.IntOrString
import io.fabric8.kubernetes.api.model.LabelSelector
import io.fabric8.kubernetes.api.model.LocalObjectReference
import io.fabric8.kubernetes.api.model.Namespace
import io.fabric8.kubernetes.api.model.ObjectMeta
import io.fabric8.kubernetes.api.model.ObjectFieldSelectorBuilder
import io.fabric8.kubernetes.api.model.Pod
import io.fabric8.kubernetes.api.model.PodList
import io.fabric8.kubernetes.api.model.PodSpec
import io.fabric8.kubernetes.api.model.PodTemplateSpec
import io.fabric8.kubernetes.api.model.Quantity
import io.fabric8.kubernetes.api.model.ResourceFieldSelectorBuilder
import io.fabric8.kubernetes.api.model.ResourceRequirements
import io.fabric8.kubernetes.api.model.Secret
import io.fabric8.kubernetes.api.model.SecretEnvSource
import io.fabric8.kubernetes.api.model.SecretKeySelectorBuilder
import io.fabric8.kubernetes.api.model.SecretVolumeSource
import io.fabric8.kubernetes.api.model.SecurityContext
import io.fabric8.kubernetes.api.model.Service
import io.fabric8.kubernetes.api.model.ServiceAccount
import io.fabric8.kubernetes.api.model.ServicePort
import io.fabric8.kubernetes.api.model.ServiceSpec
import io.fabric8.kubernetes.api.model.Volume
import io.fabric8.kubernetes.api.model.VolumeMount
import io.fabric8.kubernetes.api.model.apps.Deployment as K8sDeployment
import io.fabric8.kubernetes.api.model.apps.DaemonSet as K8sDaemonSet
import io.fabric8.kubernetes.api.model.batch.Job as K8sJob
import io.fabric8.kubernetes.api.model.apps.DaemonSetList
import io.fabric8.kubernetes.api.model.apps.DaemonSetSpec
import io.fabric8.kubernetes.api.model.apps.DeploymentList
import io.fabric8.kubernetes.api.model.apps.DeploymentSpec
import io.fabric8.kubernetes.api.model.apps.DoneableDaemonSet
import io.fabric8.kubernetes.api.model.apps.DoneableDeployment
import io.fabric8.kubernetes.api.model.batch.DoneableJob
import io.fabric8.kubernetes.api.model.apps.DoneableStatefulSet
import io.fabric8.kubernetes.api.model.apps.StatefulSetList
import io.fabric8.kubernetes.api.model.apps.StatefulSet as K8sStatefulSet
import io.fabric8.kubernetes.api.model.batch.JobList
import io.fabric8.kubernetes.api.model.batch.JobSpec
import io.fabric8.kubernetes.api.model.policy.HostPortRange
import io.fabric8.kubernetes.api.model.policy.PodSecurityPolicy
import io.fabric8.kubernetes.api.model.policy.PodSecurityPolicyBuilder
import io.fabric8.kubernetes.api.model.networking.NetworkPolicyBuilder
import io.fabric8.kubernetes.api.model.networking.NetworkPolicyEgressRuleBuilder
import io.fabric8.kubernetes.api.model.networking.NetworkPolicyIngressRuleBuilder
import io.fabric8.kubernetes.api.model.networking.NetworkPolicyPeerBuilder
import io.fabric8.kubernetes.api.model.rbac.ClusterRole
import io.fabric8.kubernetes.api.model.rbac.ClusterRoleBinding
import io.fabric8.kubernetes.api.model.rbac.PolicyRule
import io.fabric8.kubernetes.api.model.rbac.Role
import io.fabric8.kubernetes.api.model.rbac.RoleBinding
import io.fabric8.kubernetes.api.model.rbac.RoleRef
import io.fabric8.kubernetes.api.model.rbac.Subject
import io.fabric8.kubernetes.client.Callback
import io.fabric8.kubernetes.client.DefaultKubernetesClient
import io.fabric8.kubernetes.client.KubernetesClient
import io.fabric8.kubernetes.client.KubernetesClientException
import io.fabric8.kubernetes.client.dsl.Deletable
import io.fabric8.kubernetes.client.dsl.ExecListener
import io.fabric8.kubernetes.client.dsl.ExecWatch
import io.fabric8.kubernetes.client.dsl.MixedOperation
import io.fabric8.kubernetes.client.dsl.Resource
import io.fabric8.kubernetes.client.dsl.ScalableResource
import io.fabric8.kubernetes.client.utils.BlockingInputStreamPumper
import io.kubernetes.client.models.V1beta1ValidatingWebhookConfiguration
import objects.ConfigMapKeyRef
import objects.DaemonSet
import objects.Deployment
import objects.Job
import objects.K8sPolicyRule
import objects.K8sRole
import objects.K8sRoleBinding
import objects.K8sServiceAccount
import objects.K8sSubject
import objects.NetworkPolicy
import objects.NetworkPolicyTypes
import objects.Node
import objects.SecretKeyRef
import okhttp3.Response
import util.Timer

import java.util.concurrent.CountDownLatch
import java.util.concurrent.Executors
import java.util.concurrent.Future
import java.util.concurrent.ScheduledExecutorService
import java.util.concurrent.TimeUnit
import java.util.stream.Collectors

class Kubernetes implements OrchestratorMain {

    final int sleepDurationSeconds = 5
    final int maxWaitTimeSeconds = 90
    final int lbWaitTimeSeconds = 600
    final int intervalTime = 1
    String namespace
    KubernetesClient client

    MixedOperation<K8sDaemonSet, DaemonSetList, DoneableDaemonSet, Resource<K8sDaemonSet, DoneableDaemonSet>> daemonsets

    MixedOperation<K8sDeployment, DeploymentList, DoneableDeployment,
            ScalableResource<K8sDeployment, DoneableDeployment>> deployments

    MixedOperation<K8sStatefulSet, StatefulSetList, DoneableStatefulSet,
            ScalableResource<K8sStatefulSet, DoneableStatefulSet>> statefulsets

    MixedOperation<K8sJob, JobList, DoneableJob,
            ScalableResource<K8sJob, DoneableJob>> jobs

    Kubernetes(String ns) {
        this.namespace = ns
        this.client = new DefaultKubernetesClient()
        // On OpenShift, the namespace config is typically non-null (set to the default project), which causes all
        // "any namespace" requests to be scoped to the default project.
        this.client.configuration.namespace = null
        this.client.configuration.setRollingTimeout(60 * 60 * 1000)
        this.deployments = this.client.apps().deployments()
        this.daemonsets = this.client.apps().daemonSets()
        this.statefulsets = this.client.apps().statefulSets()
        this.jobs = this.client.batch().jobs()
    }

    Kubernetes() {
        Kubernetes("default")
    }

    def ensureNamespaceExists(String ns) {
        Namespace namespace = new Namespace("v1", null, new ObjectMeta(name: ns), null, null)
        try {
            client.namespaces().create(namespace)
            defaultPspForNamespace(ns)
            println "Created namespace ${ns}"
        } catch (KubernetesClientException kce) {
            // 409 is already exists
            if (kce.code != 409) {
                throw kce
            }
        }
    }

    def setup() {
        ensureNamespaceExists(this.namespace)
    }

    def cleanup() {
    }

    /*
        Deployment Methods
    */

    def createDeployment(Deployment deployment) {
        ensureNamespaceExists(deployment.namespace)
        createDeploymentNoWait(deployment)
        waitForDeploymentAndPopulateInfo(deployment)
    }

    boolean updateDeploymentNoWait(Deployment deployment) {
        if (deployments.inNamespace(deployment.namespace).withName(deployment.name).get()) {
            println "Deployment ${deployment.name} found in namespace ${deployment.namespace}. Updating..."
        } else {
            println "Deployment ${deployment.name} NOT found in namespace ${deployment.namespace}. Creating..."
        }
        return createDeploymentNoWait(deployment)
    }

    def updateDeployment(Deployment deployment) {
        if (deployments.inNamespace(deployment.namespace).withName(deployment.name).get()) {
            println "Deployment ${deployment.name} found in namespace ${deployment.namespace}. Updating..."
        } else {
            println "Deployment ${deployment.name} NOT found in namespace ${deployment.namespace}. Creating..."
        }
        // Our createDeployment actually uses createOrReplace so it should work for these purposes
        return createDeployment(deployment)
    }

    def batchCreateDeployments(List<Deployment> deployments) {
        for (Deployment deployment : deployments) {
            ensureNamespaceExists(deployment.namespace)
            createDeploymentNoWait(deployment)
        }
        for (Deployment deployment : deployments) {
            waitForDeploymentAndPopulateInfo(deployment)
        }
    }

    def waitForAllPodsToBeRemoved(String ns, Map<String, String>labels, int retries = 30, int intervalSeconds = 5) {
        LabelSelector selector = new LabelSelector()
        selector.matchLabels = labels
        Timer t = new Timer(retries, intervalSeconds)
        PodList list
        while (t.IsValid()) {
            list = client.pods().inNamespace(ns).withLabelSelector(selector).list()
            if (list.items.size() == 0) {
                return true
            }
        }

        println "Timed out waiting for the following pods to be removed"
        for (Pod pod : list.getItems()) {
            println "\t- ${pod.metadata.name}"
        }
        return false
    }

    List<Pod> getPods(String namespace, String appName) {
        return getPodsByLabel(namespace, ["app": appName])
    }

    List<Pod> getPodsByLabel(String namespace, Map<String, String> label) {
        def selector = new LabelSelector()
        selector.matchLabels = label
        PodList list = evaluateWithRetry(2, 3) {
            return client.pods().inNamespace(namespace).withLabelSelector(selector).list()
        }
        return list.getItems()
    }

    Boolean deletePod(String namespace, String podName, Long gracePeriodSecs) {
        Deletable<Boolean> podClient = client.pods().inNamespace(namespace).withName(podName)
        if (gracePeriodSecs != null) {
            podClient = podClient.withGracePeriod(gracePeriodSecs)
        }
        return podClient.delete()
    }

    def getAndPrintPods(String ns, String name) {
        println "Status of ${name}'s pods:"
        for (Pod pod : getPodsByLabel(ns, ["deployment": name])) {
            println "\t- ${pod.metadata.name}"
            for (ContainerStatus status : pod.status.containerStatuses) {
                println "\t  Container status: ${status.state}"
            }
        }
    }

    def waitForDeploymentDeletion(Deployment deploy, int retries = 30, int intervalSeconds = 5) {
        Timer t = new Timer(retries, intervalSeconds)

        K8sDeployment d
        while (t.IsValid()) {
            d = this.deployments.inNamespace(deploy.namespace).withName(deploy.name).get()
            if (d == null) {
                println "${deploy.name}: deployment removed."
                return
            }
            getAndPrintPods(deploy.namespace, deploy.name)
        }
        println "Timed out waiting for deployment ${deploy.name} to be deleted"
    }

    def deleteAndWaitForDeploymentDeletion(Deployment... deployments) {
        for (Deployment deployment : deployments) {
            this.deleteDeployment(deployment)
        }
        for (Deployment deployment : deployments) {
            this.waitForDeploymentDeletion(deployment)
        }
    }

    def deleteDeployment(Deployment deployment) {
        if (deployment.exposeAsService) {
            this.deleteService(deployment.name, deployment.namespace)
        }
        // Retry deletion due to race condition in sdk and controller
        // See https://github.com/fabric8io/kubernetes-client/issues/1477
        Timer t = new Timer(10, 1)
        while (t.IsValid()) {
            try {
                this.deployments.inNamespace(deployment.namespace).withName(deployment.name).delete()
                break
            } catch (KubernetesClientException ex) {
                println "Failed to delete deployment: ${ex.toString()}"
            }
        }
        println "removing the deployment: ${deployment.name}"
    }

    def createOrchestratorDeployment(K8sDeployment dep) {
        dep.setApiVersion("")
        dep.metadata.setResourceVersion("")
        return this.deployments.inNamespace(dep.metadata.namespace).create(dep)
    }

    K8sDeployment getOrchestratorDeployment(String ns, String name) {
        return this.deployments.inNamespace(ns).withName(name).get()
    }

    String getDeploymentId(Deployment deployment) {
        return this.deployments.inNamespace(deployment.namespace)
                .withName(deployment.name)
                .get()?.metadata?.uid
    }

    def getDeploymentReplicaCount(Deployment deployment) {
        K8sDeployment d = this.deployments.inNamespace(deployment.namespace)
                .withName(deployment.name)
                .get()
        if (d != null) {
            println "${deployment.name}: Replicas=${d.getSpec().getReplicas()}"
            return d.getSpec().getReplicas()
        }
    }

    def getDeploymentUnavailableReplicaCount(Deployment deployment) {
        K8sDeployment d = this.deployments
                .inNamespace(deployment.namespace)
                .withName(deployment.name)
                .get()
        if (d != null) {
            println "${deployment.name}: Unavailable Replicas=${d.getStatus().getUnavailableReplicas()}"
            return d.getStatus().getUnavailableReplicas()
        }
    }

    def getDeploymentNodeSelectors(Deployment deployment) {
        K8sDeployment d = this.deployments
                .inNamespace(deployment.namespace)
                .withName(deployment.name)
                .get()
        if (d != null) {
            println "${deployment.name}: Host=${d.getSpec().getTemplate().getSpec().getNodeSelector()}"
            return d.getSpec().getTemplate().getSpec().getNodeSelector()
        }
    }

    Set<String> getDeploymentSecrets(Deployment deployment) {
        K8sDeployment d = this.deployments
                .inNamespace(deployment.namespace)
                .withName(deployment.name)
                .get()
        Set<String> secretSet = [] as Set
        if (d != null) {
            d.getSpec()?.getTemplate()?.getSpec()?.getVolumes()?.each { secretSet.add(it.secret.secretName) }
            d.getSpec()?.getTemplate()?.getSpec()?.getContainers()?.getAt(0)?.getEnvFrom()?.each {
                // Only care about secrets for now.
                if (it.getSecretRef() != null) {
                    secretSet.add(it.secretRef.name)
                }
            }
        }
        return secretSet
    }

    def getDeploymentCount(String ns = null) {
        return this.deployments.inNamespace(ns).list().getItems().collect { it.metadata.name }
    }

    def createPortForward(int port, Deployment deployment, String podName = "") {
        if (deployment.pods.size() == 0) {
            throw new KubernetesClientException(
                    "Error creating port-forward: Could not get pod details from deployment.")
        }
        if (deployment.pods.size() > 1 && podName.equals("")) {
            throw new KubernetesClientException(
                    "Error creating port-forward: Deployment contains more than 1 pod, but no pod was specified.")
        }
        return deployment.pods.size() == 1 ?
                this.client.pods()
                        .inNamespace(deployment.namespace)
                        .withName(deployment.pods.get(0).name)
                        .portForward(port) :
                this.client.pods()
                        .inNamespace(deployment.namespace)
                        .withName(podName)
                        .portForward(port)
    }

    /*
        DaemonSet Methods
    */

    def createDaemonSet(DaemonSet daemonSet) {
        ensureNamespaceExists(daemonSet.namespace)
        createDaemonSetNoWait(daemonSet)
        waitForDaemonSetAndPopulateInfo(daemonSet)
    }

    def deleteDaemonSet(DaemonSet daemonSet) {
        this.daemonsets.inNamespace(daemonSet.namespace).withName(daemonSet.name).delete()
        println "${daemonSet.name}: daemonset removed."
    }

    def createJob(Job job) {
        ensureNamespaceExists(job.namespace)

        job.getNamespace() != null ?: job.setNamespace(this.namespace)

        K8sJob k8sJob = new K8sJob(
                metadata: new ObjectMeta(
                        name: job.name,
                        namespace: job.namespace,
                        labels: job.labels
                ),
                spec: new JobSpec(
                        template: new PodTemplateSpec(
                                metadata: new ObjectMeta(
                                        name: job.name,
                                        namespace: job.namespace,
                                        labels: job.labels
                                ),
                                spec: generatePodSpec(job)
                        )
                )
        )
        // Jobs cannot be Always
        k8sJob.spec.template.spec.restartPolicy = "Never"

        try {
            println("Told the orchestrator to create job " + job.getName())
            return this.jobs.inNamespace(job.namespace).createOrReplace(k8sJob)
        } catch (Exception e) {
            println("Error creating k8s job" + e.toString())
        }
        return null
    }

    def deleteJob(Job job) {
        this.jobs.inNamespace(job.namespace).withName(job.name).delete()
        println "${job.name}: job removed."
    }

    def waitForDaemonSetDeletion(String name, String ns = namespace) {
        Timer t = new Timer(30, 5)

        while (t.IsValid()) {
            if (this.daemonsets.inNamespace(ns).withName(name).get() == null) {
                println "Daemonset ${name} has been deleted"
                return
            }
        }
        println "Timed out waiting for daemonset ${name} to stop"
    }

    def getDaemonSetReplicaCount(DaemonSet daemonSet) {
        K8sDaemonSet d = this.daemonsets
                .inNamespace(daemonSet.namespace)
                .withName(daemonSet.name)
                .get()
        if (d != null) {
            println "${daemonSet.name}: Replicas=${d.getStatus().getDesiredNumberScheduled()}"
            return d.getStatus().getDesiredNumberScheduled()
        }
        return null
    }

    def getDaemonSetUnavailableReplicaCount(DaemonSet daemonSet) {
        K8sDaemonSet d = this.daemonsets
                .inNamespace(daemonSet.namespace)
                .withName(daemonSet.name)
                .get()
        if (d != null) {
            println "${daemonSet.name}: Unavailable Replicas=${d.getStatus().getNumberUnavailable()}"
            return d.getStatus().getNumberUnavailable() == null ? 0 : d.getStatus().getNumberUnavailable()
        }
        return null
    }

    def getDaemonSetNodeSelectors(DaemonSet daemonSet) {
        K8sDaemonSet d = this.daemonsets
                .inNamespace(daemonSet.namespace)
                .withName(daemonSet.name)
                .get()
        if (d != null) {
            println "${daemonSet.name}: Host=${d.getSpec().getTemplate().getSpec().getNodeSelector()}"
            return d.getSpec().getTemplate().getSpec().getNodeSelector()
        }
        return null
    }

    def getDaemonSetCount(String ns = null) {
        return this.daemonsets.inNamespace(ns).list().getItems().collect { it.metadata.name }
    }

    String getDaemonSetId(DaemonSet daemonSet) {
        return this.daemonsets.inNamespace(daemonSet.namespace)
                .withName(daemonSet.name)
                .get()?.metadata?.uid
    }

    /*
        StatefulSet Methods
    */

    def getStatefulSetCount(String ns = null) {
        return this.statefulsets.inNamespace(ns).list().getItems().collect { it.metadata.name }
    }

    /*
        Container Methods
    */

    def deleteContainer(String containerName, String namespace = this.namespace) {
        withRetry(2, 3) {
            client.pods().inNamespace(namespace).withName(containerName).delete()
        }
    }

    def wasContainerKilled(String containerName, String namespace = this.namespace) {
        Timer t = new Timer(20, 3)

        Pod pod
        while (t.IsValid()) {
            try {
                pod = client.pods().inNamespace(namespace).withName(containerName).get()
                if (pod == null) {
                    println "Could not query K8S for pod details, assuming pod was killed"
                    return true
                }
                println "Pod Deletion Timestamp: ${pod.metadata.deletionTimestamp}"
                if (pod.metadata.deletionTimestamp != null ) {
                    return true
                }
            } catch (Exception e) {
                println "wasContainerKilled: error fetching pod details - retrying"
            }
        }
        println "wasContainerKilled: did not determine container was killed before 60s timeout"
        println "container details were found:\n${containerName}: ${pod.toString()}"
        return false
    }

    def isKubeProxyPresent() {
        return evaluateWithRetry(2, 3) {
            PodList pods = client.pods().inAnyNamespace().list()
            return pods.getItems().findAll {
                it.getSpec().getContainers().find {
                    it.getImage().contains("kube-proxy")
                }
            }
        }
    }

    def isKubeDashboardRunning() {
        return evaluateWithRetry(2, 3) {
            PodList pods = client.pods().inAnyNamespace().list()
            List<Pod> kubeDashboards = pods.getItems().findAll {
                it.getSpec().getContainers().find {
                    it.getImage().contains("kubernetes-dashboard")
                }
            }
            return kubeDashboards.size() > 0
        }
    }

    def getContainerlogs(Deployment deployment) {
        withRetry(2, 3) {
            PodList pods = client.pods().inNamespace(deployment.namespace).list()
            Pod pod = pods.getItems().find { it.getMetadata().getName().startsWith(deployment.name) }

            try {
                println client.pods()
                        .inNamespace(pod.metadata.namespace)
                        .withName(pod.metadata.name)
                        .tailingLines(5000)
                        .watchLog(System.out)
            } catch (Exception e) {
                println "Error getting container logs: ${e.toString()}"
            }
        }
    }

    def getStaticPodCount(String ns = null) {
        return evaluateWithRetry(2, 3) {
            // This method assumes that a static pod name will contain the node name that the pod is running on
            def nodeNames = client.nodes().list().items.collect { it.metadata.name }
            Set<String> staticPods = [] as Set
            client.pods().inNamespace(ns).list().items.each {
                for (String node : nodeNames) {
                    if (it.metadata.name.contains(node)) {
                        staticPods.add(it.metadata.name[0..it.metadata.name.indexOf(node) - 2])
                    }
                }
            }
            return staticPods
        }
    }

    /*
        Service Methods
    */

    def createService(Deployment deployment) {
        withRetry(2, 3) {
            Service service = new Service(
                    metadata: new ObjectMeta(
                            name: deployment.serviceName ? deployment.serviceName : deployment.name,
                            namespace: deployment.namespace,
                            labels: deployment.labels
                    ),
                    spec: new ServiceSpec(
                            ports: deployment.getPorts().collect {
                                k, v ->
                                    new ServicePort(
                                            name: k as String,
                                            port: k as Integer,
                                            protocol: v,
                                            targetPort: new IntOrString(deployment.targetport) ?:
                                                    new IntOrString(k as Integer)
                                    )
                            },
                            selector: deployment.labels,
                            type: deployment.createLoadBalancer ? "LoadBalancer" : "ClusterIP"
                    )
            )
            client.services().inNamespace(deployment.namespace).createOrReplace(service)
            println(deployment.serviceName ?: deployment.name + " service created")
            if (deployment.createLoadBalancer) {
                deployment.loadBalancerIP = waitForLoadBalancer(deployment.serviceName ?:
                        deployment.name, deployment.namespace)
            }
        }
    }

    def createService(objects.Service s) {
        withRetry(2, 3) {
            Service service = new Service(
                    metadata: new ObjectMeta(
                            name: s.name,
                            namespace: s.namespace,
                            labels: s.labels
                    ),
                    spec: new ServiceSpec(
                            ports: s.getPorts().collect {
                                k, v ->
                                    new ServicePort(
                                            name: k as String,
                                            port: k as Integer,
                                            protocol: v,
                                            targetPort:
                                                    new IntOrString(s.targetport) ?: new IntOrString(k as Integer)
                                    )
                            },
                            selector: s.labels,
                            type: s.type.toString()
                    )
            )
            client.services().inNamespace(s.namespace).createOrReplace(service)
        }
        println "${s.name}: Service created"
        if (objects.Service.Type.LOADBALANCER == s.type) {
            s.loadBalancerIP = waitForLoadBalancer(s.name, s.namespace)
        }
    }

    def deleteService(String name, String namespace = this.namespace) {
        withRetry(2, 3) {
            println "${name}: Service deleting..."
            client.services().inNamespace(namespace).withName(name).delete()
        }
        println "${name}: Service deleted"
    }

    def waitForServiceDeletion(objects.Service service) {
        boolean beenDeleted = false

        int retries = maxWaitTimeSeconds / sleepDurationSeconds
        Timer t = new Timer(retries, sleepDurationSeconds)
        while (!beenDeleted && t.IsValid()) {
            Service s = client.services().inNamespace(service.namespace).withName(service.name).get()
            beenDeleted = true

            println "Waiting for service ${service.name} to be deleted"
            if (s != null) {
                beenDeleted = false
            }
        }

        if (beenDeleted) {
            println service.name + ": service removed."
        } else {
            println "Timed out waiting for service ${service.name} to be removed"
        }
    }

    def waitForLoadBalancer(Deployment deployment) {
        "Creating a load balancer"
        if (deployment.createLoadBalancer) {
            deployment.loadBalancerIP = waitForLoadBalancer(deployment.serviceName ?:
                                        deployment.name, deployment.namespace)
        }
    }

    /**
     * This is an overloaded method for creating load balancer for a given service or deployment
     *
     * @param service
     */
    String waitForLoadBalancer(String serviceName, String namespace) {
        Service service
        String loadBalancerIP
        int iterations = lbWaitTimeSeconds/intervalTime
        println "Waiting for LB external IP for " + serviceName
        Timer t = new Timer(iterations, intervalTime)
        while (t.IsValid()) {
            service = client.services().inNamespace(namespace).withName(serviceName).get()
            if (service?.status?.loadBalancer?.ingress?.size()) {
                loadBalancerIP = service.status.loadBalancer.ingress.get(0).
                                  ip ?: service.status.loadBalancer.ingress.get(0).hostname
                println "LB IP: " + loadBalancerIP
                break
            }
        }
        if (loadBalancerIP == null) {
            println("Could not get loadBalancer IP in ${t.SecondsSince()} and ${iterations}")
        }
        return loadBalancerIP
    }
    /*
        Secrets Methods
    */
    def waitForSecretCreation(String secretName, String namespace = this.namespace) {
        int retries = maxWaitTimeSeconds / sleepDurationSeconds
        Timer t = new Timer(retries, sleepDurationSeconds)
        while (t.IsValid()) {
            Secret secret = client.secrets().inNamespace(namespace).withName(secretName).get()
            if (secret != null) {
                println secretName + ": secret created."
                return secret
            }
        }
        println "Timed out waiting for secret ${secretName} to be created"
    }

    String createImagePullSecret(String name, String username, String password, String namespace = this.namespace) {
        return createImagePullSecret(new objects.Secret(
            name: name,
            server: "https://docker.io",
            username: username,
            password: password,
            namespace: namespace
        ))
    }

    String createImagePullSecret(objects.Secret secret) {
        if (!secret.username || !secret.password) {
            throw new RuntimeException("Secret requires a username and password: " +
                    "username provided: ${secret.username}, " +
                    "password provided: ${secret.password}")
        }
        def namespace = secret.namespace ?: this.namespace

        def auth = secret.username + ":" + secret.password
        def b64Password = Base64.getEncoder().encodeToString(auth.getBytes())
        def dockerConfigJSON =  "{\"auths\":{\"" + secret.server + "\": {\"auth\": \"" + b64Password + "\"}}}"
        Map<String, String> data = new HashMap<String, String>()
        data.put(".dockerconfigjson", Base64.getEncoder().encodeToString(dockerConfigJSON.getBytes()))

        Secret k8sSecret = new Secret(
                apiVersion: "v1",
                kind: "Secret",
                type: "kubernetes.io/dockerconfigjson",
                data: data,
                metadata: new ObjectMeta(
                        name: secret.name,
                        namespace: namespace
                )
        )

        Secret createdSecret = client.secrets().inNamespace(namespace).createOrReplace(k8sSecret)
        if (createdSecret != null) {
            createdSecret = waitForSecretCreation(secret.name, namespace)
            return createdSecret.metadata.uid
        }
        throw new RuntimeException("Couldn't create secret")
    }

    String createSecret(String name, String namespace = this.namespace) {
        return evaluateWithRetry(2, 3) {
            Map<String, String> data = new HashMap<String, String>()
            data.put("username", "YWRtaW4=")
            data.put("password", "MWYyZDFlMmU2N2Rm")

            Secret secret = new Secret(
                    apiVersion: "v1",
                    kind: "Secret",
                    type: "Opaque",
                    data: data,
                    metadata: new ObjectMeta(
                            name: name
                    )
            )

            try {
                Secret createdSecret = client.secrets().inNamespace(namespace).createOrReplace(secret)
                if (createdSecret != null) {
                    createdSecret = waitForSecretCreation(name, namespace)
                    return createdSecret.metadata.uid
                }
            } catch (Exception e) {
                println("Error creating secret" + e.toString())
            }
            return null
        }
    }

    String updateSecret(Secret secret) {
        return withRetry(2, 3) {
            client.secrets().inNamespace(secret.metadata.namespace).createOrReplace(secret)
        }
    }

    def deleteSecret(String name, String namespace = this.namespace) {
        withRetry(2, 3) {
            client.secrets().inNamespace(namespace).withName(name).delete()
        }
        sleep(sleepDurationSeconds * 1000)
        println name + ": Secret removed."
    }

    def getSecretCount(String ns = null) {
        return evaluateWithRetry(2, 3) {
            return client.secrets().inNamespace(ns).list().getItems().findAll {
                !it.type.startsWith("kubernetes.io/service-account-token")
            }.size()
        }
    }

    Secret getSecret(String name, String namespace) {
        return evaluateWithRetry(2, 3) {
            return client.secrets().inNamespace(namespace).withName(name).get()
        }
    }

    /*
        Network Policy Methods
    */

    String applyNetworkPolicy(NetworkPolicy policy) {
        return evaluateWithRetry(2, 3) {
            io.fabric8.kubernetes.api.model.networking.NetworkPolicy networkPolicy =
                    createNetworkPolicyObject(policy)

            println "${networkPolicy.metadata.name}: NetworkPolicy created:"
            println YamlGenerator.toYaml(networkPolicy)
            io.fabric8.kubernetes.api.model.networking.NetworkPolicy createdPolicy =
                    client.network().networkPolicies()
                            .inNamespace(networkPolicy.metadata.namespace ?
                                    networkPolicy.metadata.namespace :
                                    this.namespace).createOrReplace(networkPolicy)
            policy.uid = createdPolicy.metadata.uid
            return createdPolicy.metadata.uid
        }
    }

    boolean deleteNetworkPolicy(NetworkPolicy policy) {
        return evaluateWithRetry(2, 3) {
            Boolean status = client.network().networkPolicies()
                    .inNamespace(policy.namespace ? policy.namespace : this.namespace)
                    .withName(policy.name)
                    .delete()
            if (status) {
                println "${policy.name}: NetworkPolicy removed."
                return true
            }
            println "${policy.name}: Failed to remove NetworkPolicy."
            return false
        }
    }

    def getNetworkPolicyCount(String ns) {
        return evaluateWithRetry(2, 3) {
            return client.network().networkPolicies().inNamespace(ns).list().items.size()
        }
    }

    def getAllNetworkPoliciesNamesByNamespace(Boolean ignoreUndoneStackroxGenerated = false) {
        return evaluateWithRetry(2, 3) {
            Map<String, List<String>> networkPolicies = [:]
            client.network().networkPolicies().inAnyNamespace().list().items.each {
                boolean skip = false
                if (ignoreUndoneStackroxGenerated) {
                    if (it.spec.podSelector.matchLabels?.get("network-policies.stackrox.io/disable") == "nomatch") {
                        skip = true
                    }
                }
                skip ?: networkPolicies.containsKey(it.metadata.namespace) ?
                        networkPolicies.get(it.metadata.namespace).add(it.metadata.name) :
                        networkPolicies.put(it.metadata.namespace, [it.metadata.name])
            }
            return networkPolicies
        }
    }

    /*
        Node Methods
     */

    def getNodeCount() {
        return evaluateWithRetry(2, 3) {
            return client.nodes().list().getItems().size()
        }
    }

    List<Node> getNodeDetails() {
        return evaluateWithRetry(2, 3) {
            return client.nodes().list().items.collect {
                new Node(
                        uid: it.metadata.uid,
                        name: it.metadata.name,
                        labels: it.metadata.labels,
                        annotations: it.metadata.annotations,
                        internalIps: it.status.addresses.findAll { it.type == "InternalIP" }*.address,
                        externalIps: it.status.addresses.findAll { it.type == "ExternalIP" }*.address,
                        containerRuntimeVersion: it.status.nodeInfo.containerRuntimeVersion,
                        kernelVersion: it.status.nodeInfo.kernelVersion,
                        osImage: it.status.nodeInfo.osImage
                )
            }
        }
    }

    def isGKE() {
        return evaluateWithRetry(2, 3) {
            List<Node> gkeNodes = client.nodes().list().getItems().findAll {
                it.getStatus().getNodeInfo().getKubeletVersion().contains("gke")
            }
            return gkeNodes.size() > 0
        }
    }

    def supportsNetworkPolicies() {
        return evaluateWithRetry(2, 3) {
            List<Node> gkeNodes = client.nodes().list().getItems().findAll {
                it.getStatus().getNodeInfo().getKubeletVersion().contains("gke")
            }
            return gkeNodes.size() > 0
        }
    }

    /*
        Namespace Methods
     */

    List<objects.Namespace> getNamespaceDetails() {
        return evaluateWithRetry(2, 3) {
            return client.namespaces().list().items.collect {
                new objects.Namespace(
                        uid: it.metadata.uid,
                        name: it.metadata.name,
                        labels: it.metadata.labels,
                        deploymentCount: getDeploymentCount(it.metadata.name) +
                                getDaemonSetCount(it.metadata.name) +
                                getStaticPodCount(it.metadata.name) +
                                getStatefulSetCount(it.metadata.name) +
                                getJobCount(it.metadata.name),
                        secretsCount: getSecretCount(it.metadata.name),
                        networkPolicyCount: getNetworkPolicyCount(it.metadata.name)
                )
            }
        }
    }

    /*
        Service Accounts
     */

    List<ServiceAccount> getServiceAccounts() {
        return evaluateWithRetry(1, 2) {
            def serviceAccounts = []
            client.serviceAccounts().inAnyNamespace().list().items.each {
                // Ingest the K8s service account to a K8sServiceAccount() in a manner similar to the SR product.
                def annotations = it.metadata.annotations
                if (annotations) {
                    annotations.remove("kubectl.kubernetes.io/last-applied-configuration")
                }
                serviceAccounts.add(new K8sServiceAccount(
                        name: it.metadata.name,
                        namespace: it.metadata.namespace,
                        labels: it.metadata.labels ? it.metadata.labels : [:],
                        annotations: annotations ?: [:],
                        secrets: it.secrets*.name,
                        imagePullSecrets: it.imagePullSecrets*.name,
                        automountToken: it.automountServiceAccountToken == null
                                ? true : it.automountServiceAccountToken,
                ))
            }
            return serviceAccounts
        }
    }

    def createServiceAccount(K8sServiceAccount serviceAccount) {
        withRetry(1, 2) {
            ServiceAccount sa = new ServiceAccount(
                    metadata: new ObjectMeta(
                            name: serviceAccount.name,
                            namespace: serviceAccount.namespace,
                            labels: serviceAccount.labels,
                            annotations: serviceAccount.annotations
                    ),
                    secrets: serviceAccount.secrets,
                    imagePullSecrets: serviceAccount.imagePullSecrets
            )
            client.serviceAccounts().inNamespace(sa.metadata.namespace).createOrReplace(sa)
        }
    }

    def deleteServiceAccount(K8sServiceAccount serviceAccount) {
        withRetry(1, 2) {
            client.serviceAccounts().inNamespace(serviceAccount.namespace).withName(serviceAccount.name).delete()
        }
    }

    def addServiceAccountImagePullSecret(String accountName, String secretName, String namespace = this.namespace) {
        withRetry(1, 2) {
            ServiceAccount serviceAccount = client.serviceAccounts()
                    .inNamespace(namespace)
                    .withName(accountName)
                    .get()

            Set<LocalObjectReference> imagePullSecretsSet = new HashSet<>(serviceAccount.getImagePullSecrets())
            imagePullSecretsSet.add(new LocalObjectReference(secretName))
            List<LocalObjectReference> imagePullSecretsList = []
            imagePullSecretsList.addAll(imagePullSecretsSet)
            serviceAccount.setImagePullSecrets(imagePullSecretsList)

            client.serviceAccounts().inNamespace(namespace).withName(accountName).createOrReplace(serviceAccount)
        }
    }

    def removeServiceAccountImagePullSecret(String accountName, String secretName, String namespace = this.namespace) {
        ServiceAccount serviceAccount = client.serviceAccounts()
                .inNamespace(namespace)
                .withName(accountName)
                .get()

        Set<LocalObjectReference> imagePullSecretsSet = new HashSet<>(serviceAccount.getImagePullSecrets())
        imagePullSecretsSet.remove(new LocalObjectReference(secretName))
        List<LocalObjectReference> imagePullSecretsList = []
        imagePullSecretsList.addAll(imagePullSecretsSet)
        serviceAccount.setImagePullSecrets(imagePullSecretsList)

        client.serviceAccounts().inNamespace(namespace).withName(accountName).createOrReplace(serviceAccount)
    }

    /*
        Roles
     */

    List<K8sRole> getRoles() {
        return evaluateWithRetry(1, 2) {
            def roles = []
            client.rbac().roles().inAnyNamespace().list().items.each {
                roles.add(new K8sRole(
                        name: it.metadata.name,
                        namespace: it.metadata.namespace,
                        clusterRole: false,
                        labels: it.metadata.labels ? it.metadata.labels : [:],
                        annotations: it.metadata.annotations ? it.metadata.annotations : [:],
                        rules: it.rules.collect {
                            new K8sPolicyRule(
                                    verbs: it.verbs,
                                    apiGroups: it.apiGroups,
                                    resources: it.resources,
                                    nonResourceUrls: it.nonResourceURLs,
                                    resourceNames: it.resourceNames
                            )
                        }
                ))
            }
            return roles
        }
    }

    def createRole(K8sRole role) {
        withRetry(1, 2) {
            Role r = new Role(
                    metadata: new ObjectMeta(
                            name: role.name,
                            namespace: role.namespace,
                            labels: role.labels,
                            annotations: role.annotations
                    ),
                    rules: role.rules.collect {
                        new PolicyRule(
                                verbs: it.verbs,
                                apiGroups: it.apiGroups,
                                resources: it.resources,
                                nonResourceURLs: it.nonResourceUrls,
                                resourceNames: it.resourceNames
                        )
                    }
            )
            role.uid = client.rbac().roles().inNamespace(role.namespace).createOrReplace(r).metadata.uid
        }
    }

    def deleteRole(K8sRole role) {
        withRetry(1, 2) {
            client.rbac().roles().inNamespace(role.namespace).withName(role.name).delete()
        }
    }

    /*
        RoleBindings
     */

    List<K8sRoleBinding> getRoleBindings() {
        return evaluateWithRetry(2, 3) {
            def bindings = []
            client.rbac().roleBindings().inAnyNamespace().list().items.each {
                def b = new K8sRoleBinding(
                        new K8sRole(
                                name: it.metadata.name,
                                namespace: it.metadata.namespace,
                                clusterRole: false,
                                labels: it.metadata.labels ? it.metadata.labels : [:],
                                annotations: it.metadata.annotations ? it.metadata.annotations : [:]
                        ),
                        it.subjects.collect {
                    new K8sSubject(kind: it.kind, name: it.name, namespace: it.namespace ?: "")
                        }
                )
                def uid = it.roleRef.kind == "Role" ?
                        client.rbac().roles()
                                .inNamespace(it.metadata.namespace)
                                .withName(it.roleRef.name).get()?.metadata?.uid :
                        client.rbac().clusterRoles().withName(it.roleRef.name).get()?.metadata?.uid
                b.roleRef.uid = uid ?: ""
                bindings.add(b)
            }
            return bindings
        }
    }

    def createRoleBinding(K8sRoleBinding roleBinding) {
        withRetry(1, 2) {
            RoleBinding r = new RoleBinding(
                    metadata: new ObjectMeta(
                            name: roleBinding.name,
                            namespace: roleBinding.namespace,
                            labels: roleBinding.labels,
                            annotations: roleBinding.annotations
                    ),
                    subjects: roleBinding.subjects.collect {
                        new Subject(kind: it.kind, name: it.name, namespace: it.namespace)
                    },
                    roleRef: new RoleRef(
                            name: roleBinding.roleRef.name,
                            kind: roleBinding.roleRef.clusterRole ? "ClusterRole" : "Role"
                    )
            )
            client.rbac().roleBindings().inNamespace(roleBinding.namespace).createOrReplace(r)
        }
    }

    def deleteRoleBinding(K8sRoleBinding roleBinding) {
        withRetry(1, 2) {
            client.rbac().roleBindings()
                    .inNamespace(roleBinding.namespace)
                    .withName(roleBinding.name)
                    .delete()
        }
    }

    /*
        ClusterRoles
     */

    List<K8sRole> getClusterRoles() {
        return evaluateWithRetry(2, 3) {
            def clusterRoles = []
            client.rbac().clusterRoles().inAnyNamespace().list().items.each {
                clusterRoles.add(new K8sRole(
                        name: it.metadata.name,
                        namespace: "",
                        clusterRole: true,
                        labels: it.metadata.labels ? it.metadata.labels : [:],
                        annotations: it.metadata.annotations ? it.metadata.annotations : [:],
                        rules: it.rules.collect {
                            new K8sPolicyRule(
                                    verbs: it.verbs,
                                    apiGroups: it.apiGroups,
                                    resources: it.resources,
                                    nonResourceUrls: it.nonResourceURLs,
                                    resourceNames: it.resourceNames
                            )
                        }
                ))
            }
            return clusterRoles
        }
    }

    def createClusterRole(K8sRole role) {
        withRetry(2, 3) {
            ClusterRole r = new ClusterRole(
                    metadata: new ObjectMeta(
                            name: role.name,
                            labels: role.labels,
                            annotations: role.annotations
                    ),
                    rules: role.rules.collect {
                        new PolicyRule(
                                verbs: it.verbs,
                                apiGroups: it.apiGroups,
                                resources: it.resources,
                                nonResourceURLs: it.nonResourceUrls,
                                resourceNames: it.resourceNames
                        )
                    }
            )
            role.uid = client.rbac().clusterRoles().createOrReplace(r).metadata.uid
        }
    }

    def deleteClusterRole(K8sRole role) {
        withRetry(2, 3) {
            client.rbac().clusterRoles().withName(role.name).delete()
        }
    }

    /*
        ClusterRoleBindings
     */

    List<K8sRoleBinding> getClusterRoleBindings() {
        return evaluateWithRetry(2, 3) {
            def clusterBindings = []
            client.rbac().clusterRoleBindings().inAnyNamespace().list().items.each {
                def b = new K8sRoleBinding(
                        new K8sRole(
                                name: it.metadata.name,
                                namespace: "",
                                clusterRole: true,
                                labels: it.metadata.labels ? it.metadata.labels : [:],
                                annotations: it.metadata.annotations ? it.metadata.annotations : [:]
                        ),
                        it.subjects.collect {
                    new K8sSubject(kind: it.kind, name: it.name, namespace: it.namespace ?: "")
                        }
                )
                def uid = client.rbac().clusterRoles().withName(it.roleRef.name).get()?.metadata?.uid ?:
                        client.rbac().roles()
                                .inNamespace(it.metadata.namespace)
                                .withName(it.roleRef.name).get()?.metadata?.uid
                b.roleRef.uid = uid ?: ""
                clusterBindings.add(b)
            }
            return clusterBindings
        }
    }

    def createClusterRoleBinding(K8sRoleBinding roleBinding) {
        withRetry(2, 3) {
            ClusterRoleBinding r = new ClusterRoleBinding(
                    metadata: new ObjectMeta(
                            name: roleBinding.name,
                            labels: roleBinding.labels,
                            annotations: roleBinding.annotations
                    ),
                    subjects: roleBinding.subjects.collect {
                        new Subject(kind: it.kind, name: it.name, namespace: it.namespace)
                    },
                    roleRef: new RoleRef(
                            name: roleBinding.roleRef.name,
                            kind: roleBinding.roleRef.clusterRole ? "ClusterRole" : "Role"
                    )
            )
            client.rbac().clusterRoleBindings().createOrReplace(r)
        }
    }

    def deleteClusterRoleBinding(K8sRoleBinding roleBinding) {
        withRetry(2, 3) {
            client.rbac().clusterRoleBindings().withName(roleBinding.name).delete()
        }
    }

    /*
        PodSecurityPolicies
    */

    protected generatePspRole() {
        def rules = [new K8sPolicyRule(
                apiGroups: ["policy"],
                resources: ["podsecuritypolicies"],
                resourceNames: ["allow-all-for-test"],
                verbs: ["use"]
        ),]
        return new K8sRole(
                name: "allow-all-for-test",
//                namespace: namespace,
                clusterRole: true,
                rules: rules
        )
    }

    protected generatePspRoleBinding(String namespace) {
        def roleBinding =  new K8sRoleBinding(
                name: "allow-all-for-test-" + namespace,
                namespace: namespace,
                roleRef: generatePspRole(),
                subjects: [new K8sSubject(
                        name: "default",
                        namespace: namespace,
                        kind: "ServiceAccount"
                )]
        )
        return roleBinding
    }

    protected defaultPspForNamespace(String namespace) {
        PodSecurityPolicy psp = new PodSecurityPolicyBuilder().withNewMetadata()
                .withName("allow-all-for-test")
                .endMetadata()
                .withNewSpec()
                .withPrivileged(true)
                .withAllowPrivilegeEscalation(true)
                .withAllowedCapabilities("*")
                .withVolumes("*")
                .withHostNetwork(true)
                .withHostPorts(new HostPortRange(65535, 0))
                .withHostIPC(true)
                .withHostPID(true)
                .withNewRunAsUser().withRule("RunAsAny").endRunAsUser()
                .withNewSeLinux().withRule("RunAsAny").endSeLinux()
                .withNewSupplementalGroups().withRule("RunAsAny").endSupplementalGroups()
                .withNewFsGroup().withRule("RunAsAny").endFsGroup()
                .endSpec()
                .build()
        client.policy().podSecurityPolicies().createOrReplace(psp)
        createClusterRole(generatePspRole())
        createClusterRoleBinding(generatePspRoleBinding(namespace))
    }

    /*
        Jobs
     */

    def getJobCount(String ns = null) {
        return evaluateWithRetry(2, 3) {
            return client.batch().jobs().inNamespace(ns).list().getItems().collect { it.metadata.name }
        }
    }

    /*
        ConfigMaps
    */

    def createConfigMap(String name, Map<String,String> data, String namespace = this.namespace) {
        ConfigMap configMap = new ConfigMap(
                apiVersion: "v1",
                kind: "ConfigMap",
                data: data,
                metadata: new ObjectMeta(
                        name: name
                )
        )

        try {
            client.configMaps().inNamespace(namespace).createOrReplace(configMap)
            Timer t = new Timer(20, 3)
            while (t.IsValid()) {
                ConfigMap foundAfterCreate = client.configMaps().inNamespace(namespace).withName(name).get()
                if (foundAfterCreate != null) {
                    println name + ": config map created."
                    return foundAfterCreate
                }
            }
            throw new RuntimeException("Config map not found after create")
        } catch (Exception e) {
            println("Error creating configMap: " + e.toString())
        }
    }

    /*
        Misc/Helper Methods
    */

    def execInContainer(Deployment deployment, String cmd) {
        // Wait for container 0 to be running first.
        def timer = new Timer(30, 1)
        while (timer.IsValid()) {
            println "First container in pod ${deployment.pods.get(0).name} not yet running ..."
            def p = client.pods().inNamespace(deployment.namespace).withName(deployment.pods.get(0).name).get()
            if (p == null || p.status.containerStatuses.size() == 0) {
                continue
            }
            def status = p.status.containerStatuses.get(0)
            if (status.state.running != null) {
                break
            }
        }

        ScheduledExecutorService executorService = Executors.newScheduledThreadPool(20)
        try {
            CountDownLatch latch = new CountDownLatch(1)
            ExecWatch watch =
                    client.pods().inNamespace(deployment.namespace).withName(deployment.pods.get(0).name)
                            .redirectingOutput().usingListener(new ExecListener() {
                @Override
                void onOpen(Response response) {
                }

                @Override
                void onFailure(Throwable t, Response response) {
                    latch.countDown()
                }

                @Override
                void onClose(int code, String reason) {
                    latch.countDown()
                }
            }).exec(cmd.split(" "))
            BlockingInputStreamPumper pump = new BlockingInputStreamPumper(watch.getOutput(), new SystemOutCallback())
            Future<String> outPumpFuture = executorService.submit(pump, "Done")
            executorService.scheduleAtFixedRate(new FutureChecker("Exec", cmd, outPumpFuture), 0, 2, TimeUnit.SECONDS)

            latch.await(30, TimeUnit.SECONDS)
            watch.close()
            pump.close()
        } catch (Exception e) {
            println "Error exec'ing in pod: ${e.toString()}"
            return false
        }
        executorService.shutdown()
        return true
    }

    String generateYaml(Object orchestratorObject) {
        if (orchestratorObject instanceof NetworkPolicy) {
            return YamlGenerator.toYaml(createNetworkPolicyObject(orchestratorObject))
        }

        return ""
    }

    String getNameSpace() {
        return this.namespace
    }

    String getSensorContainerName() {
        return evaluateWithRetry(2, 3) {
            return client.pods().inNamespace("stackrox").list().items.find {
                it.metadata.name.startsWith("sensor")
            }.metadata.name
        }
    }

    def waitForSensor() {
        def start = System.currentTimeMillis()
        def running = client.apps().deployments()
                .inNamespace("stackrox")
                .withName("sensor")
                .get().status.readyReplicas < 1
        while (!running && (System.currentTimeMillis() - start) < 30000) {
            println "waiting for sensor to come back online. Trying again in 1s..."
            sleep 1000
            running = client.apps().deployments()
                    .inNamespace("stackrox")
                    .withName("sensor")
                    .get().status.readyReplicas < 1
        }
        if (!running) {
            println "Failed to detect sensor came back up within 30s... Future tests may be impacted."
        }
    }

    int getAllDeploymentTypesCount(String ns = null) {
        return getDeploymentCount(ns).size() +
                getDaemonSetCount(ns).size() +
                getStaticPodCount(ns).size() +
                getStatefulSetCount(ns).size() +
                getJobCount(ns).size()
    }

    /*
        Private K8S Support functions
    */

    def createDeploymentNoWait(Deployment deployment) {
        deployment.getNamespace() != null ?: deployment.setNamespace(this.namespace)

        // Create service if needed
        if (deployment.exposeAsService) {
            createService(deployment)
        }

        K8sDeployment d = new K8sDeployment(
                metadata: new ObjectMeta(
                        name: deployment.name,
                        namespace: deployment.namespace,
                        labels: deployment.labels,
                        annotations: deployment.annotation,
                ),
                spec: new DeploymentSpec(
                        selector: new LabelSelector(null, deployment.labels),
                        replicas: deployment.replicas,
                        minReadySeconds: 15,
                        template: new PodTemplateSpec(
                                metadata: new ObjectMeta(
                                        name: deployment.name,
                                        namespace: deployment.namespace,
                                        labels: deployment.labels + ["deployment": deployment.name],
                                ),
                                spec: generatePodSpec(deployment)
                        )

                )
        )

        try {
            client.apps().deployments().inNamespace(deployment.namespace).createOrReplace(d)
            println("Told the orchestrator to create " + deployment.name)
            if (deployment.createLoadBalancer) {
                waitForLoadBalancer(deployment)
            }
            return true
        } catch (Exception e) {
            println("Error creating k8s deployment: " + e.toString())
            return false
        }
    }

    def waitForDeploymentAndPopulateInfo(Deployment deployment) {
        try {
            deployment.deploymentUid = waitForDeploymentCreation(
                    deployment.getName(),
                    deployment.getNamespace(),
                    deployment.skipReplicaWait
            )
            updateDeploymentDetails(deployment)
        } catch (Exception e) {
            println("Error while waiting for deployment/populating deployment info: " + e.toString())
        }
    }

    def waitForDeploymentCreation(String deploymentName, String namespace, Boolean skipReplicaWait = false) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for ${deploymentName} to start"
            K8sDeployment d = this.deployments.inNamespace(namespace).withName(deploymentName).get()
            getAndPrintPods(namespace, deploymentName)
            if (d == null) {
                println "${deploymentName} not found yet"
                continue
            } else if (skipReplicaWait) {
                // If skipReplicaWait is set, we still want to sleep for a few seconds to allow the deployment
                // to work its way through the system.
                sleep(sleepDurationSeconds * 1000)
                println "${deploymentName}: deployment created (skipped replica wait)."
                return
            }
            if (d.getStatus().getReadyReplicas() == d.getSpec().getReplicas()) {
                println "All ${d.getSpec().getReplicas()} replicas found " +
                        "in ready state for ${deploymentName}"
                println "Took ${t.SecondsSince()} seconds for k8s deployment ${deploymentName}"
                return d.getMetadata().getUid()
            }
            println "${d.getStatus().getReadyReplicas()}/" +
                    "${d.getSpec().getReplicas()} are in the ready state for ${deploymentName}"
        }
    }

    def createDaemonSetNoWait(DaemonSet daemonSet) {
        daemonSet.getNamespace() != null ?: daemonSet.setNamespace(this.namespace)

        K8sDaemonSet ds = new K8sDaemonSet(
                metadata: new ObjectMeta(
                        name: daemonSet.name,
                        namespace: daemonSet.namespace,
                        labels: daemonSet.labels
                ),
                spec: new DaemonSetSpec(
                        minReadySeconds: 15,
                        selector: new LabelSelector(null, daemonSet.labels),
                        template: new PodTemplateSpec(
                                metadata: new ObjectMeta(
                                        name: daemonSet.name,
                                        namespace: daemonSet.namespace,
                                        labels: daemonSet.labels
                                ),
                                spec: generatePodSpec(daemonSet)
                        )
                )
        )

        try {
            this.daemonsets.inNamespace(daemonSet.namespace).createOrReplace(ds)
            println("Told the orchestrator to create " + daemonSet.getName())
        } catch (Exception e) {
            println("Error creating k8s deployment" + e.toString())
        }
    }

    def waitForDaemonSetAndPopulateInfo(DaemonSet daemonSet) {
        try {
            daemonSet.deploymentUid = waitForDaemonSetCreation(
                    daemonSet.getName(),
                    daemonSet.getNamespace(),
                    daemonSet.skipReplicaWait
            )
            updateDeploymentDetails(daemonSet)
        } catch (Exception e) {
            println("Error while waiting for daemonset/populating daemonset info: " + e.toString())
        }
    }

    def waitForDaemonSetCreation(String name, String namespace, Boolean skipReplicaWait = false) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for ${name} to start"
            K8sDaemonSet d = this.daemonsets.inNamespace(namespace).withName(name).get()
            getAndPrintPods(namespace, name)
            if (d == null) {
                println "${name} not found yet"
                continue
            } else if (skipReplicaWait) {
                // If skipReplicaWait is set, we still want to sleep for a few seconds to allow the deployment
                // to work its way through the system.
                sleep(sleepDurationSeconds * 1000)
                println "${name}: daemonset created (skipped replica wait)."
                return
            }
            if (d.getStatus().getCurrentNumberScheduled() == d.getStatus().getDesiredNumberScheduled()) {
                println "All ${d.getStatus().getDesiredNumberScheduled()} replicas found in ready state for ${name}"
                return d.getMetadata().getUid()
            }
            println "${d.getStatus().getCurrentNumberScheduled()}/" +
                    "${d.getStatus().getDesiredNumberScheduled()} are in the ready state for ${name}"
        }
    }

    def generatePodSpec(Deployment deployment) {
        List<ContainerPort> depPorts = deployment.ports.collect {
            k, v -> new ContainerPort(
                    k as Integer,
                    null,
                    null,
                    "port" + (k as String),
                    v as String
            )
        }

        List<LocalObjectReference> imagePullSecrets = new LinkedList<>()
        for (String str : deployment.getImagePullSecret()) {
            LocalObjectReference obj = new LocalObjectReference(name: str)
            imagePullSecrets.add(obj)
        }

        List<EnvVar> envVars = deployment.env.collect {
            k, v -> new EnvVar(k, v, null)
        }

        deployment.envValueFromSecretKeyRef.forEach {
            String k, SecretKeyRef v -> envVars.add(new EnvVarBuilder()
                    .withName(k)
                    .withValueFrom(new EnvVarSourceBuilder()
                            .withSecretKeyRef(
                                    new SecretKeySelectorBuilder().withKey(v.key).withName(v.name).build())
                            .build())
                    .build())
        }

        deployment.envValueFromConfigMapKeyRef.forEach {
            String k, ConfigMapKeyRef v -> envVars.add(new EnvVarBuilder()
                    .withName(k)
                    .withValueFrom(new EnvVarSourceBuilder()
                            .withConfigMapKeyRef(
                                    new ConfigMapKeySelectorBuilder().withKey(v.key).withName(v.name).build())
                            .build())
                    .build())
        }

        deployment.envValueFromFieldRef.forEach {
            String k, String fieldPath -> envVars.add(new EnvVarBuilder()
                    .withName(k)
                    .withValueFrom(new EnvVarSourceBuilder()
                            .withFieldRef(
                                    new ObjectFieldSelectorBuilder().withFieldPath(fieldPath).build())
                            .build())
                    .build())
        }

        deployment.envValueFromResourceFieldRef.forEach {
            String k, String resource -> envVars.add(new EnvVarBuilder()
                    .withName(k)
                    .withValueFrom(new EnvVarSourceBuilder()
                            .withResourceFieldRef(
                                    new ResourceFieldSelectorBuilder().withResource(resource).build())
                            .build())
                    .build())
        }

        List<EnvFromSource> envFrom = new LinkedList<>()
        for (String secret : deployment.getEnvFromSecrets()) {
            envFrom.add(new EnvFromSource(null, null, new SecretEnvSource(secret, false)))
        }
        for (String configMapName : deployment.getEnvFromConfigMaps()) {
            envFrom.add(new EnvFromSource(new ConfigMapEnvSource(configMapName, false), null, null))
        }

        List<Volume> volumes = []
        deployment.volumes.each {
            v -> Volume vol = new Volume(
                    name: v.name,
                    hostPath: v.hostPath ? new HostPathVolumeSource(
                            path: v.mountPath,
                            type: "Directory") :
                            null,
                    secret: deployment.secretNames.get(v.name) ?
                            new SecretVolumeSource(secretName: deployment.secretNames.get(v.name)) :
                            null
            )
            volumes.add(vol)
        }

        List<VolumeMount> volMounts = []
        deployment.volumeMounts.each {
            v -> VolumeMount volMount = new VolumeMount(
                    mountPath: v.mountPath,
                    name: v.name,
                    readOnly: v.readOnly
            )
            volMounts.add(volMount)
        }

        Map<String , Quantity> limits = new HashMap<>()
        for (String key:deployment.limits.keySet()) {
            Quantity quantity = new Quantity(deployment.limits.get(key))
            limits.put(key, quantity)
        }

        Map<String , Quantity> requests = new HashMap<>()
        for (String key:deployment.request.keySet()) {
            Quantity quantity = new Quantity(deployment.request.get(key))
            requests.put(key, quantity)
        }

        Container container = new Container(
                name: deployment.name,
                image: deployment.image,
                command: deployment.command,
                args: deployment.args,
                ports: depPorts,
                volumeMounts: volMounts,
                env: envVars,
                envFrom: envFrom,
                resources: new ResourceRequirements(limits, requests),
                securityContext: new SecurityContext(privileged: deployment.isPrivileged,
                                                     readOnlyRootFilesystem: deployment.readOnlyRootFilesystem,
                                                     capabilities: new Capabilities(add: deployment.addCapabilities,
                                                                                    drop: deployment.dropCapabilities))
        )

        PodSpec podSpec = new PodSpec(
                containers: [container],
                volumes: volumes,
                imagePullSecrets: imagePullSecrets,
                hostNetwork: deployment.hostNetwork,
                serviceAccountName: deployment.serviceAccountName
        )
        return podSpec
    }

    def updateDeploymentDetails(Deployment deployment) {
        // Filtering pod query by using the "name=<name>" because it should always be present in the deployment
        // object - IF this is ever missing, it may cause problems fetching pod details
        def deployedPods = evaluateWithRetry(2, 3) {
            return client.pods().inNamespace(deployment.namespace).withLabel("name", deployment.name).list()
        }
        for (Pod pod : deployedPods.getItems()) {
            deployment.addPod(
                    pod.getMetadata().getName(),
                    pod.getMetadata().getUid(),
                    pod.getStatus().getContainerStatuses() != null ?
                            pod.getStatus().getContainerStatuses().stream().map {
                                container -> container.getContainerID()
                            }.collect(Collectors.toList()) :
                            [],
                    pod.getStatus().getPodIP()
            )
        }
    }

    protected io.fabric8.kubernetes.api.model.networking.NetworkPolicy createNetworkPolicyObject(NetworkPolicy policy) {
        def networkPolicy = new NetworkPolicyBuilder()
                .withApiVersion("networking.k8s.io/v1")
                .withKind("NetworkPolicy")
                .withNewMetadata()
                .withName(policy.name)

        if (policy.namespace) {
            networkPolicy.withNamespace(policy.namespace)
        }

        if (policy.labels != null) {
            networkPolicy.withLabels(policy.labels)
        }

        networkPolicy = networkPolicy.endMetadata().withNewSpec()

        if (policy.metadataPodSelector != null) {
            networkPolicy.withNewPodSelector().withMatchLabels(policy.metadataPodSelector).endPodSelector()
        }

        if (policy.types != null) {
            def polTypes = []
            for (NetworkPolicyTypes type : policy.types) {
                polTypes.add(type.toString())
            }
            networkPolicy.withPolicyTypes(polTypes)
        }

        if (policy.ingressPodSelector != null) {
            networkPolicy.withIngress(
                    new NetworkPolicyIngressRuleBuilder().withFrom(
                            new NetworkPolicyPeerBuilder()
                                    .withNewPodSelector()
                                    .withMatchLabels(policy.ingressPodSelector)
                                    .endPodSelector().build()
                    ).build()
            )
        }

        if (policy.egressPodSelector != null) {
            networkPolicy.withEgress(
                    new NetworkPolicyEgressRuleBuilder().withTo(
                            new NetworkPolicyPeerBuilder()
                                    .withNewPodSelector()
                                    .withMatchLabels(policy.egressPodSelector)
                                    .endPodSelector().build()
                    ).build()
            )
        }

        if (policy.ingressNamespaceSelector != null) {
            networkPolicy.withIngress(
                    new NetworkPolicyIngressRuleBuilder().withFrom(
                            new NetworkPolicyPeerBuilder()
                                    .withNewNamespaceSelector()
                                    .withMatchLabels(policy.ingressNamespaceSelector)
                                    .endNamespaceSelector().build()
                    ).build()
            )
        }

        if (policy.egressNamespaceSelector != null) {
            networkPolicy.withEgress(
                    new NetworkPolicyEgressRuleBuilder().withTo(
                            new NetworkPolicyPeerBuilder()
                                    .withNewNamespaceSelector()
                                    .withMatchLabels(policy.egressNamespaceSelector)
                                    .endNamespaceSelector().build()
                    ).build()
            )
        }

        return networkPolicy.endSpec().build()
    }

    V1beta1ValidatingWebhookConfiguration getAdmissionController() {
        println "get admission controllers stub"
    }

    def deleteAdmissionController(String name) {
        println "delete admission controllers stub: ${name}"
    }

    def createAdmissionController(V1beta1ValidatingWebhookConfiguration config) {
        println "create admission controllers stub: ${config}"
    }

    String createNamespace(String ns) {
        return evaluateWithRetry(2, 3) {
            Namespace namespace = new Namespace("v1", null, new ObjectMeta(name: ns), null, null)
            def namespaceId = client.namespaces().createOrReplace(namespace).metadata.getUid()
            defaultPspForNamespace(ns)
            return namespaceId
        }
    }

    def deleteNamespace(String ns) {
        withRetry(2, 3) {
            client.namespaces().withName(ns).delete()
        }
    }

    def waitForNamespaceDeletion(String ns, int retries = 20, int intervalSeconds = 3) {
        println "Waiting for namespace ${ns} to be deleted"
        Timer t = new Timer(retries, intervalSeconds)
        while (t.IsValid()) {
            if (client.namespaces().withName(ns).get() == null ) {
                println "K8s found that namespace ${ns} was deleted"
                return true
            }
            println "Retrying in ${intervalSeconds}..."
        }
        println "K8s did not detect that namespace ${ns} was deleted"
        return false
    }

    private static class SystemOutCallback implements Callback<byte[]> {
        @Override
        void call(byte[] data) {
            System.out.print(new String(data))
        }
    }

    private static class FutureChecker implements Runnable {
        private final String name
        private final String cmd
        private final Future<String> future

        private FutureChecker(String name, String cmd, Future<String> future) {
            this.name = name
            this.cmd = cmd
            this.future = future
        }

        @Override
        void run() {
            if (!future.isDone()) {
                System.out.println(name + ":[" + cmd + "] is not done yet")
            }
        }
    }
}
