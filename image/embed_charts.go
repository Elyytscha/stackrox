package image

import (
	"embed"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/image/sensor"
	"github.com/stackrox/rox/pkg/helmtpl"
	"github.com/stackrox/rox/pkg/helmutil"
	"github.com/stackrox/rox/pkg/k8sutil/k8sobjects"
	"github.com/stackrox/rox/pkg/namespaces"
	rendererUtils "github.com/stackrox/rox/pkg/renderer/utils"
	"github.com/stackrox/rox/pkg/templates"
	"github.com/stackrox/rox/pkg/utils"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//go:embed templates/* assets/* templates/helm/stackrox-central/* templates/helm/stackrox-central/templates/* templates/helm/stackrox-secured-cluster/templates/* templates/helm/stackrox-secured-cluster/*

// AssetFS holds the helm charts
var AssetFS embed.FS

const (
	templatePath = "templates"

	// CentralServicesChartPrefix points to the new stackrox-central-services Helm Chart.
	CentralServicesChartPrefix = "templates/helm/stackrox-central"
	// SecuredClusterServicesChartPrefix points to the new stackrox-secured-cluster-services Helm Chart.
	SecuredClusterServicesChartPrefix = "templates/helm/stackrox-secured-cluster"
)

// These are the go based files from embedded chart filesystem
var (
	k8sScriptsFileMap = map[string]string{
		"templates/sensor/kubernetes/sensor.sh":        "templates/sensor.sh",
		"templates/sensor/kubernetes/delete-sensor.sh": "templates/delete-sensor.sh",
		"templates/common/ca-setup.sh":                 "templates/ca-setup-sensor.sh",
		"templates/common/delete-ca.sh":                "templates/delete-ca-sensor.sh",
	}

	osScriptsFileMap = map[string]string{
		"templates/sensor/openshift/sensor.sh":        "templates/sensor.sh",
		"templates/sensor/openshift/delete-sensor.sh": "templates/delete-sensor.sh",
		"templates/common/ca-setup.sh":                "templates/ca-setup-sensor.sh",
		"templates/common/delete-ca.sh":               "templates/delete-ca-sensor.sh",
	}
)

// Image holds the filesystem
type Image struct {
	fs fs.FS
}

// NewImage returns a new image instance, if a nil filesystem is given the default FS is used
func NewImage(fs fs.FS) *Image {
	return &Image{fs: fs}
}

var defaultImage = NewImage(AssetFS)

// GetDefaultImage returns an image with it's default embedded filesystem
func GetDefaultImage() *Image {
	return defaultImage
}

// LoadFileContents resolves a given file's contents.
func (i *Image) LoadFileContents(filename string) (string, error) {
	content, err := fs.ReadFile(AssetFS, filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ReadFileAndTemplate reads and renders the template for the file
func (i *Image) ReadFileAndTemplate(pathToFile string, funcs template.FuncMap) (*template.Template, error) {
	templatePath := path.Join(templatePath, pathToFile)
	contents, err := i.LoadFileContents(templatePath)
	if err != nil {
		return nil, err
	}

	tpl := template.New(templatePath)
	if funcs != nil {
		tpl = tpl.Funcs(funcs)
	}
	return tpl.Parse(contents)
}

func (i *Image) getChartTemplate(prefix string) (*helmtpl.ChartTemplate, error) {
	chartTplFiles, err := i.GetFiles(prefix)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s chart files from embedded filesystem", prefix)
	}
	chartTpl, err := helmtpl.Load(chartTplFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "loading %s helmtpl", prefix)
	}

	return chartTpl, nil
}

func (i *Image) mustGetSensorChart(values map[string]interface{}, certs *sensor.Certs) *chart.Chart {
	ch, err := i.getSensorChart(values, certs)
	utils.Must(err)
	return ch
}

// GetSensorChart returns the Helm chart for sensor
func (i *Image) GetSensorChart(values map[string]interface{}, certs *sensor.Certs) *chart.Chart {
	return i.mustGetSensorChart(values, certs)
}

// GetCentralServicesChartTemplate retrieves the StackRox Central Services Helm chart template.
func (i *Image) GetCentralServicesChartTemplate() (*helmtpl.ChartTemplate, error) {
	return i.getChartTemplate(CentralServicesChartPrefix)
}

// GetSecuredClusterServicesChartTemplate retrieves the StackRox Secured Cluster Services Helm chart template.
func (i *Image) GetSecuredClusterServicesChartTemplate() (*helmtpl.ChartTemplate, error) {
	return i.getChartTemplate(SecuredClusterServicesChartPrefix)
}

var (
	secretGVK = schema.GroupVersionKind{Version: "v1", Kind: "Secret"}
	// SensorCertObjectRefs are the objects in the sensor bundle that represents tls certs.
	SensorCertObjectRefs = map[k8sobjects.ObjectRef]struct{}{
		{GVK: secretGVK, Name: "sensor-tls", Namespace: namespaces.StackRox}:            {},
		{GVK: secretGVK, Name: "collector-tls", Namespace: namespaces.StackRox}:         {},
		{GVK: secretGVK, Name: "admission-control-tls", Namespace: namespaces.StackRox}: {},
	}
	// AdditionalCASensorSecretRef is the object in the sensor bundle that represents additional ca certs.
	AdditionalCASensorSecretRef = k8sobjects.ObjectRef{
		GVK:       secretGVK,
		Name:      "additional-ca-sensor",
		Namespace: namespaces.StackRox,
	}
)

// LoadAndInstantiateChartTemplate loads a Helm chart (meta-)template from an embed.FS, and instantiates
// it, using default chart values.
func (i *Image) LoadAndInstantiateChartTemplate(prefix string, metaVals map[string]interface{}) ([]*loader.BufferedFile, error) {
	chartTplFiles, err := i.GetFiles(prefix)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s chart files from embedded filesystems", prefix)
	}
	chartTpl, err := helmtpl.Load(chartTplFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "loading %s helmtpl", prefix)
	}

	// Render template files.
	renderedChartFiles, err := chartTpl.InstantiateRaw(metaVals)
	if err != nil {
		return nil, errors.Wrapf(err, "instantiating %s helmtpl", prefix)
	}

	// Apply .helmignore filtering rules, to be on the safe side (but keep .helmignore).
	renderedChartFiles, err = helmutil.FilterFiles(renderedChartFiles)
	if err != nil {
		return nil, errors.Wrap(err, "filtering instantiated helm chart files")
	}

	return renderedChartFiles, nil
}

