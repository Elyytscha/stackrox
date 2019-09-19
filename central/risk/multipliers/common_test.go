package multipliers

import (
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/protoconv"
)

func getMockImages() []*storage.Image {
	return []*storage.Image{
		{
			Name: &storage.ImageName{
				FullName: "docker.io/library/nginx:1.10",
				Registry: "docker.io",
				Remote:   "library/nginx",
				Tag:      "1.10",
			},
			Scan: &storage.ImageScan{
				Components: []*storage.EmbeddedImageScanComponent{
					{
						Name:    "comp1",
						Version: "1.1.1",
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cvss: 5,
							},
							{
								Cvss: 5,
							},
						},
					},
					{
						Name:    "comp1",
						Version: "1.1.1",
						Vulns: []*storage.EmbeddedVulnerability{
							{
								Cvss: 5,
							},
							{
								Cvss: 5,
							},
						},
					},
				},
			},
			Metadata: &storage.ImageMetadata{
				V1: &storage.V1Metadata{
					Created: protoconv.ConvertTimeToTimestamp(time.Now().Add(-(180 * 24 * time.Hour))),
				},
			},
		},
		{
			Metadata: &storage.ImageMetadata{
				V1: &storage.V1Metadata{
					Created: protoconv.ConvertTimeToTimestamp(time.Now().Add(-(90 * 24 * time.Hour))),
				},
			},
		},
	}
}

func getMockDeployment() *storage.Deployment {
	return &storage.Deployment{
		Name:                         "mock-deployment",
		ServiceAccount:               "service-account",
		ClusterId:                    "cluster",
		Namespace:                    "namespace",
		AutomountServiceAccountToken: true,
		Ports: []*storage.PortConfig{
			{
				Name:          "Port1",
				ContainerPort: 22,
				Exposure:      storage.PortConfig_EXTERNAL,
			},
			{
				Name:          "Port2",
				ContainerPort: 23,
				Exposure:      storage.PortConfig_INTERNAL,
			},
			{
				Name:          "Port3",
				ContainerPort: 24,
				Exposure:      storage.PortConfig_NODE,
			},
		},
		Containers: []*storage.Container{
			{
				Name: "containerName",
				Volumes: []*storage.Volume{
					{
						Name:     "readonly",
						ReadOnly: true,
					},
				},
				Secrets: []*storage.EmbeddedSecret{
					{
						Name: "secret",
					},
				},
				SecurityContext: &storage.SecurityContext{
					AddCapabilities: []string{
						"ALL",
					},
					Privileged: true,
				},
				Image: types.ToContainerImage(getMockImages()[0]),
			},
			{
				Name:  "Container2",
				Image: types.ToContainerImage(getMockImages()[1]),
				Volumes: []*storage.Volume{
					{
						Name: "rw volume",
					},
				},
				SecurityContext: &storage.SecurityContext{},
			},
		},
	}
}
