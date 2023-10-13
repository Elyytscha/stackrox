// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/auth_machine_to_machine.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AuthMachineToMachineConfig_Type int32

const (
	AuthMachineToMachineConfig_GITHUB_ACTIONS AuthMachineToMachineConfig_Type = 0
	AuthMachineToMachineConfig_GENERIC        AuthMachineToMachineConfig_Type = 1
)

var AuthMachineToMachineConfig_Type_name = map[int32]string{
	0: "GITHUB_ACTIONS",
	1: "GENERIC",
}

var AuthMachineToMachineConfig_Type_value = map[string]int32{
	"GITHUB_ACTIONS": 0,
	"GENERIC":        1,
}

func (x AuthMachineToMachineConfig_Type) String() string {
	return proto.EnumName(AuthMachineToMachineConfig_Type_name, int32(x))
}

func (AuthMachineToMachineConfig_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5f143e8b019d79, []int{0, 0}
}

// AuthMachineToMachineConfig is the storage representation of auth machine to machine configs in Central.
//
// Refer to v1.AuthMachineToMachineConfig for a more detailed doc.
// Next tag: 6.
type AuthMachineToMachineConfig struct {
	Id                      string                                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk,type(uuid)"`
	Type                    AuthMachineToMachineConfig_Type       `protobuf:"varint,2,opt,name=type,proto3,enum=AuthMachineToMachineConfig_Type" json:"type,omitempty"`
	TokenExpirationDuration string                                `protobuf:"bytes,3,opt,name=token_expiration_duration,json=tokenExpirationDuration,proto3" json:"token_expiration_duration,omitempty"`
	Mappings                []*AuthMachineToMachineConfig_Mapping `protobuf:"bytes,4,rep,name=mappings,proto3" json:"mappings,omitempty"`
	// Types that are valid to be assigned to IssuerConfig:
	//	*AuthMachineToMachineConfig_Generic
	IssuerConfig         isAuthMachineToMachineConfig_IssuerConfig `protobuf_oneof:"IssuerConfig"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *AuthMachineToMachineConfig) Reset()         { *m = AuthMachineToMachineConfig{} }
func (m *AuthMachineToMachineConfig) String() string { return proto.CompactTextString(m) }
func (*AuthMachineToMachineConfig) ProtoMessage()    {}
func (*AuthMachineToMachineConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5f143e8b019d79, []int{0}
}
func (m *AuthMachineToMachineConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthMachineToMachineConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthMachineToMachineConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthMachineToMachineConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthMachineToMachineConfig.Merge(m, src)
}
func (m *AuthMachineToMachineConfig) XXX_Size() int {
	return m.Size()
}
func (m *AuthMachineToMachineConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthMachineToMachineConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AuthMachineToMachineConfig proto.InternalMessageInfo

type isAuthMachineToMachineConfig_IssuerConfig interface {
	isAuthMachineToMachineConfig_IssuerConfig()
	MarshalTo([]byte) (int, error)
	Size() int
	Clone() isAuthMachineToMachineConfig_IssuerConfig
}

type AuthMachineToMachineConfig_Generic struct {
	Generic *AuthMachineToMachineConfig_GenericIssuer `protobuf:"bytes,5,opt,name=generic,proto3,oneof" json:"generic,omitempty"`
}

func (*AuthMachineToMachineConfig_Generic) isAuthMachineToMachineConfig_IssuerConfig() {}
func (m *AuthMachineToMachineConfig_Generic) Clone() isAuthMachineToMachineConfig_IssuerConfig {
	if m == nil {
		return nil
	}
	cloned := new(AuthMachineToMachineConfig_Generic)
	*cloned = *m

	cloned.Generic = m.Generic.Clone()
	return cloned
}

func (m *AuthMachineToMachineConfig) GetIssuerConfig() isAuthMachineToMachineConfig_IssuerConfig {
	if m != nil {
		return m.IssuerConfig
	}
	return nil
}

func (m *AuthMachineToMachineConfig) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AuthMachineToMachineConfig) GetType() AuthMachineToMachineConfig_Type {
	if m != nil {
		return m.Type
	}
	return AuthMachineToMachineConfig_GITHUB_ACTIONS
}

func (m *AuthMachineToMachineConfig) GetTokenExpirationDuration() string {
	if m != nil {
		return m.TokenExpirationDuration
	}
	return ""
}

func (m *AuthMachineToMachineConfig) GetMappings() []*AuthMachineToMachineConfig_Mapping {
	if m != nil {
		return m.Mappings
	}
	return nil
}

func (m *AuthMachineToMachineConfig) GetGeneric() *AuthMachineToMachineConfig_GenericIssuer {
	if x, ok := m.GetIssuerConfig().(*AuthMachineToMachineConfig_Generic); ok {
		return x.Generic
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AuthMachineToMachineConfig) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AuthMachineToMachineConfig_Generic)(nil),
	}
}

func (m *AuthMachineToMachineConfig) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthMachineToMachineConfig) Clone() *AuthMachineToMachineConfig {
	if m == nil {
		return nil
	}
	cloned := new(AuthMachineToMachineConfig)
	*cloned = *m

	if m.Mappings != nil {
		cloned.Mappings = make([]*AuthMachineToMachineConfig_Mapping, len(m.Mappings))
		for idx, v := range m.Mappings {
			cloned.Mappings[idx] = v.Clone()
		}
	}
	if m.IssuerConfig != nil {
		cloned.IssuerConfig = m.IssuerConfig.Clone()
	}
	return cloned
}

type AuthMachineToMachineConfig_Mapping struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Role                 string   `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthMachineToMachineConfig_Mapping) Reset()         { *m = AuthMachineToMachineConfig_Mapping{} }
func (m *AuthMachineToMachineConfig_Mapping) String() string { return proto.CompactTextString(m) }
func (*AuthMachineToMachineConfig_Mapping) ProtoMessage()    {}
func (*AuthMachineToMachineConfig_Mapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5f143e8b019d79, []int{0, 0}
}
func (m *AuthMachineToMachineConfig_Mapping) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthMachineToMachineConfig_Mapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthMachineToMachineConfig_Mapping.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthMachineToMachineConfig_Mapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthMachineToMachineConfig_Mapping.Merge(m, src)
}
func (m *AuthMachineToMachineConfig_Mapping) XXX_Size() int {
	return m.Size()
}
func (m *AuthMachineToMachineConfig_Mapping) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthMachineToMachineConfig_Mapping.DiscardUnknown(m)
}

