package rocksdb

import (
	"github.com/stackrox/rox/pkg/dackbox/transactions"
	generic "github.com/stackrox/rox/pkg/rocksdb/crud"
	"github.com/stackrox/rox/pkg/sliceutils"
	"github.com/tecbot/gorocksdb"
)

type rocksDBWrapper struct {
	db *gorocksdb.DB
}

func (b *rocksDBWrapper) NewTransaction(update bool) transactions.DBTransaction {
	snapshot := b.db.NewSnapshot()
	readOpts := gorocksdb.NewDefaultReadOptions()
	readOpts.SetSnapshot(snapshot)

	itOpts := gorocksdb.NewDefaultReadOptions()
	itOpts.SetSnapshot(snapshot)
	itOpts.SetPrefixSameAsStart(true)
	itOpts.SetFillCache(false)

	wrapper := &txnWrapper{
		db:       b.db,
		isUpdate: update,

		snapshot: snapshot,
		readOpts: readOpts,
		itOpts:   itOpts,
	}
	if update {
		wrapper.batch = gorocksdb.NewWriteBatch()
	}
	return wrapper
}

// NewRocksDBWrapper is a wrapper around a rocksDB so it implements the DBTransactionFactory interface
func NewRocksDBWrapper(db *gorocksdb.DB) transactions.DBTransactionFactory {
	return &rocksDBWrapper{
		db: db,
	}
}

type txnWrapper struct {
	db       *gorocksdb.DB
	isUpdate bool

	batch *gorocksdb.WriteBatch

	readOpts *gorocksdb.ReadOptions
	itOpts   *gorocksdb.ReadOptions
	snapshot *gorocksdb.Snapshot
}

func (t *txnWrapper) Delete(keys ...[]byte) error {
	if !t.isUpdate {
		panic("trying to delete a key during a read txn")
	}
	for _, k := range keys {
		t.batch.Delete(k)
	}
	return nil
}

func (t *txnWrapper) Get(key []byte) ([]byte, bool, error) {
	slice, err := t.db.Get(t.readOpts, key)
	if err != nil {
		return nil, false, err
	}
	defer slice.Free()
	if !slice.Exists() {
		return nil, false, nil
	}
	// Copy before returning as the slice is freed in defer
	return sliceutils.ByteClone(slice.Data()), true, nil
}

func (t *txnWrapper) Set(key, value []byte) error {
	if !t.isUpdate {
		panic("trying to set during a read txn")
	}

	t.batch.Put(key, value)
	return nil
}

func (t *txnWrapper) BucketForEach(graphPrefix []byte, stripPrefix bool, fn func(k, v []byte) error) error {
	return generic.BucketForEach(t.db, t.itOpts, graphPrefix, stripPrefix, fn)
}

func (t *txnWrapper) BucketKeyForEach(graphPrefix []byte, stripPrefix bool, fn func(k []byte) error) error {
	return generic.BucketKeyForEach(t.db, t.itOpts, graphPrefix, stripPrefix, fn)
}

func (t *txnWrapper) BucketKeyCount(prefix []byte) (int, error) {
	var count int
	err := generic.BucketKeyForEach(t.db, t.itOpts, prefix, false, func(k []byte) error {
		count++
		return nil
	})
	return count, err
}

func (t *txnWrapper) Commit() error {
	writeOpts := generic.DefaultWriteOptions()
	defer writeOpts.Destroy()

	return t.db.Write(writeOpts, t.batch)
}

func (t *txnWrapper) Discard() {
	if t.batch != nil {
		t.batch.Destroy()
		t.batch = nil
	}
	if t.readOpts != nil {
		t.readOpts.Destroy()
		t.readOpts = nil
	}
	if t.itOpts != nil {
		t.itOpts.Destroy()
		t.itOpts = nil
	}
	if t.snapshot != nil {
		t.db.ReleaseSnapshot(t.snapshot)
		t.snapshot = nil
	}
}
