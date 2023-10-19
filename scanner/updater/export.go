package updater

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/quay/claircore/libvuln/driver"
	"github.com/quay/claircore/libvuln/jsonblob"
	"github.com/quay/claircore/libvuln/updates"
	"github.com/quay/zlog"
	"github.com/stackrox/rox/scanner/updater/manual"
	"golang.org/x/time/rate"

	// default updaters
	_ "github.com/quay/claircore/updater/defaults"
)

// Export is responsible for triggering the updaters to download Common Vulnerabilities and Exposures (CVEs) data
// and then outputting the result as a zstd-compressed file with .zst extension
func Export(ctx context.Context, outputDir string) error {

	err := os.MkdirAll(outputDir, 0700)
	if err != nil {
		return err
	}
	// create output json file
	outputFile, err := os.Create(filepath.Join(outputDir, "output.zst"))
	if err != nil {
		return err
	}

	limiter := rate.NewLimiter(rate.Every(time.Second), 5)
	httpClient := &http.Client{
		Transport: &rateLimitedTransport{
			limiter:   limiter,
			transport: http.DefaultTransport,
		},
	}

	zstdWriter, err := zstd.NewWriter(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := zstdWriter.Close()
		if closeErr != nil {
			zlog.Error(ctx).Err(closeErr).Msg("Failed to close zstd writer")
		}
	}()

	updaterSet, err := manual.UpdaterSet(ctx, nil)
	if err != nil {
		return err
	}
	outOfTree := [][]driver.Updater{
		make([]driver.Updater, 0),
	}
	outOfTree = append(outOfTree, updaterSet.Updaters())

	for i, uSet := range [][]string{
		{"oracle", "photon", "suse", "osv", "rhcc"},
		{"alpine", "rhel", "ubuntu", "aws", "debian"},
	} {
		jsonStore, err := jsonblob.New()
		if err != nil {
			return err
		}

		updateMgr, err := updates.NewManager(ctx, jsonStore, updates.NewLocalLockSource(), httpClient,
			updates.WithEnabled(uSet),
			updates.WithOutOfTree(outOfTree[i]),
		)
		if err != nil {
			return err
		}

		if err := updateMgr.Run(ctx); err != nil {
			return err
		}

		err = jsonStore.Store(zstdWriter)
		if err != nil {
			return err
		}

		// Flush the zstd writer to make sure compressed data is sent to outputFile
		if err := zstdWriter.Flush(); err != nil {
			return err
		}

		// Ensure data is flushed from OS buffers to the disk
		if err := outputFile.Sync(); err != nil {
			return err
		}
	}

	return nil
}

type rateLimitedTransport struct {
	limiter   *rate.Limiter
	transport http.RoundTripper
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req)
}