var xxx_messageInfo_AuthMachineToMachineConfig_Mapping proto.InternalMessageInfo

func (m *AuthMachineToMachineConfig_Mapping) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *AuthMachineToMachineConfig_Mapping) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *AuthMachineToMachineConfig_Mapping) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

func (m *AuthMachineToMachineConfig_Mapping) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthMachineToMachineConfig_Mapping) Clone() *AuthMachineToMachineConfig_Mapping {
	if m == nil {
		return nil
	}
	cloned := new(AuthMachineToMachineConfig_Mapping)
	*cloned = *m

	return cloned
}

type AuthMachineToMachineConfig_GenericIssuer struct {
	Issuer               string   `protobuf:"bytes,1,opt,name=issuer,proto3" json:"issuer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthMachineToMachineConfig_GenericIssuer) Reset() {
	*m = AuthMachineToMachineConfig_GenericIssuer{}
}
func (m *AuthMachineToMachineConfig_GenericIssuer) String() string { return proto.CompactTextString(m) }
func (*AuthMachineToMachineConfig_GenericIssuer) ProtoMessage()    {}
func (*AuthMachineToMachineConfig_GenericIssuer) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5f143e8b019d79, []int{0, 1}
}
func (m *AuthMachineToMachineConfig_GenericIssuer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthMachineToMachineConfig_GenericIssuer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthMachineToMachineConfig_GenericIssuer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthMachineToMachineConfig_GenericIssuer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthMachineToMachineConfig_GenericIssuer.Merge(m, src)
}
func (m *AuthMachineToMachineConfig_GenericIssuer) XXX_Size() int {
	return m.Size()
}
func (m *AuthMachineToMachineConfig_GenericIssuer) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthMachineToMachineConfig_GenericIssuer.DiscardUnknown(m)
}

var xxx_messageInfo_AuthMachineToMachineConfig_GenericIssuer proto.InternalMessageInfo

func (m *AuthMachineToMachineConfig_GenericIssuer) GetIssuer() string {
	if m != nil {
		return m.Issuer
	}
	return ""
}

func (m *AuthMachineToMachineConfig_GenericIssuer) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthMachineToMachineConfig_GenericIssuer) Clone() *AuthMachineToMachineConfig_GenericIssuer {
	if m == nil {
		return nil
	}
	cloned := new(AuthMachineToMachineConfig_GenericIssuer)
	*cloned = *m

	return cloned
}

func init() {
	proto.RegisterEnum("AuthMachineToMachineConfig_Type", AuthMachineToMachineConfig_Type_name, AuthMachineToMachineConfig_Type_value)
	proto.RegisterType((*AuthMachineToMachineConfig)(nil), "AuthMachineToMachineConfig")
	proto.RegisterType((*AuthMachineToMachineConfig_Mapping)(nil), "AuthMachineToMachineConfig.Mapping")
	proto.RegisterType((*AuthMachineToMachineConfig_GenericIssuer)(nil), "AuthMachineToMachineConfig.GenericIssuer")
}

func init() {
	proto.RegisterFile("storage/auth_machine_to_machine.proto", fileDescriptor_ce5f143e8b019d79)
}

var fileDescriptor_ce5f143e8b019d79 = []byte{
	// 407 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xb3, 0xb1, 0x5b, 0xd3, 0x09, 0x44, 0xd1, 0x52, 0x51, 0xd7, 0x07, 0x63, 0x19, 0xa1,
	0xb8, 0x12, 0x72, 0xa5, 0xd0, 0x53, 0x2f, 0xa8, 0x0e, 0x56, 0xea, 0x43, 0x8b, 0xb4, 0x98, 0x0b,
	0x17, 0xcb, 0xc4, 0x8b, 0xb3, 0x72, 0xea, 0x35, 0xf6, 0x1a, 0x35, 0x6f, 0xc2, 0x23, 0x71, 0xe4,
	0x09, 0x10, 0x0a, 0x17, 0xce, 0x3c, 0x01, 0xb2, 0x77, 0x5b, 0x89, 0x43, 0x73, 0xfb, 0x66, 0xe6,
	0x9f, 0x7f, 0x34, 0xa3, 0x81, 0x97, 0x8d, 0xe0, 0x75, 0x9a, 0xd3, 0xd3, 0xb4, 0x15, 0xab, 0xe4,
	0x26, 0x5d, 0xae, 0x58, 0x49, 0x13, 0xc1, 0xef, 0xd0, 0xaf, 0x6a, 0x2e, 0xb8, 0x75, 0x98, 0xf3,
	0x9c, 0xf7, 0x78, 0xda, 0x91, 0xcc, 0xba, 0x7f, 0x34, 0xb0, 0x2e, 0x5a, 0xb1, 0xba, 0x92, 0xda,
	0x98, 0x2b, 0x98, 0xf3, 0xf2, 0x33, 0xcb, 0xf1, 0x14, 0x86, 0x2c, 0x33, 0x91, 0x83, 0xbc, 0x83,
	0xe0, 0xe8, 0xef, 0xcf, 0xe7, 0x4f, 0x9b, 0x2f, 0xeb, 0x73, 0xb7, 0x2a, 0x5e, 0x89, 0x4d, 0x45,
	0xbd, 0xb6, 0x65, 0xd9, 0x89, 0x4b, 0x86, 0x2c, 0xc3, 0x67, 0xa0, 0x77, 0x29, 0x73, 0xe8, 0x20,
	0x6f, 0x3c, 0x73, 0xfc, 0x87, 0x3d, 0xfd, 0x78, 0x53, 0x51, 0xd2, 0xab, 0xf1, 0x39, 0x1c, 0x0b,
	0x5e, 0xd0, 0x32, 0xa1, 0xb7, 0x15, 0xab, 0x53, 0xc1, 0x78, 0x99, 0x64, 0xad, 0x04, 0x53, 0xeb,
	0xa6, 0x92, 0xa3, 0x5e, 0x10, 0xde, 0xd7, 0xdf, 0xaa, 0x32, 0x7e, 0x03, 0x8f, 0x6e, 0xd2, 0xaa,
	0x62, 0x65, 0xde, 0x98, 0xba, 0xa3, 0x79, 0xa3, 0xd9, 0x8b, 0x5d, 0x53, 0xaf, 0xa4, 0x96, 0xdc,
	0x37, 0xe1, 0x10, 0x8c, 0x9c, 0x96, 0xb4, 0x66, 0x4b, 0x73, 0xcf, 0x41, 0xde, 0x68, 0x76, 0xb2,
	0xab, 0x7f, 0x21, 0xa5, 0x51, 0xd3, 0xb4, 0xb4, 0xbe, 0x1c, 0x90, 0xbb, 0x5e, 0x2b, 0x04, 0x43,
	0x79, 0xe3, 0x09, 0x68, 0x05, 0xdd, 0xc8, 0x73, 0x91, 0x0e, 0xf1, 0x21, 0xec, 0x7d, 0x4d, 0xd7,
	0xad, 0xbc, 0xcb, 0x01, 0x91, 0x01, 0xc6, 0xa0, 0xd7, 0x7c, 0x4d, 0xd5, 0x86, 0x3d, 0x5b, 0x53,
	0x78, 0xf2, 0xdf, 0x08, 0xfc, 0x0c, 0xf6, 0x59, 0x4f, 0xca, 0x4f, 0x45, 0xee, 0x14, 0xf4, 0xee,
	0x82, 0x18, 0xc3, 0x78, 0x11, 0xc5, 0x97, 0x1f, 0x82, 0xe4, 0x62, 0x1e, 0x47, 0xef, 0xae, 0xdf,
	0x4f, 0x06, 0x78, 0x04, 0xc6, 0x22, 0xbc, 0x0e, 0x49, 0x34, 0x9f, 0xa0, 0x60, 0x0c, 0x8f, 0xa5,
	0x95, 0xdc, 0x20, 0x38, 0xfb, 0xbe, 0xb5, 0xd1, 0x8f, 0xad, 0x8d, 0x7e, 0x6d, 0x6d, 0xf4, 0xed,
	0xb7, 0x3d, 0x80, 0x63, 0xc6, 0xfd, 0x46, 0xa4, 0xcb, 0xa2, 0xe6, 0xb7, 0xf2, 0x1f, 0x7c, 0xf5,
	0x4b, 0x1f, 0x0d, 0x05, 0x9f, 0xf6, 0xfb, 0xfc, 0xeb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3b,
	0x58, 0xcc, 0x24, 0x66, 0x02, 0x00, 0x00,
}

func (m *AuthMachineToMachineConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthMachineToMachineConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthMachineToMachineConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.IssuerConfig != nil {
		{
			size := m.IssuerConfig.Size()
			i -= size
			if _, err := m.IssuerConfig.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if len(m.Mappings) > 0 {
		for iNdEx := len(m.Mappings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Mappings[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.TokenExpirationDuration) > 0 {
		i -= len(m.TokenExpirationDuration)
		copy(dAtA[i:], m.TokenExpirationDuration)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.TokenExpirationDuration)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Type != 0 {
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthMachineToMachineConfig_Generic) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthMachineToMachineConfig_Generic) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Generic != nil {
		{
			size, err := m.Generic.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	return len(dAtA) - i, nil
}
func (m *AuthMachineToMachineConfig_Mapping) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthMachineToMachineConfig_Mapping) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthMachineToMachineConfig_Mapping) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Role) > 0 {
		i -= len(m.Role)
		copy(dAtA[i:], m.Role)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.Role)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AuthMachineToMachineConfig_GenericIssuer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthMachineToMachineConfig_GenericIssuer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthMachineToMachineConfig_GenericIssuer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Issuer) > 0 {
		i -= len(m.Issuer)
		copy(dAtA[i:], m.Issuer)
		i = encodeVarintAuthMachineToMachine(dAtA, i, uint64(len(m.Issuer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAuthMachineToMachine(dAtA []byte, offset int, v uint64) int {
	offset -= sovAuthMachineToMachine(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AuthMachineToMachineConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovAuthMachineToMachine(uint64(m.Type))
	}
	l = len(m.TokenExpirationDuration)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	if len(m.Mappings) > 0 {
		for _, e := range m.Mappings {
			l = e.Size()
			n += 1 + l + sovAuthMachineToMachine(uint64(l))
		}
	}
	if m.IssuerConfig != nil {
		n += m.IssuerConfig.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *AuthMachineToMachineConfig_Generic) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Generic != nil {
		l = m.Generic.Size()
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	return n
}
func (m *AuthMachineToMachineConfig_Mapping) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	l = len(m.Role)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *AuthMachineToMachineConfig_GenericIssuer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Issuer)
	if l > 0 {
		n += 1 + l + sovAuthMachineToMachine(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAuthMachineToMachine(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAuthMachineToMachine(x uint64) (n int) {
	return sovAuthMachineToMachine(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AuthMachineToMachineConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthMachineToMachine
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AuthMachineToMachineConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AuthMachineToMachineConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= AuthMachineToMachineConfig_Type(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenExpirationDuration", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenExpirationDuration = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mappings", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Mappings = append(m.Mappings, &AuthMachineToMachineConfig_Mapping{})
			if err := m.Mappings[len(m.Mappings)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Generic", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &AuthMachineToMachineConfig_GenericIssuer{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.IssuerConfig = &AuthMachineToMachineConfig_Generic{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthMachineToMachine(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AuthMachineToMachineConfig_Mapping) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthMachineToMachine
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Mapping: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Mapping: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Role", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Role = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthMachineToMachine(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AuthMachineToMachineConfig_GenericIssuer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthMachineToMachine
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenericIssuer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenericIssuer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Issuer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Issuer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAuthMachineToMachine(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthMachineToMachine
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAuthMachineToMachine(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuthMachineToMachine
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAuthMachineToMachine
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAuthMachineToMachine
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAuthMachineToMachine
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAuthMachineToMachine
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAuthMachineToMachine        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuthMachineToMachine          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAuthMachineToMachine = fmt.Errorf("proto: unexpected end of group")
)
