package augmentedobjs

// This block enumerates custom tags.
const (
	DockerfileLineCustomTag      = "Dockerfile Line"
	ComponentAndVersionCustomTag = "Component And Version"
	NotWhitelistedCustomTag      = "Not Whitelisted"
	ContainerNameCustomTag       = "Container Name"
	ImageScanCustomTag           = "Image Scan"
)

type dockerfileLine struct {
	Line string `search:"Dockerfile Line"`
}

type componentAndVersion struct {
	ComponentAndVersion string `search:"Component And Version"`
}

type whitelistResult struct {
	NotWhitelisted bool `search:"Not Whitelisted"`
}
