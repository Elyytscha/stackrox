package k8sintrospect

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/k8sutil"
	"github.com/stackrox/rox/pkg/logging"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

const (
	logWindow = 20 * time.Minute
)

var (
	log = logging.LoggerForModule()
)

type collector struct {
	started concurrency.Flag

	ctx    concurrency.ErrorWaitable
	filesC chan<- File

	cfg Config

	client        kubernetes.Interface
	dynamicClient dynamic.Interface

	errors []error
}

func newCollector(ctx concurrency.ErrorWaitable, k8sRESTConfig *rest.Config, cfg Config, filesC chan<- File) (*collector, error) {
	restConfigShallowCopy := *k8sRESTConfig
	oldWrapTransport := restConfigShallowCopy.WrapTransport
	restConfigShallowCopy.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
		if oldWrapTransport != nil {
			rt = oldWrapTransport(rt)
		}
		return httputil.ContextBoundRoundTripper(ctx, rt)
	}

	k8sClient, err := kubernetes.NewForConfig(&restConfigShallowCopy)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Kubernetes client set")
	}
	dynamicClient, err := dynamic.NewForConfig(&restConfigShallowCopy)
	if err != nil {
		return nil, errors.Wrap(err, "could not create dynamic Kubernetes client")
	}

	return &collector{
		ctx:           ctx,
		filesC:        filesC,
		cfg:           cfg,
		client:        k8sClient,
		dynamicClient: dynamicClient,
	}, nil
}

func generateFileName(obj k8sutil.Object, suffix string) string {
	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "_global"
	}

	app := obj.GetLabels()["app"]
	if app == "" {
		app = obj.GetLabels()["app.kubernetes.io/name"]
	}
	if app == "" {
		app = "_ungrouped"
	}
	return fmt.Sprintf("%s/%s/%s-%s%s", namespace, app, strings.ToLower(obj.GetObjectKind().GroupVersionKind().Kind), obj.GetName(), suffix)
}

func (c *collector) emitFile(obj k8sutil.Object, suffix string, data []byte) error {
	return c.emitFileRaw(generateFileName(obj, suffix), data)
}

func (c *collector) emitFileRaw(filePath string, data []byte) error {
	file := File{
		Path:     filePath,
		Contents: data,
	}

	select {
	case c.filesC <- file:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *collector) createDynamicClients() map[schema.GroupVersionKind]dynamic.NamespaceableResourceInterface {
	gvkSet := make(map[schema.GroupVersionKind]struct{})
	for _, objCfg := range c.cfg.Objects {
		gvkSet[objCfg.GVK] = struct{}{}
	}

	_, apiResourceLists, err := c.client.Discovery().ServerGroupsAndResources()
	if err != nil {
		c.recordError(errors.Wrap(err, "failed to obtain server resources"))
		return nil
	}

	clientMap := make(map[schema.GroupVersionKind]dynamic.NamespaceableResourceInterface)
	for _, apiResourceList := range apiResourceLists {
		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			c.recordError(errors.Wrap(err, "failed to parse group/version for API resource list"))
			continue
		}
		for _, apiResource := range apiResourceList.APIResources {
			if strings.ContainsRune(apiResource.Name, '/') {
				continue
			}

			gvk := schema.GroupVersionKind{
				Group:   gv.Group,
				Version: gv.Version,
				Kind:    apiResource.Kind,
			}
			if _, ok := gvkSet[gvk]; !ok {
				log.Infof("Resource %v not relevant", gvk)
				continue
			}
			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: apiResource.Name,
			}
			log.Infof("Creating client for resource %v", gvr)
			clientMap[gvk] = c.dynamicClient.Resource(gvr)
		}
	}

	return clientMap
}

