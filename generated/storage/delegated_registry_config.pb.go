// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/delegated_registry_config.proto

package storage

import (
	fmt "fmt"
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

type DelegatedRegistryConfig_EnabledFor int32

const (
	DelegatedRegistryConfig_NONE     DelegatedRegistryConfig_EnabledFor = 0
	DelegatedRegistryConfig_ALL      DelegatedRegistryConfig_EnabledFor = 1
	DelegatedRegistryConfig_SPECIFIC DelegatedRegistryConfig_EnabledFor = 2
)

var DelegatedRegistryConfig_EnabledFor_name = map[int32]string{
	0: "NONE",
	1: "ALL",
	2: "SPECIFIC",
}

var DelegatedRegistryConfig_EnabledFor_value = map[string]int32{
	"NONE":     0,
	"ALL":      1,
	"SPECIFIC": 2,
}

func (x DelegatedRegistryConfig_EnabledFor) String() string {
	return proto.EnumName(DelegatedRegistryConfig_EnabledFor_name, int32(x))
}

func (DelegatedRegistryConfig_EnabledFor) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e8bf9a7bdbf38232, []int{0, 0}
}

// DelegatedRegistryConfig determines how to handle scan requests.
//
// Refer to v1.DelegatedRegistryConfig for more detailed docs.
//
// Any changes made to this message must also be reflected in central/delegatedregistryconfig/convert/convert.go.
type DelegatedRegistryConfig struct {
	EnabledFor           DelegatedRegistryConfig_EnabledFor           `protobuf:"varint,1,opt,name=enabled_for,json=enabledFor,proto3,enum=storage.DelegatedRegistryConfig_EnabledFor" json:"enabled_for,omitempty"`
	DefaultClusterId     string                                       `protobuf:"bytes,2,opt,name=default_cluster_id,json=defaultClusterId,proto3" json:"default_cluster_id,omitempty"`
	Registries           []*DelegatedRegistryConfig_DelegatedRegistry `protobuf:"bytes,3,rep,name=registries,proto3" json:"registries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                     `json:"-"`
	XXX_unrecognized     []byte                                       `json:"-"`
	XXX_sizecache        int32                                        `json:"-"`
}

func (m *DelegatedRegistryConfig) Reset()         { *m = DelegatedRegistryConfig{} }
func (m *DelegatedRegistryConfig) String() string { return proto.CompactTextString(m) }
func (*DelegatedRegistryConfig) ProtoMessage()    {}
func (*DelegatedRegistryConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_e8bf9a7bdbf38232, []int{0}
}
func (m *DelegatedRegistryConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DelegatedRegistryConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DelegatedRegistryConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DelegatedRegistryConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelegatedRegistryConfig.Merge(m, src)
}
func (m *DelegatedRegistryConfig) XXX_Size() int {
	return m.Size()
}
func (m *DelegatedRegistryConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DelegatedRegistryConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DelegatedRegistryConfig proto.InternalMessageInfo

func (m *DelegatedRegistryConfig) GetEnabledFor() DelegatedRegistryConfig_EnabledFor {
	if m != nil {
		return m.EnabledFor
	}
	return DelegatedRegistryConfig_NONE
}

func (m *DelegatedRegistryConfig) GetDefaultClusterId() string {
	if m != nil {
		return m.DefaultClusterId
	}
	return ""
}

func (m *DelegatedRegistryConfig) GetRegistries() []*DelegatedRegistryConfig_DelegatedRegistry {
	if m != nil {
		return m.Registries
	}
	return nil
}

func (m *DelegatedRegistryConfig) MessageClone() proto.Message {
	return m.Clone()
}
func (m *DelegatedRegistryConfig) Clone() *DelegatedRegistryConfig {
	if m == nil {
		return nil
	}
	cloned := new(DelegatedRegistryConfig)
	*cloned = *m

	if m.Registries != nil {
		cloned.Registries = make([]*DelegatedRegistryConfig_DelegatedRegistry, len(m.Registries))
		for idx, v := range m.Registries {
			cloned.Registries[idx] = v.Clone()
		}
	}
	return cloned
}

type DelegatedRegistryConfig_DelegatedRegistry struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	ClusterId            string   `protobuf:"bytes,2,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) Reset() {
	*m = DelegatedRegistryConfig_DelegatedRegistry{}
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) String() string {
	return proto.CompactTextString(m)
}
func (*DelegatedRegistryConfig_DelegatedRegistry) ProtoMessage() {}
func (*DelegatedRegistryConfig_DelegatedRegistry) Descriptor() ([]byte, []int) {
	return fileDescriptor_e8bf9a7bdbf38232, []int{0, 0}
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DelegatedRegistryConfig_DelegatedRegistry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DelegatedRegistryConfig_DelegatedRegistry.Merge(m, src)
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) XXX_Size() int {
	return m.Size()
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) XXX_DiscardUnknown() {
	xxx_messageInfo_DelegatedRegistryConfig_DelegatedRegistry.DiscardUnknown(m)
}

