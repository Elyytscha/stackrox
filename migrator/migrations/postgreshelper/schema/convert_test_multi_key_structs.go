// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/lib/pq"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

// ConvertTestMultiKeyStructFromProto converts a `*storage.TestMultiKeyStruct` to Gorm model
func ConvertTestMultiKeyStructFromProto(obj *storage.TestMultiKeyStruct) (*TestMultiKeyStructs, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestMultiKeyStructs{
		Key1:              obj.GetKey1(),
		Key2:              obj.GetKey2(),
		StringSlice:       pq.Array(obj.GetStringSlice()).(*pq.StringArray),
		Bool:              obj.GetBool(),
		Uint64:            obj.GetUint64(),
		Int64:             obj.GetInt64(),
		Float:             obj.GetFloat(),
		Labels:            obj.GetLabels(),
		Timestamp:         pgutils.NilOrTime(obj.GetTimestamp()),
		Enum:              obj.GetEnum(),
		Enums:             pq.Array(pgutils.ConvertEnumSliceToIntArray(obj.GetEnums())).(*pq.Int32Array),
		String:            obj.GetString_(),
		Int32Slice:        pq.Array(obj.GetInt32Slice()).(*pq.Int32Array),
		OneofnestedNested: obj.GetOneofnested().GetNested(),
		Serialized:        serialized,
	}
	return model, nil
}

// ConvertTestMultiKeyStruct_NestedFromProto converts a `*storage.TestMultiKeyStruct_Nested` to Gorm model
func ConvertTestMultiKeyStruct_NestedFromProto(obj *storage.TestMultiKeyStruct_Nested, idx int, test_multi_key_structs_Key1 string, test_multi_key_structs_Key2 string) (*TestMultiKeyStructsNesteds, error) {
	model := &TestMultiKeyStructsNesteds{
		TestMultiKeyStructsKey1: test_multi_key_structs_Key1,
		TestMultiKeyStructsKey2: test_multi_key_structs_Key2,
		Idx:                     idx,
		Nested:                  obj.GetNested(),
		IsNested:                obj.GetIsNested(),
		Int64:                   obj.GetInt64(),
		Nested2Nested2:          obj.GetNested2().GetNested2(),
		Nested2IsNested:         obj.GetNested2().GetIsNested(),
		Nested2Int64:            obj.GetNested2().GetInt64(),
	}
	return model, nil
}

// ConvertTestMultiKeyStructToProto converts Gorm model `TestMultiKeyStructs` to its protobuf type object
func ConvertTestMultiKeyStructToProto(m *TestMultiKeyStructs) (*storage.TestMultiKeyStruct, error) {
	var msg storage.TestMultiKeyStruct
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}