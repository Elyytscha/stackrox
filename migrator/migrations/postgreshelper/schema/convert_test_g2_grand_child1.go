// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/stackrox/rox/generated/storage"
)

// ConvertTestG2GrandChild1FromProto converts a `*storage.TestG2GrandChild1` to Gorm model
func ConvertTestG2GrandChild1FromProto(obj *storage.TestG2GrandChild1) (*TestG2GrandChild1, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestG2GrandChild1{
		Id:         obj.GetId(),
		ParentId:   obj.GetParentId(),
		ChildId:    obj.GetChildId(),
		Val:        obj.GetVal(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestG2GrandChild1ToProto converts Gorm model `TestG2GrandChild1` to its protobuf type object
func ConvertTestG2GrandChild1ToProto(m *TestG2GrandChild1) (*storage.TestG2GrandChild1, error) {
	var msg storage.TestG2GrandChild1
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}