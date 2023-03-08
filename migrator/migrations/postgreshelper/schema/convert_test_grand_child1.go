// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/stackrox/rox/generated/storage"
)

// ConvertTestGrandChild1FromProto converts a `*storage.TestGrandChild1` to Gorm model
func ConvertTestGrandChild1FromProto(obj *storage.TestGrandChild1) (*TestGrandChild1, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestGrandChild1{
		Id:         obj.GetId(),
		ParentId:   obj.GetParentId(),
		ChildId:    obj.GetChildId(),
		Val:        obj.GetVal(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestGrandChild1ToProto converts Gorm model `TestGrandChild1` to its protobuf type object
func ConvertTestGrandChild1ToProto(m *TestGrandChild1) (*storage.TestGrandChild1, error) {
	var msg storage.TestGrandChild1
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}