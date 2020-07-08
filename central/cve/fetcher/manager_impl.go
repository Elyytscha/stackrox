package fetcher

import (
	"archive/zip"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/facebookincubator/nvdtools/cvefeed/nvd/schema"
	"github.com/pkg/errors"
	"github.com/stackrox/k8s-istio-cve-pusher/nvd"
	"github.com/stackrox/rox/central/cve/converter"
	cveDataStore "github.com/stackrox/rox/central/cve/datastore"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/sac"
	pkgSearch "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/throttle"
)

var (
	cveElevatedCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Cluster, resources.Image),
		))

	connectionDropThrottle = throttle.NewDropThrottle(10 * time.Minute)
)

type mode int

const (
	online = iota
	offline
	unknown
	k8sIstioCveZipName = "k8s-istio.zip"
)

// Init copies build time CVEs to persistent volume
func (m *k8sIstioCVEManagerImpl) initialize() {
	offlineModeSetting := env.OfflineModeEnv.Setting()
	if offlineModeSetting == "true" {
		m.mgrMode = offline
	} else {
		m.mgrMode = online
	}

	if err := copyCVEsFromPreloadedToPersistentDirIfAbsent(converter.K8s); err != nil {
		log.Errorf("could not copy preloaded k8s CVE files to persistent volume %q:%v", path.Join(persistentCVEsPath, commonCveDir, k8sCVEsDir), err)
		return
	}
	log.Infof("successfully copied preloaded k8s CVE files to persistent volume: %q", path.Join(persistentCVEsPath, commonCveDir, k8sCVEsDir))

	if err := copyCVEsFromPreloadedToPersistentDirIfAbsent(converter.Istio); err != nil {
		log.Errorf("could not copy preloaded istio CVE files to persistent volume %q: %v", path.Join(persistentCVEsPath, commonCveDir, istioCVEsDir), err)
		return
	}
	log.Infof("successfully copied preloaded CVE istio files to persistent volume: %q", path.Join(persistentCVEsPath, commonCveDir, istioCVEsDir))

	m.k8sCVEMgr.initialize()
	m.istioCVEMgr.initialize()
}

// Fetch (works only in online mode) fetches new CVEs and reconciles them
func (m *k8sIstioCVEManagerImpl) Start() {
	if m.mgrMode != online {
		log.Error("can't fetch in non-online mode")
		return
	}

	ticker := time.NewTicker(fetchDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.reconcileAllCVEsInOnlineMode(true)
		case <-m.updateSignal.Done():
			m.updateSignal.Reset()
			m.reconcileAllCVEsInOnlineMode(true)
		}
	}
}

func (m *k8sIstioCVEManagerImpl) HandleClusterConnection() {
	connectionDropThrottle.Run(func() {
		m.updateSignal.Signal()
	})
}

// Update (works only in offline mode) updates new CVEs and reconciles them based on data from scanner bundle
func (m *k8sIstioCVEManagerImpl) Update(zipPath string, forceUpdate bool) {
	if m.mgrMode != offline {
		log.Error("can't fetch in non-offline mode")
		return
	}
	m.reconcileAllCVEsInOfflineMode(zipPath, forceUpdate)
}

// GetNVDCVE returns current istio CVEs loaded in memory
func (m *k8sIstioCVEManagerImpl) GetNVDCVE(id string) *schema.NVDCVEFeedJSON10DefCVEItem {
	cve := m.k8sCVEMgr.getNVDCVE(id)
	if cve == nil {
		return m.istioCVEMgr.getNVDCVE(id)
	}
	return cve
}

func (m *k8sIstioCVEManagerImpl) updateCVEs(nvdCVEs []*schema.NVDCVEFeedJSON10DefCVEItem, ct converter.CVEType) error {
	if ct == converter.K8s {
		return m.k8sCVEMgr.updateCVEs(nvdCVEs)
	} else if ct == converter.Istio {
		return m.istioCVEMgr.updateCVEs(nvdCVEs)
	}
	return errors.Errorf("unknown CVE type: %d", ct)
}

func (m *k8sIstioCVEManagerImpl) reconcileAllCVEsInOnlineMode(forceUpdate bool) {
	if err := m.reconcileOnlineModeCVEs(converter.K8s, forceUpdate); err != nil {
		log.Errorf("reconcile failed for k8s CVEs with error %v", err)
	}
	if err := m.reconcileOnlineModeCVEs(converter.Istio, forceUpdate); err != nil {
		log.Errorf("reconcile failed for istio CVEs with error %v", err)
	}
}

func (m *k8sIstioCVEManagerImpl) reconcileAllCVEsInOfflineMode(zipPath string, forceUpdate bool) {
	if err := m.reconcileOfflineModeCVEs(converter.K8s, zipPath, forceUpdate); err != nil {
		log.Errorf("reconcile failed for k8s CVEs with error %v", err)
	}
	if err := m.reconcileOfflineModeCVEs(converter.Istio, zipPath, forceUpdate); err != nil {
		log.Errorf("reconcile failed for istio CVEs with error %v", err)
	}
}

