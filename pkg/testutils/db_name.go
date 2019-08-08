package testutils

import (
	"strings"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/etcd-io/bbolt"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stretchr/testify/require"
)

// Suite abstracts a testify Suite.
type Suite interface {
	T() *testing.T
}

// DBFileName returns an appropriate, unique, DB file name for the given suite.
// It works equally well whether used in SetupSuite or SetupTest.
func DBFileName(suite Suite) string {
	return DBFileNameForT(suite.T())
}

// DBFileNameForT returns an appropriate, unique, DB file name for the given test.
func DBFileNameForT(t *testing.T) string {
	return strings.Replace(t.Name(), "/", "_", -1) + ".db"
}

// DBForT creates and returns a new DB to use for this test, failing the test if creating/opening
// the DB fails.
func DBForT(t *testing.T) *bbolt.DB {
	db, err := bolthelper.NewTemp(DBFileNameForT(t))
	require.NoError(t, err)
	require.NotNil(t, db)
	return db
}

// BadgerDBForT creates and returns a new DB to use for this test, failing the test if creating/opening
// the DB fails.
func BadgerDBForT(t *testing.T) *badger.DB {
	db, _, err := badgerhelper.NewTemp(DBFileNameForT(t))
	require.NoError(t, err)
	require.NotNil(t, db)
	return db
}

// DBForSuite creates and returns a new, temporary DB for use with the given suite.
func DBForSuite(suite Suite) *bbolt.DB {
	return DBForT(suite.T())
}
