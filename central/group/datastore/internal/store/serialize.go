package store

import (
	"fmt"
	"strings"

	serializePkg "github.com/stackrox/rox/central/group/datastore/serialize"
	"github.com/stackrox/rox/generated/storage"
)

// Serialization
////////////////

func serialize(group *storage.Group) ([]byte, []byte) {
	return serializeKey(group.GetProps()), serializeValue(group)
}

func serializeKey(props *storage.GroupProperties) []byte {
	return []byte(serializePkg.PropsKey(props))
}

func serializeValue(group *storage.Group) []byte {
	return []byte(group.GetRoleName())
}

// Deserialization
////////////////

func deserialize(key, value []byte) (*storage.Group, error) {
	props, err := deserializeKey(key)
	if err != nil {
		return nil, err
	}
	return &storage.Group{
		Props:    props,
		RoleName: string(value),
	}, nil
}

func deserializeKey(key []byte) (*storage.GroupProperties, error) {
	str := string(key)
	props := strings.Split(str, serializePkg.KeySeparator)
	if len(props) != 3 {
		return nil, fmt.Errorf("unable to deserialize key: %s", str)
	}
	// If no values (which is totally ok), then just return nil.
	if len(props[0]) == 0 && len(props[1]) == 0 && len(props[2]) == 0 {
		return nil, nil
	}
	return &storage.GroupProperties{
		AuthProviderId: props[0],
		Key:            props[1],
		Value:          props[2],
	}, nil
}