// reconcileOnlineModeCVEs fetches new CVEs from definitions.stackrox.io and reconciles them
func (m *k8sIstioCVEManagerImpl) reconcileOnlineModeCVEs(ct converter.CVEType, forceUpdate bool) error {
	paths, err := getPaths(ct)
	if err != nil {
		return err
	}

	urls, err := getUrls(ct)
	if err != nil {
		return err
	}

	localCVEChecksum, err := getLocalCVEChecksum(paths.persistentCveChecksumFile)
	if err != nil {
		return nil
	}

	remoteCVEChecksum, err := fetchRemote(urls.cveChecksumURL)
	if err != nil {
		return err
	}

	// If CVEs have been loaded before and checksums are same, no need to update CVEs
	if !forceUpdate && localCVEChecksum == remoteCVEChecksum {
		log.Infof("local and remote CVE checksums are same, skipping download of new %s CVEs", cveTypeToString[ct])
		return nil
	}

	data, err := fetchRemote(urls.cveURL)
	if err != nil {
		return err
	}

	if err := overwriteCVEs(paths.persistentCveFile, paths.persistentCveChecksumFile, remoteCVEChecksum, data); err != nil {
		return err
	}

	newCVEs, err := getLocalCVEs(paths.persistentCveFile)
	if err != nil {
		return err
	}

	if err := m.updateCVEs(newCVEs, ct); err != nil {
		return err
	}

	if localCVEChecksum != remoteCVEChecksum {
		log.Infof("%s CVEs have been updated, %d new CVEs found", cveTypeToString[ct], len(newCVEs))
	}
	return nil
}

// reconcileOfflineModeCVEs reads the scanner bundle zip and updates the CVEs
func (m *k8sIstioCVEManagerImpl) reconcileOfflineModeCVEs(ct converter.CVEType, zipPath string, forceUpdate bool) error {
	paths, err := getPaths(ct)
	if err != nil {
		return err
	}

	bundlePath, err := extractK8sIstioCVEsInScannerBundleZip(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(bundlePath)
		if err != nil {
			log.Errorf("error while deleting the temp bundle dir, error: %v", err)
		}
	}()

	var bundledCVEFile, bundledCVEChecksumFile string

	if ct == converter.K8s {
		bundledCVEFile = filepath.Join(bundlePath, nvd.Feeds[nvd.Kubernetes].CVEFilename)
		bundledCVEChecksumFile = filepath.Join(bundlePath, nvd.Feeds[nvd.Kubernetes].ChecksumFilename)
	} else if ct == converter.Istio {
		bundledCVEFile = filepath.Join(bundlePath, nvd.Feeds[nvd.Istio].CVEFilename)
		bundledCVEChecksumFile = filepath.Join(bundlePath, nvd.Feeds[nvd.Istio].ChecksumFilename)
	} else {
		return errors.Errorf("unknown CVE type: %d", ct)
	}

	oldCveChecksum, err := getLocalCVEChecksum(paths.persistentCveChecksumFile)
	if err != nil {
		return nil
	}

	newCveChecksum, err := getLocalCVEChecksum(bundledCVEChecksumFile)
	if err != nil {
		return err
	}

	// If CVEs have been loaded before and checksums are same, no need to update CVEs
	if !forceUpdate && oldCveChecksum == newCveChecksum {
		log.Infof("local and bundled CVE checksums are same, skipping reconciliation of of new %s CVEs", cveTypeToString[ct])
		return nil
	}

	data, err := ioutil.ReadFile(bundledCVEFile)
	if err != nil {
		return err
	}

	if err := overwriteCVEs(paths.persistentCveFile, paths.persistentCveChecksumFile, newCveChecksum, string(data)); err != nil {
		return err
	}

	newCVEs, err := getLocalCVEs(paths.persistentCveFile)
	if err != nil {
		return err
	}

	if err := m.updateCVEs(newCVEs, ct); err != nil {
		return err
	}

	if oldCveChecksum != newCveChecksum {
		log.Infof("%s CVEs have been updated, %d new CVEs found", cveTypeToString[ct], len(newCVEs))
	}
	return nil
}

func extractK8sIstioCVEsInScannerBundleZip(zipPath string) (string, error) {
	tmpPath, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	if err := unzip(zipPath, tmpPath); err != nil {
		return "", err
	}

	k8sIstioZipPath := filepath.Join(tmpPath, k8sIstioCveZipName)
	if err := unzip(k8sIstioZipPath, tmpPath); err != nil {
		return "", err
	}

	return tmpPath, nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func reconcileCVEsInDB(cveDataStore cveDataStore.DataStore, cveType storage.CVE_CVEType, newCVEs set.StringSet) error {
	results, err := cveDataStore.Search(cveElevatedCtx,
		pkgSearch.NewQueryBuilder().AddExactMatches(pkgSearch.CVEType, cveType.String()).ProtoQuery())
	if err != nil {
		return err
	}

	// Identify the cluster cves that do not affect the infra
	discardCVEs := pkgSearch.ResultsToIDSet(results).Difference(newCVEs)
	if len(discardCVEs) == 0 {
		return nil
	}
	// delete all the cluster cves that do not affect the infra
	return cveDataStore.Delete(cveElevatedCtx, discardCVEs.AsSlice()...)
}
