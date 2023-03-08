// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/stackrox/rox/generated/storage"
)

// ConvertTestParent1FromProto converts a `*storage.TestParent1` to Gorm model
func ConvertTestParent1FromProto(obj *storage.TestParent1) (*TestParent1, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestParent1{
		Id:         obj.GetId(),
		ParentId:   obj.GetParentId(),
		Val:        obj.GetVal(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestParent1_Child1RefFromProto converts a `*storage.TestParent1_Child1Ref` to Gorm model
func ConvertTestParent1_Child1RefFromProto(obj *storage.TestParent1_Child1Ref, idx int, test_parent1_Id string) (*TestParent1Childrens, error) {
	model := &TestParent1Childrens{
		TestParent1Id: test_parent1_Id,
		Idx:           idx,
		ChildId:       obj.GetChildId(),
	}
	return model, nil
}

// ConvertTestParent1ToProto converts Gorm model `TestParent1` to its protobuf type object
func ConvertTestParent1ToProto(m *TestParent1) (*storage.TestParent1, error) {
	var msg storage.TestParent1
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}