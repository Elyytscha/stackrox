package store

import (
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/bolthelper"
	"github.com/boltdb/bolt"
)

const (
	apiTokensBucket = "apiTokens"
)

// Store is the (bolt-backed) store for API tokens.
// We don't store the tokens themselves, but do store metadata.
// Importantly, the Store persists token revocations.
type Store interface {
	AddToken(*v1.TokenMetadata) error
	GetToken(id string) (token *v1.TokenMetadata, exists bool, err error)
	GetTokens(*v1.GetAPITokensRequest) ([]*v1.TokenMetadata, error)
	RevokeToken(id string) (exists bool, err error)
}

// New returns a ready-to-use store.
func New(db *bolt.DB) Store {
	bolthelper.RegisterBucketOrPanic(db, apiTokensBucket)
	return &storeImpl{DB: db}
}
