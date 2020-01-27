package index

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/batcher"
	"github.com/stackrox/rox/pkg/blevehelper"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
	mappings "github.com/stackrox/rox/pkg/search/options/images"
	"github.com/stackrox/rox/pkg/utils"
)

const batchSize = 5000

const resourceName = "Image"

type indexerImpl struct {
	index *blevehelper.BleveWrapper
}

type imageWrapper struct {
	*storage.Image `json:"image"`
	Type           string `json:"type"`
}

func getComponentPath(s string) (string, []string) {
	return fmt.Sprintf("image.scan.components.%s", s), []string{"image", "scan", "components", s}
}

func getVulnPath(s string) (string, []string) {
	return fmt.Sprintf("image.scan.components.vulns.%s", s), []string{"image", "scan", "components", "vulns", s}
}

func getSubMappingOrPanic(mapping *mapping.DocumentMapping, subPath string) *mapping.DocumentMapping {
	subMapping := mapping.Properties[subPath]
	if subMapping == nil {
		utils.Should(errors.Errorf("no mapping with name %q", subPath))
	}
	return subMapping
}

func getFieldOrPanic(mapping *mapping.DocumentMapping) *mapping.FieldMapping {
	if len(mapping.Fields) == 0 {
		utils.Should(errors.Errorf("no fields are available for mapping: %+v", mapping))
	}
	return mapping.Fields[0]
}

func mapComponents(im *mapping.IndexMappingImpl, components []*storage.EmbeddedImageScanComponent, doc *document.Document) {
	imageMapping := getSubMappingOrPanic(im.TypeMapping[v1.SearchCategory_IMAGES.String()], "image")
	scanMapping := getSubMappingOrPanic(imageMapping, "scan")
	componentMapping := getSubMappingOrPanic(scanMapping, "components")

	componentNameMapping := getFieldOrPanic(getSubMappingOrPanic(componentMapping, "name"))
	componentNamePathStr, componentNamePath := getComponentPath("name")

	componentVersionMapping := getFieldOrPanic(getSubMappingOrPanic(componentMapping, "version"))
	componentVersionPathStr, componentVersionPath := getComponentPath("version")

	vulnMapping := getSubMappingOrPanic(componentMapping, "vulns")

	cveMapping := getFieldOrPanic(getSubMappingOrPanic(vulnMapping, "cve"))
	cvePathStr, cvePath := getVulnPath("cve")

	cvssMapping := getFieldOrPanic(getSubMappingOrPanic(vulnMapping, "cvss"))
	cvssPathStr, cvssPath := getVulnPath("cvss")

	fixedMapping := vulnMapping.Properties["SetFixedBy"].Properties["fixed_by"].Fields[0]

	fixedPathStr := "image.scan.components.vulns.SetFixedBy.fixed_by"
	fixedPath := strings.Split("image.scan.components.vulns.SetFixedBy.fixed_by", ".")

	walkContext := im.NewWalkContext(doc, imageMapping)

	for i, c := range components {
		componentIndex := []uint64{uint64(i)}

		componentNameMapping.ProcessString(c.GetName(), componentNamePathStr, componentNamePath, componentIndex, walkContext)
		componentVersionMapping.ProcessString(c.GetVersion(), componentVersionPathStr, componentVersionPath, componentIndex, walkContext)

		for j, vuln := range c.GetVulns() {
			vulnIndex := []uint64{uint64(i), uint64(j)}

			cveMapping.ProcessString(vuln.GetCve(), cvePathStr, cvePath, vulnIndex, walkContext)
			cvssMapping.ProcessFloat64(float64(vuln.GetCvss()), cvssPathStr, cvssPath, vulnIndex, walkContext)
			fixedMapping.ProcessString(vuln.GetFixedBy(), fixedPathStr, fixedPath, vulnIndex, walkContext)
		}
	}
}

func (b *indexerImpl) optimizedMapDocument(wrapper *imageWrapper) (*document.Document, error) {
	doc := document.NewDocument(wrapper.GetId())

	components := wrapper.GetScan().GetComponents()
	if wrapper.GetScan() != nil {
		wrapper.Scan.Components = nil
		defer func() {
			wrapper.Scan.Components = components
		}()
	}
	if err := b.index.Mapping().MapDocument(doc, wrapper); err != nil {
		return nil, err
	}

	mapComponents(b.index.Mapping().(*mapping.IndexMappingImpl), components, doc)
	return doc, nil
}

func (b *indexerImpl) AddImage(image *storage.Image) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Image")

	wrapper := &imageWrapper{
		Image: image,
		Type:  v1.SearchCategory_IMAGES.String(),
	}
	if err := b.index.Index.Index(image.GetId(), wrapper); err != nil {
		return err
	}
	return b.index.IncTxnCount()
}

func (b *indexerImpl) AddImages(images []*storage.Image) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Image")
	batchManager := batcher.New(len(images), batchSize)
	for {
		start, end, ok := batchManager.Next()
		if !ok {
			break
		}
		if err := b.processBatch(images[start:end]); err != nil {
			return err
		}
	}
	return b.index.IncTxnCount()
}

func (b *indexerImpl) processBatch(images []*storage.Image) error {
	batch := b.index.NewBatch()
	for _, image := range images {
		if err := batch.Index(image.GetId(), &imageWrapper{
			Image: image,
			Type:  v1.SearchCategory_IMAGES.String(),
		}); err != nil {
			return err
		}
	}
	return b.index.Batch(batch)
}

func (b *indexerImpl) DeleteImage(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Image")
	if err := b.index.Delete(id); err != nil {
		return err
	}
	return b.index.IncTxnCount()
}

func (b *indexerImpl) DeleteImages(ids []string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.RemoveMany, "Image")
	batch := b.index.NewBatch()
	for _, id := range ids {
		batch.Delete(id)
	}
	if err := b.index.Batch(batch); err != nil {
		return err
	}
	return b.index.IncTxnCount()
}

func (b *indexerImpl) GetTxnCount() uint64 {
	return b.index.GetTxnCount()
}

func (b *indexerImpl) ResetIndex() error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Reset, "Image")
	return blevesearch.ResetIndex(v1.SearchCategory_IMAGES, b.index.Index)
}

func (b *indexerImpl) Search(q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Image")
	return blevesearch.RunSearchRequest(v1.SearchCategory_IMAGES, q, b.index.Index, mappings.OptionsMap, opts...)
}

func (b *indexerImpl) SetTxnCount(seq uint64) error {
	return b.index.SetTxnCount(seq)
}
