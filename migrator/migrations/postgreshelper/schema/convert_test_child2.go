// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/stackrox/rox/generated/storage"
)

// ConvertTestChild2FromProto converts a `*storage.TestChild2` to Gorm model
func ConvertTestChild2FromProto(obj *storage.TestChild2) (*TestChild2, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestChild2{
		Id:            obj.GetId(),
		ParentId:      obj.GetParentId(),
		GrandparentId: obj.GetGrandparentId(),
		Val:           obj.GetVal(),
		Serialized:    serialized,
	}
	return model, nil
}

// ConvertTestChild2ToProto converts Gorm model `TestChild2` to its protobuf type object
func ConvertTestChild2ToProto(m *TestChild2) (*storage.TestChild2, error) {
	var msg storage.TestChild2
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}