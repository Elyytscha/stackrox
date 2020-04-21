package clairify

import (
	"net"
	"net/http"
	"time"

	gogoProto "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	clairConv "github.com/stackrox/rox/pkg/clair"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/httputil/proxy"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/registries"
	scannerTypes "github.com/stackrox/rox/pkg/scanners/types"
	"github.com/stackrox/rox/pkg/urlfmt"
	clairV1 "github.com/stackrox/scanner/api/v1"
	"github.com/stackrox/scanner/pkg/clairify/client"
	"github.com/stackrox/scanner/pkg/clairify/types"
)

const (
	typeString = "clairify"

	clientTimeout      = 5 * time.Minute
	maxConcurrentScans = 30
)

var (
	log = logging.LoggerForModule()
)

// Creator provides the type an scanners.Creator to add to the scanners Registry.
func Creator(set registries.Set) (string, func(integration *storage.ImageIntegration) (scannerTypes.ImageScanner, error)) {
	return typeString, func(integration *storage.ImageIntegration) (scannerTypes.ImageScanner, error) {
		scan, err := newScanner(integration, set)
		return scan, err
	}
}

type clairify struct {
	scannerTypes.ScanSemaphore

	client                *client.Clairify
	conf                  *storage.ClairifyConfig
	protoImageIntegration *storage.ImageIntegration
	activeRegistries      registries.Set
}

func newScanner(protoImageIntegration *storage.ImageIntegration, activeRegistries registries.Set) (*clairify, error) {
	clairifyConfig, ok := protoImageIntegration.IntegrationConfig.(*storage.ImageIntegration_Clairify)
	if !ok {
		return nil, errors.New("Clairify configuration required")
	}
	conf := clairifyConfig.Clairify
	if err := validateConfig(conf); err != nil {
		return nil, err
	}
	endpoint := urlfmt.FormatURL(conf.Endpoint, urlfmt.InsecureHTTP, urlfmt.NoTrailingSlash)

	dialer := net.Dialer{
		Timeout: 2 * time.Second,
	}

	tlsConfig, err := clientconn.TLSConfig(mtls.ScannerSubject, clientconn.TLSConfigOptions{
		UseClientCert: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize TLS config")
	}

	httpClient := &http.Client{
		Timeout: clientTimeout,
		Transport: &http.Transport{
			DialContext:     dialer.DialContext,
			TLSClientConfig: tlsConfig,
			Proxy:           proxy.FromConfig(),
		},
	}

	client := client.NewWithClient(endpoint, true, httpClient)
	scanner := &clairify{
		client:                client,
		conf:                  conf,
		protoImageIntegration: protoImageIntegration,
		activeRegistries:      activeRegistries,

		ScanSemaphore: scannerTypes.NewSemaphoreWithValue(maxConcurrentScans),
	}
	return scanner, nil
}

// Test initiates a test of the Clairify Scanner which verifies that we have the proper scan permissions
func (c *clairify) Test() error {
	return c.client.Ping()
}

func validateConfig(c *storage.ClairifyConfig) error {
	if c.GetEndpoint() == "" {
		return errors.New("endpoint parameter must be defined for Clairify")
	}
	return nil
}

func convertLayerToImageScan(image *storage.Image, layerEnvelope *clairV1.LayerEnvelope) *storage.ImageScan {
	if layerEnvelope == nil || layerEnvelope.Layer == nil {
		return nil
	}

	return &storage.ImageScan{
		ScanTime:   gogoProto.TimestampNow(),
		Components: clairConv.ConvertFeatures(image, layerEnvelope.Layer.Features),
	}
}

func v1ImageToClairifyImage(i *storage.Image) *types.Image {
	return &types.Image{
		SHA:      i.GetId(),
		Registry: i.GetName().GetRegistry(),
		Remote:   i.GetName().GetRemote(),
		Tag:      i.GetName().GetTag(),
	}
}

// Try many ways to retrieve a sha
func (c *clairify) getScan(image *storage.Image) (*clairV1.LayerEnvelope, error) {
	if env, err := c.client.RetrieveImageDataBySHA(utils.GetSHA(image), true, true); err == nil {
		return env, nil
	}
	return c.client.RetrieveImageDataByName(v1ImageToClairifyImage(image), true, true)
}

// GetScan retrieves the most recent scan
func (c *clairify) GetScan(image *storage.Image) (*storage.ImageScan, error) {
	// If we haven't retrieved any metadata then we won't be able to scan with Clairify
	if image.GetMetadata() == nil {
		return nil, nil
	}
	env, err := c.getScan(image)
	// If not found, then should trigger a scan
	if err != nil {
		if err != client.ErrorScanNotFound {
			return nil, err
		}
		if err := c.scan(image); err != nil {
			return nil, err
		}
		env, err = c.getScan(image)
		if err != nil {
			return nil, err
		}
	}
	scan := convertLayerToImageScan(image, env)
	if scan == nil {
		return nil, errors.New("malformed response from scanner")
	}
	return scan, nil
}

func (c *clairify) scan(image *storage.Image) error {
	rc := c.activeRegistries.GetRegistryMetadataByImage(image)
	if rc == nil {
		return nil
	}

	_, err := c.client.AddImage(rc.Username, rc.Password, &types.ImageRequest{
		Image:    utils.GetFullyQualifiedFullName(image),
		Registry: rc.URL,
		Insecure: rc.Insecure})
	return err
}

// Match decides if the image is contained within this scanner
func (c *clairify) Match(image *storage.ImageName) bool {
	return c.activeRegistries.Match(image)
}

func (c *clairify) Type() string {
	return typeString
}

func (c *clairify) Name() string {
	return c.protoImageIntegration.GetName()
}
