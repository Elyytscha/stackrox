package gcp

import (
	"github.com/spf13/cobra"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/renderer"
	"github.com/stackrox/rox/pkg/version"
	centralgenerate "github.com/stackrox/rox/roxctl/central/generate"
	"github.com/stackrox/rox/roxctl/common/util"
)

// Generate defines the gcp generate command.
func Generate() *cobra.Command {
	c := &cobra.Command{
		Use: "generate",
		RunE: util.RunENoArgs(func(c *cobra.Command) error {
			valuesFilename, _ := c.Flags().GetString("values-file")
			outputDir, _ := c.Flags().GetString("output-dir")
			values, err := loadValues(valuesFilename)
			if err != nil {
				return err
			}

			return generate(outputDir, values)
		}),
	}
	c.PersistentFlags().String("values-file", "/data/final_values.yaml", "path to the input yaml values file")
	c.PersistentFlags().String("output-dir", "", "path to the output helm chart directory")

	return c
}

func generate(outputDir string, values *Values) error {
	config := renderer.Config{
		Version:        version.GetMainVersion(),
		OutputDir:      outputDir,
		GCPMarketplace: true,
		ClusterType:    storage.ClusterType_KUBERNETES_CLUSTER,
		K8sConfig: &renderer.K8sConfig{
			AppName: values.Name,
			CommonConfig: renderer.CommonConfig{
				MainImage:      values.MainImage,
				ScannerImage:   values.ScannerImage,
				ScannerDBImage: values.ScannerDBImage,
			},
			DeploymentFormat: v1.DeploymentFormat_HELM,
			OfflineMode:      false,
		},
		External: &renderer.ExternalPersistence{
			Name:         values.PVCName,
			StorageClass: values.PVCStorageclass,
			Size:         values.PVCSize,
		},
		LicenseData: []byte(values.License),
	}

	switch values.Network {
	case "Load Balancer":
		config.K8sConfig.LoadBalancerType = v1.LoadBalancerType_LOAD_BALANCER
	case "Node Port":
		config.K8sConfig.LoadBalancerType = v1.LoadBalancerType_NODE_PORT
	case "None":
		config.K8sConfig.LoadBalancerType = v1.LoadBalancerType_NONE
	}

	return centralgenerate.OutputZip(config)
}
