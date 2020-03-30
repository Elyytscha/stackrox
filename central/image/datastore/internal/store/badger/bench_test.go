package badger

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/central/image/datastore/internal/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/require"
)

const maxGRPCSize = 4194304

func getImageStore(b *testing.B) store.Store {
	db, _, err := badgerhelper.NewTemp(b.Name() + ".db")
	if err != nil {
		b.Fatal(err)
	}
	return New(db, false)
}

func BenchmarkAddImage(b *testing.B) {
	store := getImageStore(b)
	image := fixtures.GetImage()
	for i := 0; i < b.N; i++ {
		require.NoError(b, store.Upsert(image))
	}
}

func BenchmarkGetImage(b *testing.B) {
	store := getImageStore(b)
	image := fixtures.GetImage()
	require.NoError(b, store.Upsert(image))
	for i := 0; i < b.N; i++ {
		_, exists, err := store.GetImage(image.GetId(), true)
		require.True(b, exists)
		require.NoError(b, err)
	}
}

func BenchmarkListImage(b *testing.B) {
	store := getImageStore(b)
	image := fixtures.GetImage()
	require.NoError(b, store.Upsert(image))
	for i := 0; i < b.N; i++ {
		_, exists, err := store.ListImage(image.GetId())
		require.True(b, exists)
		require.NoError(b, err)
	}
}

// This really isn't a benchmark, but just prints out how many ListImages can be returned in an API call
func BenchmarkMaxListImage(b *testing.B) {
	listImage := &storage.ListImage{
		Id:   uuid.NewDummy().String(),
		Name: "quizzical_cat",
		SetComponents: &storage.ListImage_Components{
			Components: 10,
		},
		SetCves: &storage.ListImage_Cves{
			Cves: 10,
		},
		SetFixable: &storage.ListImage_FixableCves{
			FixableCves: 10,
		},
		Created: types.TimestampNow(),
	}

	bytes, _ := proto.Marshal(listImage)
	fmt.Printf("Max ListImages that can be returned: %d\n", maxGRPCSize/len(bytes))
}
