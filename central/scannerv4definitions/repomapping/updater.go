package repomapping

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/scannerdefinitions/file"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
)

type repoMappingUpdater struct {
	files []*file.File

	client      *http.Client
	downloadURL string
	interval    time.Duration
	once        sync.Once
	stopSig     concurrency.Signal
}

const (
	baseURL = "https://storage.googleapis.com/scanner-v4-test/redhat-repository-mappings/"
)

// NewUpdater creates a new updater.
func NewUpdater(files []*file.File, client *http.Client, downloadURL string, interval time.Duration) *repoMappingUpdater {
	return &repoMappingUpdater{
		files:       files,
		client:      client,
		downloadURL: downloadURL,
		interval:    interval,
		stopSig:     concurrency.NewSignal(),
	}
}

// Stop stops the updater.
func (u *repoMappingUpdater) Stop() {
	u.stopSig.Signal()
}

// Start starts the updater.
// The updater is only started once.
func (u *repoMappingUpdater) Start() {
	u.once.Do(func() {
		// Run the first update in a blocking-manner.
		u.update()
		go u.runForever()
	})
}

func (u *repoMappingUpdater) runForever() {
	t := time.NewTicker(u.interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			u.update()
		case <-u.stopSig.Done():
			return
		}
	}
}

func (u *repoMappingUpdater) update() error {
	if err := u.doUpdate(); err != nil {
		log.Errorf("Failed to update Scanner v4 repository mapping from endpoint %q: %v", u.downloadURL, err)
		return err
	}
	return nil
}

func (u *repoMappingUpdater) doUpdate() error {
	tempDir, err := os.MkdirTemp("", "repomapping")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filesToDownload := []string{name2Cpe, repo2Cpe}
	for i, file := range filesToDownload {
		if i >= len(u.files) {
			return errors.New("Insufficient files available to store the downloaded JSON.")
		}

		out, err := downloadFromURL(baseURL+file, tempDir, file)
		if err != nil {
			return fmt.Errorf("failed to download %s: %v", file, err)
		}
		// Seek to the beginning of the outZip
		_, err = out.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("error seeking to the beginning of zip outZip: %w", err)
		}
		err = u.files[i].WriteContent(out)
		if err != nil {
			return fmt.Errorf("error writing file to content: %w", err)
		}
		log.Infof("Successfully write to the content: %v", u.files[i].Path())
	}

	return nil
}

func downloadFromURL(url, dir, filename string) (*os.File, error) {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second * 3)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK { // Success
			out, err := os.Create(filepath.Join(dir, filename))
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return nil, err
			}

			return out, nil
		} else {
			time.Sleep(time.Second * 3)
			continue
		}
	}
	return nil, lastErr
}
