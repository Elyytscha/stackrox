package fixtures

import "github.com/stackrox/rox/generated/api/v1"

// GetDeployment returns a Mock Deployment
func GetDeployment() *v1.Deployment {
	return &v1.Deployment{
		Name:        "nginx_server",
		Id:          "s79mdvmb6dsl",
		ClusterId:   "prod cluster",
		ClusterName: "prod cluster",
		Namespace:   "stackrox",
		Labels: []*v1.Deployment_KeyValue{
			{
				Key:   "com.docker.stack.namespace",
				Value: "prevent",
			},
			{
				Key:   "com.docker.swarm.service.name",
				Value: "prevent_sensor",
			},
		},
		Containers: []*v1.Container{
			{
				Image: &v1.Image{
					Name: &v1.ImageName{
						Sha:      "sha256:SHA1",
						Registry: "docker.io",
						Remote:   "library/nginx",
						Tag:      "1.10",
					},
					Metadata: &v1.ImageMetadata{
						Layers: []*v1.ImageLayer{
							{
								Instruction: "ADD",
								Value:       "FILE:blah",
							},
						},
					},
					Scan: &v1.ImageScan{
						Components: []*v1.ImageScanComponent{
							{
								Name: "name",
								Vulns: []*v1.Vulnerability{
									{
										Cve:     "cve",
										Cvss:    10,
										Summary: "Vuln summary",
									},
								},
							},
						},
					},
				},
				SecurityContext: &v1.SecurityContext{
					Privileged:       true,
					AddCapabilities:  []string{"SYS_ADMIN"},
					DropCapabilities: []string{"SYS_MODULE"},
				},
				Config: &v1.ContainerConfig{
					Env: []*v1.ContainerConfig_EnvironmentConfig{
						{
							Key:   "envkey",
							Value: "envvalue",
						},
					},
				},
				Volumes: []*v1.Volume{
					{
						Name:        "vol1",
						Source:      "/vol1",
						Destination: "/vol2",
						Type:        "host",
						ReadOnly:    true,
					},
				},
			},
			{
				Image: GetImage(),
			},
		},
	}
}
