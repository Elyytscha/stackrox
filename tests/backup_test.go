package tests

import (
	"archive/zip"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	deploymentBadgerStore "github.com/stackrox/rox/central/deployment/store/badger"
	"github.com/stackrox/rox/central/globaldb/badgerutils"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/tar"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tecbot/gorocksdb"
)

const scratchPath = "backuptest"

// This doesn't actually check for the existence of a specific deployment, but any deployment
func checkDeploymentExists(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)

	service := v1.NewDeploymentServiceClient(conn)

	// 10 seconds should be enough for a deployment to be returned
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		resp, err := service.ListDeployments(ctx, &v1.RawQuery{})
		cancel()
		require.NoError(t, err)
		if len(resp.Deployments) != 0 {
			return
		}
		time.Sleep(2 * time.Second)
	}
	assert.Fail(t, "Failed to find a deployment")
}

// Grab the backup DB and open it, ensuring that there are values for deployments
func TestBackup(t *testing.T) {
	setupNginxLatestTagDeployment(t)
	defer teardownNginxLatestTagDeployment(t)

	checkDeploymentExists(t)

	out, err := os.Create("backup.zip")
	require.NoError(t, err)
	defer func() {
		_ = os.Remove("backup.zip")
	}()

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://"+testutils.RoxAPIEndpoint(t)+"/db/backup", nil)
	require.NoError(t, err)
	req.SetBasicAuth(testutils.RoxUsername(t), testutils.RoxPassword(t))
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer utils.IgnoreError(resp.Body.Close)
	_, err = io.Copy(out, resp.Body)
	require.NoError(t, err)

	defer utils.IgnoreError(out.Close)

	zipFile, err := zip.OpenReader("backup.zip")
	require.NoError(t, err)
	defer utils.IgnoreError(zipFile.Close)

	if features.RocksDB.Enabled() {
		checkZipForRocks(t, zipFile)
	} else {
		checkZipForBadger(t, zipFile)
	}
}

func checkZipForBadger(t *testing.T, zipFile *zip.ReadCloser) {
	badgerFileEntry := getFileWithName(zipFile, "badger.db")
	require.NotNil(t, badgerFileEntry)

	badgerFile, err := badgerFileEntry.Open()
	require.NoError(t, err)

	b, err := badgerhelper.New("backup.db")
	require.NoError(t, err)

	require.NoError(t, badgerutils.Load(badgerFile, b))

	depStore, err := deploymentBadgerStore.New(b)
	require.NoError(t, err)
	deployments, err := depStore.GetDeployments()
	require.NoError(t, err)
	assert.NotEmpty(t, deployments)
}

func checkZipForRocks(t *testing.T, zipFile *zip.ReadCloser) {
	// Open the tar file holding the rocks DB backup.
	rocksFileEntry := getFileWithName(zipFile, "rocks.db")
	require.NotNil(t, rocksFileEntry)
	rocksFile, err := rocksFileEntry.Open()
	require.NoError(t, err)

	// Dump the untar'd rocks file to a scratch directory.
	tmpBackupDir, err := ioutil.TempDir("", scratchPath)
	require.NoError(t, err)
	defer utils.IgnoreError(func() error { return os.RemoveAll(tmpBackupDir) })

	err = tar.ToPath(tmpBackupDir, rocksFile)
	require.NoError(t, err)
	require.NoError(t, rocksFile.Close())

	// Generate the backup files in the directory.
	opts := gorocksdb.NewDefaultOptions()
	backupEngine, err := gorocksdb.OpenBackupEngine(opts, tmpBackupDir)
	require.NoError(t, err)

	// Restore the db to another temp directory
	tmpDBDir, err := ioutil.TempDir("", scratchPath)
	require.NoError(t, err)
	defer utils.IgnoreError(func() error { return os.RemoveAll(tmpDBDir) })
	err = backupEngine.RestoreDBFromLatestBackup(tmpDBDir, tmpDBDir, gorocksdb.NewRestoreOptions())
	require.NoError(t, err)

	// Check for errors on cleanup.
	require.NoError(t, os.RemoveAll(tmpBackupDir))
	require.NoError(t, os.RemoveAll(tmpDBDir))
}

func getFileWithName(zipFile *zip.ReadCloser, name string) *zip.File {
	var ret *zip.File
	for _, f := range zipFile.File {
		if f.Name == name {
			ret = f
			break
		}
	}
	return ret
}