func (c *collector) collectPodData(pod *v1.Pod) error {
	yamlData, err := yaml.Marshal(pod)
	if err != nil {
		yamlData = []byte(fmt.Sprintf("Error marshaling pod to YAML: %v", err))
	}
	if err := c.emitFile(pod, ".yaml", yamlData); err != nil {
		return err
	}

	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Running != nil {
			podLogOpts := &v1.PodLogOptions{
				Container:    container.Name,
				SinceSeconds: &[]int64{int64(logWindow / time.Second)}[0],
			}
			logsData, err := c.client.CoreV1().Pods(pod.GetNamespace()).GetLogs(pod.GetName(), podLogOpts).DoRaw()
			if err != nil {
				logsData = []byte(fmt.Sprintf("Error retrieving container logs: %v", err))
			}
			if err := c.emitFile(pod, fmt.Sprintf("-logs-%s.txt", container.Name), logsData); err != nil {
				return err
			}
		}

		if container.LastTerminationState.Terminated != nil {
			since := metav1.NewTime(container.LastTerminationState.Terminated.FinishedAt.Add(-logWindow))
			podLogOpts := &v1.PodLogOptions{
				Container: container.Name,
				Previous:  true,
				SinceTime: &since,
			}
			logsData, err := c.client.CoreV1().Pods(pod.GetNamespace()).GetLogs(pod.GetName(), podLogOpts).DoRaw()
			if err != nil {
				logsData = []byte(fmt.Sprintf("Error retrieving previous container logs: %v", err))
			}
			if err := c.emitFile(pod, fmt.Sprintf("-logs-%s-previous.txt", container.Name), logsData); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *collector) recordError(err error) {
	if err != nil {
		c.errors = append(c.errors, err)
	}
}

func (c *collector) collectObjectsData(ns string, cfg ObjectConfig, resourceClient dynamic.NamespaceableResourceInterface) error {
	log.Infof("Collecting data for object type %v", cfg.GVK)
	objList, err := resourceClient.Namespace(ns).List(metav1.ListOptions{})
	if err != nil {
		c.recordError(err)
		return nil
	}

	for _, obj := range objList.Items {
		if cfg.RedactionFunc != nil {
			cfg.RedactionFunc(&obj)
		}
		objYAML, err := yaml.Marshal(obj)
		if err != nil {
			objYAML = []byte(fmt.Sprintf("Failed to marshal object to YAML: %v", err))
		}
		if err := c.emitFile(&obj, ".yaml", objYAML); err != nil {
			return err
		}
	}

	return nil
}

func (c *collector) collectNamespaceData(ns string) (bool, error) {
	namespace, err := c.client.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if err != nil && k8sErrors.IsNotFound(err) {
		return false, nil
	}
	var nsYAML []byte
	if err == nil && namespace != nil {
		namespace.TypeMeta = metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		}
		nsYAML, err = yaml.Marshal(namespace)
	}

	if err != nil && len(nsYAML) == 0 {
		nsYAML = []byte(fmt.Sprintf("Failed to retrieve namespace: %v", err))
	}

	return true, c.emitFileRaw(fmt.Sprintf("%s/namespace-spec.yaml", ns), nsYAML)
}

func (c *collector) collectEventsData(ns string) error {
	eventList, err := c.client.CoreV1().Events(ns).List(metav1.ListOptions{})

	var eventsYAML []byte
	if err == nil && eventList != nil {
		eventsYAML, err = yaml.Marshal(eventList)
	}
	if err != nil && len(eventsYAML) == 0 {
		eventsYAML = []byte(fmt.Sprintf("Failed to retrieve events: %v", err))
	}

	return c.emitFileRaw(fmt.Sprintf("%s/event-list.yaml", ns), eventsYAML)
}

func (c *collector) collectPodsData(ns string) error {
	podList, err := c.client.CoreV1().Pods(ns).List(metav1.ListOptions{})
	if err != nil {
		c.recordError(errors.Wrapf(err, "could not list pods in namespace %q", ns))
		return nil
	}

	for _, pod := range podList.Items {
		pod.TypeMeta = metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		}

		if err := c.collectPodData(&pod); err != nil {
			return err
		}
	}

	return nil
}

func (c *collector) collectErrors() error {
	if len(c.errors) == 0 {
		return nil
	}

	var errorsText bytes.Buffer
	for _, err := range c.errors {
		fmt.Fprintln(&errorsText, err.Error())
	}

	return c.emitFileRaw("errors.txt", errorsText.Bytes())
}

// Run performs the collection process. May only be invoked a single time.
func (c *collector) Run() error {
	if c.started.TestAndSet(true) {
		return errors.New("collector already ran once")
	}
	defer close(c.filesC)

	clientMap := c.createDynamicClients()

	for _, ns := range c.cfg.Namespaces {
		nsExists, err := c.collectNamespaceData(ns)
		if err != nil {
			return err
		}
		if !nsExists {
			continue
		}

		if err := c.collectPodsData(ns); err != nil {
			return err
		}
		for _, objCfg := range c.cfg.Objects {
			objClient := clientMap[objCfg.GVK]
			if objClient == nil {
				continue
			}
			if err := c.collectObjectsData(ns, objCfg, objClient); err != nil {
				return err
			}
		}

		if err := c.collectEventsData(ns); err != nil {
			return err
		}
	}

	return c.collectErrors()
}
