// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: internalapi/central/process_listening_on_ports_update.proto

package central

import (
	fmt "fmt"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	storage "github.com/stackrox/rox/generated/storage"
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

type ProcessListeningOnPortsUpdate struct {
	ProcessesListeningOnPorts []*storage.ProcessListeningOnPortFromSensor `protobuf:"bytes,1,rep,name=processes_listening_on_ports,json=processesListeningOnPorts,proto3" json:"processes_listening_on_ports,omitempty"`
	Time                      *types.Timestamp                            `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}                                    `json:"-"`
	XXX_unrecognized          []byte                                      `json:"-"`
	XXX_sizecache             int32                                       `json:"-"`
}

func (m *ProcessListeningOnPortsUpdate) Reset()         { *m = ProcessListeningOnPortsUpdate{} }
func (m *ProcessListeningOnPortsUpdate) String() string { return proto.CompactTextString(m) }
func (*ProcessListeningOnPortsUpdate) ProtoMessage()    {}
func (*ProcessListeningOnPortsUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d79a6d2a351caaa, []int{0}
}
func (m *ProcessListeningOnPortsUpdate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessListeningOnPortsUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessListeningOnPortsUpdate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessListeningOnPortsUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessListeningOnPortsUpdate.Merge(m, src)
}
func (m *ProcessListeningOnPortsUpdate) XXX_Size() int {
	return m.Size()
}
func (m *ProcessListeningOnPortsUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessListeningOnPortsUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessListeningOnPortsUpdate proto.InternalMessageInfo

func (m *ProcessListeningOnPortsUpdate) GetProcessesListeningOnPorts() []*storage.ProcessListeningOnPortFromSensor {
	if m != nil {
		return m.ProcessesListeningOnPorts
	}
	return nil
}

func (m *ProcessListeningOnPortsUpdate) GetTime() *types.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

func (m *ProcessListeningOnPortsUpdate) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ProcessListeningOnPortsUpdate) Clone() *ProcessListeningOnPortsUpdate {
	if m == nil {
		return nil
	}
	cloned := new(ProcessListeningOnPortsUpdate)
	*cloned = *m

	if m.ProcessesListeningOnPorts != nil {
		cloned.ProcessesListeningOnPorts = make([]*storage.ProcessListeningOnPortFromSensor, len(m.ProcessesListeningOnPorts))
		for idx, v := range m.ProcessesListeningOnPorts {
			cloned.ProcessesListeningOnPorts[idx] = v.Clone()
		}
	}
	cloned.Time = m.Time.Clone()
	return cloned
}

func init() {
	proto.RegisterType((*ProcessListeningOnPortsUpdate)(nil), "central.ProcessListeningOnPortsUpdate")
}

func init() {
	proto.RegisterFile("internalapi/central/process_listening_on_ports_update.proto", fileDescriptor_4d79a6d2a351caaa)
}

var fileDescriptor_4d79a6d2a351caaa = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0xce, 0xcc, 0x2b, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0x49, 0x2c, 0xc8, 0xd4, 0x4f, 0x4e, 0xcd, 0x2b, 0x29, 0x4a, 0xcc, 0xd1,
	0x2f, 0x28, 0xca, 0x4f, 0x4e, 0x2d, 0x2e, 0x8e, 0xcf, 0xc9, 0x2c, 0x2e, 0x49, 0xcd, 0xcb, 0xcc,
	0x4b, 0x8f, 0xcf, 0xcf, 0x8b, 0x2f, 0xc8, 0x2f, 0x2a, 0x29, 0x8e, 0x2f, 0x2d, 0x48, 0x49, 0x2c,
	0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x6a, 0x90, 0x92, 0x4f, 0xcf, 0xcf,
	0x4f, 0xcf, 0x49, 0xd5, 0x07, 0x0b, 0x27, 0x95, 0xa6, 0xe9, 0x97, 0x64, 0xe6, 0xa6, 0x16, 0x97,
	0x24, 0xe6, 0x16, 0x40, 0x54, 0x4a, 0xa9, 0x17, 0x97, 0xe4, 0x17, 0x25, 0xa6, 0xa7, 0xe2, 0x36,
	0x1a, 0xa2, 0x50, 0x69, 0x37, 0x23, 0x97, 0x6c, 0x00, 0x44, 0x8d, 0x0f, 0x4c, 0x89, 0x7f, 0x5e,
	0x00, 0xc8, 0xee, 0x50, 0xb0, 0xd5, 0x42, 0x59, 0x5c, 0x32, 0x50, 0x43, 0x52, 0xb1, 0xb9, 0x50,
	0x82, 0x51, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x53, 0x0f, 0x6a, 0xa3, 0x1e, 0x76, 0xd3, 0xdc, 0x8a,
	0xf2, 0x73, 0x83, 0x53, 0xf3, 0x8a, 0xf3, 0x8b, 0x82, 0x24, 0xe1, 0xc6, 0xa1, 0xdb, 0x28, 0xa4,
	0xc7, 0xc5, 0x02, 0xf2, 0x89, 0x04, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x94, 0x1e, 0xc4, 0x9b,
	0x7a, 0x30, 0x6f, 0xea, 0x85, 0xc0, 0xbc, 0x19, 0x04, 0x56, 0xe7, 0x24, 0x79, 0xe2, 0x91, 0x1c,
	0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7, 0x10, 0x05, 0x0b,
	0xa2, 0x24, 0x36, 0xb0, 0x26, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x82, 0xb8, 0x72,
	0x71, 0x01, 0x00, 0x00,
}

func (m *ProcessListeningOnPortsUpdate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessListeningOnPortsUpdate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessListeningOnPortsUpdate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Time != nil {
		{
			size, err := m.Time.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintProcessListeningOnPortsUpdate(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.ProcessesListeningOnPorts) > 0 {
		for iNdEx := len(m.ProcessesListeningOnPorts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ProcessesListeningOnPorts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProcessListeningOnPortsUpdate(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintProcessListeningOnPortsUpdate(dAtA []byte, offset int, v uint64) int {
	offset -= sovProcessListeningOnPortsUpdate(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProcessListeningOnPortsUpdate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ProcessesListeningOnPorts) > 0 {
		for _, e := range m.ProcessesListeningOnPorts {
			l = e.Size()
			n += 1 + l + sovProcessListeningOnPortsUpdate(uint64(l))
		}
	}
	if m.Time != nil {
		l = m.Time.Size()
		n += 1 + l + sovProcessListeningOnPortsUpdate(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovProcessListeningOnPortsUpdate(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProcessListeningOnPortsUpdate(x uint64) (n int) {
	return sovProcessListeningOnPortsUpdate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProcessListeningOnPortsUpdate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessListeningOnPortsUpdate
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
			return fmt.Errorf("proto: ProcessListeningOnPortsUpdate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessListeningOnPortsUpdate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProcessesListeningOnPorts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPortsUpdate
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
				return ErrInvalidLengthProcessListeningOnPortsUpdate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPortsUpdate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProcessesListeningOnPorts = append(m.ProcessesListeningOnPorts, &storage.ProcessListeningOnPortFromSensor{})
			if err := m.ProcessesListeningOnPorts[len(m.ProcessesListeningOnPorts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessListeningOnPortsUpdate
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
				return ErrInvalidLengthProcessListeningOnPortsUpdate
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProcessListeningOnPortsUpdate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Time == nil {
				m.Time = &types.Timestamp{}
			}
			if err := m.Time.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessListeningOnPortsUpdate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProcessListeningOnPortsUpdate
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
func skipProcessListeningOnPortsUpdate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProcessListeningOnPortsUpdate
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
					return 0, ErrIntOverflowProcessListeningOnPortsUpdate
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
					return 0, ErrIntOverflowProcessListeningOnPortsUpdate
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
				return 0, ErrInvalidLengthProcessListeningOnPortsUpdate
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProcessListeningOnPortsUpdate
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProcessListeningOnPortsUpdate
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProcessListeningOnPortsUpdate        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProcessListeningOnPortsUpdate          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProcessListeningOnPortsUpdate = fmt.Errorf("proto: unexpected end of group")
)