var xxx_messageInfo_DelegatedRegistryConfig_DelegatedRegistry proto.InternalMessageInfo

func (m *DelegatedRegistryConfig_DelegatedRegistry) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) GetClusterId() string {
	if m != nil {
		return m.ClusterId
	}
	return ""
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) MessageClone() proto.Message {
	return m.Clone()
}
func (m *DelegatedRegistryConfig_DelegatedRegistry) Clone() *DelegatedRegistryConfig_DelegatedRegistry {
	if m == nil {
		return nil
	}
	cloned := new(DelegatedRegistryConfig_DelegatedRegistry)
	*cloned = *m

	return cloned
}

func init() {
	proto.RegisterEnum("storage.DelegatedRegistryConfig_EnabledFor", DelegatedRegistryConfig_EnabledFor_name, DelegatedRegistryConfig_EnabledFor_value)
	proto.RegisterType((*DelegatedRegistryConfig)(nil), "storage.DelegatedRegistryConfig")
	proto.RegisterType((*DelegatedRegistryConfig_DelegatedRegistry)(nil), "storage.DelegatedRegistryConfig.DelegatedRegistry")
}

func init() {
	proto.RegisterFile("storage/delegated_registry_config.proto", fileDescriptor_e8bf9a7bdbf38232)
}

var fileDescriptor_e8bf9a7bdbf38232 = []byte{
	// 298 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x3b, 0x6d, 0xb1, 0xcd, 0xad, 0x48, 0x9c, 0x8d, 0x51, 0x30, 0x84, 0x6e, 0x0c, 0xa8,
	0x11, 0xa2, 0x2f, 0xa0, 0x31, 0x81, 0x40, 0xa8, 0x32, 0xee, 0xdc, 0x84, 0x69, 0x66, 0x12, 0x83,
	0xa1, 0x53, 0x26, 0x53, 0xd0, 0x37, 0xf1, 0x91, 0x5c, 0xfa, 0x08, 0x25, 0xbe, 0x88, 0x90, 0x1f,
	0x15, 0x82, 0xb8, 0xbb, 0xdc, 0xfb, 0x9d, 0x99, 0x73, 0x38, 0x70, 0x52, 0x2a, 0x21, 0x69, 0xc6,
	0x2f, 0x18, 0x2f, 0x78, 0x46, 0x15, 0x67, 0xb1, 0xe4, 0x59, 0x5e, 0x2a, 0xf9, 0x1a, 0x27, 0x62,
	0x95, 0xe6, 0x99, 0xb3, 0x96, 0x42, 0x09, 0x3c, 0x69, 0xc1, 0xf9, 0x76, 0x08, 0x07, 0xb7, 0x1d,
	0x4c, 0x5a, 0xd6, 0xab, 0x51, 0x1c, 0xc1, 0x8c, 0xaf, 0xe8, 0xb2, 0xe0, 0x2c, 0x4e, 0x85, 0x34,
	0x90, 0x85, 0xec, 0x3d, 0xf7, 0xd4, 0x69, 0xa5, 0xce, 0x1f, 0x32, 0xc7, 0x6f, 0x34, 0x81, 0x90,
	0x04, 0xf8, 0xf7, 0x8c, 0xcf, 0x00, 0x33, 0x9e, 0xd2, 0x4d, 0xa1, 0xe2, 0xa4, 0xd8, 0x94, 0x8a,
	0xcb, 0x38, 0x67, 0xc6, 0xd0, 0x42, 0xb6, 0x46, 0xf4, 0xf6, 0xe2, 0x35, 0x87, 0x90, 0x61, 0x02,
	0xd0, 0x3a, 0xcf, 0x79, 0x69, 0x8c, 0xac, 0x91, 0x3d, 0x73, 0xdd, 0x7f, 0xbf, 0xee, 0xed, 0xc9,
	0xaf, 0x57, 0x8e, 0x02, 0xd8, 0xef, 0x01, 0x18, 0xc3, 0x78, 0x4d, 0xd5, 0x53, 0x9d, 0x4e, 0x23,
	0xf5, 0x8c, 0x8f, 0x01, 0x7a, 0x16, 0xb5, 0xa4, 0xf3, 0x36, 0x3f, 0x07, 0xf8, 0xc9, 0x88, 0xa7,
	0x30, 0x5e, 0xdc, 0x2d, 0x7c, 0x7d, 0x80, 0x27, 0x30, 0xba, 0x8e, 0x22, 0x1d, 0xe1, 0x5d, 0x98,
	0x3e, 0xdc, 0xfb, 0x5e, 0x18, 0x84, 0x9e, 0x3e, 0xbc, 0xb9, 0x7a, 0xaf, 0x4c, 0xf4, 0x51, 0x99,
	0x68, 0x5b, 0x99, 0xe8, 0xed, 0xd3, 0x1c, 0xc0, 0x61, 0x2e, 0x9c, 0x52, 0xd1, 0xe4, 0x59, 0x8a,
	0x97, 0xa6, 0x90, 0x2e, 0xd9, 0x63, 0x57, 0xcc, 0x72, 0xa7, 0xde, 0x5f, 0x7e, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x30, 0x60, 0x72, 0xdd, 0xd3, 0x01, 0x00, 0x00,
}