// GetFiles returns all files recursively under a given path.
func (i *Image) GetFiles(prefix string) ([]*loader.BufferedFile, error) {
	prefix = strings.TrimSuffix(prefix, "/")
	var files []*loader.BufferedFile
	err := fs.WalkDir(i.fs, prefix, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(i.fs, p)
		if err != nil {
			return err
		}

		newPath := strings.TrimPrefix(p, prefix+"/")
		files = append(files, &loader.BufferedFile{
			Name: newPath,
			Data: data,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetSensorChartTemplate loads the Sensor helmtpl meta-template
func (i *Image) GetSensorChartTemplate() (*helmtpl.ChartTemplate, error) {
	chartTplFiles, err := i.GetFiles(SecuredClusterServicesChartPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "fetching sensor chart files from embedded filesystem")
	}

	return helmtpl.Load(chartTplFiles)
}

func (i *Image) getSensorChart(values map[string]interface{}, certs *sensor.Certs) (*chart.Chart, error) {
	chartTpl, err := i.GetSensorChartTemplate()
	if err != nil {
		return nil, errors.Wrap(err, "loading sensor chart template")
	}

	renderedFiles, err := chartTpl.InstantiateRaw(values)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating sensor chart template")
	}

	for certPath, data := range certs.Files {
		renderedFiles = append(renderedFiles, &loader.BufferedFile{
			Name: certPath,
			Data: data,
		})
	}

	if certOnly, _ := values["CertsOnly"].(bool); !certOnly {
		scriptFiles, err := i.addScripts(values)
		if err != nil {
			return nil, err
		}

		renderedFiles = append(renderedFiles, scriptFiles...)
	}

	return loader.LoadFiles(renderedFiles)
}

func (i *Image) addScripts(values map[string]interface{}) ([]*loader.BufferedFile, error) {
	if values["ClusterType"] == storage.ClusterType_KUBERNETES_CLUSTER.String() {
		return i.scripts(values, k8sScriptsFileMap)
	} else if values["ClusterType"] == storage.ClusterType_OPENSHIFT_CLUSTER.String() || values["ClusterType"] == storage.ClusterType_OPENSHIFT4_CLUSTER.String() {
		return i.scripts(values, osScriptsFileMap)
	} else {
		return nil, errors.Errorf("unable to create sensor bundle, invalid cluster type for cluster %s",
			values["ClusterName"])
	}
}

func (i *Image) scripts(values map[string]interface{}, filenameMap map[string]string) ([]*loader.BufferedFile, error) {
	var chartFiles []*loader.BufferedFile
	for srcFile, dstFile := range filenameMap {
		fileData, err := AssetFS.ReadFile(srcFile)
		if err != nil {
			return nil, err
		}
		t, err := template.New("temp").Funcs(rendererUtils.BuiltinFuncs).Parse(string(fileData))
		if err != nil {
			return nil, err
		}
		data, err := templates.ExecuteToBytes(t, values)
		if err != nil {
			return nil, err
		}
		chartFiles = append(chartFiles, &loader.BufferedFile{
			Name: dstFile,
			Data: data,
		})
	}

	return chartFiles, nil
}
