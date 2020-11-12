package handler

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/cve/fetcher"
	"github.com/stackrox/rox/pkg/fileutils"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/utils"
	"google.golang.org/grpc/codes"
)

var (
	scannerDefinitionsSubdir   = path.Join(migrations.DBMountPath, "scannerdefinitions")
	scannerDefinitionsFilePath = path.Join(scannerDefinitionsSubdir, "scanner-defs.zip")
)

const (
	scannerDefsSubZipName = "scanner-defs.zip"
	// K8sIstioCveZipName represent the zip bundle for k8s/istio cves
	K8sIstioCveZipName = "k8s-istio.zip"
)

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		get(w, r)
		return
	}
	if r.Method == http.MethodPost {
		post(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func get(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(scannerDefinitionsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("No scanner definitions found"))
			return
		}
		httputil.WriteGRPCStyleErrorf(w, codes.Internal, "couldn't stat file: %v", err)
		return
	}

	http.ServeFile(w, r, scannerDefinitionsFilePath)
}

func updateK8sIstioCVEs(zipPath string) {
	mgr := fetcher.SingletonManager()
	mgr.Update(zipPath, false)
}

func handleScannerDefsFile(zipF *zip.File) error {
	reader, err := zipF.Open()
	if err != nil {
		return errors.Wrap(err, "opening reader")
	}
	defer utils.IgnoreError(reader.Close)

	err = os.MkdirAll(scannerDefinitionsSubdir, 0755)
	if err != nil {
		return errors.Wrap(err, "creating subdirectory for scanner defs")
	}
	scannerDefsPersistedFile, err := os.Create(scannerDefinitionsFilePath)
	if err != nil {
		return errors.Wrap(err, "creating scanner defs persisted file")
	}
	_, err = io.Copy(scannerDefsPersistedFile, reader)
	if err != nil {
		return errors.Wrap(err, "copying scanner defs zip out")
	}
	err = os.Chtimes(scannerDefinitionsFilePath, time.Now(), zipF.Modified)
	if err != nil {
		return errors.Wrap(err, "changing modified time of scanner defs")
	}
	return nil
}

func handleZipContentsFromOfflineDump(zipPath string) error {
	zipR, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.Wrap(err, "couldn't open file as zip")
	}
	defer utils.IgnoreError(zipR.Close)

	var scannerDefsFileFound bool
	for _, zipF := range zipR.File {
		if zipF.Name == scannerDefsSubZipName {
			if err := handleScannerDefsFile(zipF); err != nil {
				return errors.Wrap(err, "couldn't handle scanner-defs sub file")
			}
			scannerDefsFileFound = true
			continue
		} else if zipF.Name == K8sIstioCveZipName {
			updateK8sIstioCVEs(zipPath)
		}
	}

	if !scannerDefsFileFound {
		return errors.New("scanner defs file not found in upload zip; wrong zip uploaded?")
	}
	return nil
}

func post(w http.ResponseWriter, r *http.Request) {
	tempDir, err := ioutil.TempDir("", "scanner-definitions-handler")
	if err != nil {
		httputil.WriteGRPCStyleErrorf(w, codes.Internal, "failed to create temp dir: %v", err)
		return
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Warnf("Failed to remove temp dir for scanner defs: %v", err)
		}
	}()

	tempFile := filepath.Join(tempDir, "tempfile.zip")
	if err := fileutils.CopySrcToFile(tempFile, r.Body); err != nil {
		httputil.WriteGRPCStyleError(w, codes.Internal, errors.Wrapf(err, "copying HTTP POST body to %s", tempFile))
		return
	}

	if err := handleZipContentsFromOfflineDump(tempFile); err != nil {
		httputil.WriteGRPCStyleError(w, codes.Internal, err)
		return
	}

	_, _ = w.Write([]byte("Successfully stored the offline vulnerability definitions"))
}