func (m *DelegatedRegistryConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DelegatedRegistryConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DelegatedRegistryConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Registries) > 0 {
		for iNdEx := len(m.Registries) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Registries[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintDelegatedRegistryConfig(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.DefaultClusterId) > 0 {
		i -= len(m.DefaultClusterId)
		copy(dAtA[i:], m.DefaultClusterId)
		i = encodeVarintDelegatedRegistryConfig(dAtA, i, uint64(len(m.DefaultClusterId)))
		i--
		dAtA[i] = 0x12
	}
	if m.EnabledFor != 0 {
		i = encodeVarintDelegatedRegistryConfig(dAtA, i, uint64(m.EnabledFor))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ClusterId) > 0 {
		i -= len(m.ClusterId)
		copy(dAtA[i:], m.ClusterId)
		i = encodeVarintDelegatedRegistryConfig(dAtA, i, uint64(len(m.ClusterId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Path) > 0 {
		i -= len(m.Path)
		copy(dAtA[i:], m.Path)
		i = encodeVarintDelegatedRegistryConfig(dAtA, i, uint64(len(m.Path)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDelegatedRegistryConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovDelegatedRegistryConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DelegatedRegistryConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EnabledFor != 0 {
		n += 1 + sovDelegatedRegistryConfig(uint64(m.EnabledFor))
	}
	l = len(m.DefaultClusterId)
	if l > 0 {
		n += 1 + l + sovDelegatedRegistryConfig(uint64(l))
	}
	if len(m.Registries) > 0 {
		for _, e := range m.Registries {
			l = e.Size()
			n += 1 + l + sovDelegatedRegistryConfig(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DelegatedRegistryConfig_DelegatedRegistry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovDelegatedRegistryConfig(uint64(l))
	}
	l = len(m.ClusterId)
	if l > 0 {
		n += 1 + l + sovDelegatedRegistryConfig(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDelegatedRegistryConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDelegatedRegistryConfig(x uint64) (n int) {
	return sovDelegatedRegistryConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DelegatedRegistryConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDelegatedRegistryConfig
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
			return fmt.Errorf("proto: DelegatedRegistryConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DelegatedRegistryConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnabledFor", wireType)
			}
			m.EnabledFor = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDelegatedRegistryConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EnabledFor |= DelegatedRegistryConfig_EnabledFor(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultClusterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDelegatedRegistryConfig
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
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DefaultClusterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Registries", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDelegatedRegistryConfig
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
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Registries = append(m.Registries, &DelegatedRegistryConfig_DelegatedRegistry{})
			if err := m.Registries[len(m.Registries)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDelegatedRegistryConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
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
func (m *DelegatedRegistryConfig_DelegatedRegistry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDelegatedRegistryConfig
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
			return fmt.Errorf("proto: DelegatedRegistry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DelegatedRegistry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDelegatedRegistryConfig
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
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDelegatedRegistryConfig
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
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDelegatedRegistryConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDelegatedRegistryConfig
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
func skipDelegatedRegistryConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDelegatedRegistryConfig
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
					return 0, ErrIntOverflowDelegatedRegistryConfig
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
					return 0, ErrIntOverflowDelegatedRegistryConfig
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
				return 0, ErrInvalidLengthDelegatedRegistryConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDelegatedRegistryConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDelegatedRegistryConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDelegatedRegistryConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDelegatedRegistryConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDelegatedRegistryConfig = fmt.Errorf("proto: unexpected end of group")
)