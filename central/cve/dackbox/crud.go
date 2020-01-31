package dackbox

import (
	"github.com/gogo/protobuf/proto"
	"github.com/stackrox/rox/central/globaldb"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/crud"
)

var (
	// Bucket stores the child image vulnerabilities.
	Bucket = []byte("image_vuln")

	// BucketHandler is the bucket's handler.
	BucketHandler = &badgerhelper.BucketHandler{BucketPrefix: Bucket}

	// Reader reads storage.CVEs directly from the store.
	Reader = crud.NewReader(
		crud.WithAllocFunction(Alloc),
	)

	// Upserter writes storage.CVEs directly to the store.
	Upserter = crud.NewUpserter(crud.WithKeyFunction(KeyFunc))

	// Deleter deletes vulns from the store.
	Deleter = crud.NewDeleter()
)

func init() {
	globaldb.RegisterBucket(Bucket, "Vuln")
}

// KeyFunc returns the key for a CVE.
func KeyFunc(msg proto.Message) []byte {
	unPrefixed := []byte(msg.(interface{ GetId() string }).GetId())
	return badgerhelper.GetBucketKey(Bucket, unPrefixed)
}

// Alloc allocates a CVE.
func Alloc() proto.Message {
	return &storage.CVE{}
}
