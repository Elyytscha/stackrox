// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v1/telemetry_service.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	central "github.com/stackrox/rox/generated/internalapi/central"
	storage "github.com/stackrox/rox/generated/storage"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ConfigureTelemetryRequest struct {
	Enabled              bool     `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConfigureTelemetryRequest) Reset()         { *m = ConfigureTelemetryRequest{} }
func (m *ConfigureTelemetryRequest) String() string { return proto.CompactTextString(m) }
func (*ConfigureTelemetryRequest) ProtoMessage()    {}
func (*ConfigureTelemetryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d29ceed52498e29, []int{0}
}
func (m *ConfigureTelemetryRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConfigureTelemetryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConfigureTelemetryRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConfigureTelemetryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigureTelemetryRequest.Merge(m, src)
}
func (m *ConfigureTelemetryRequest) XXX_Size() int {
	return m.Size()
}
func (m *ConfigureTelemetryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigureTelemetryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigureTelemetryRequest proto.InternalMessageInfo

func (m *ConfigureTelemetryRequest) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *ConfigureTelemetryRequest) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ConfigureTelemetryRequest) Clone() *ConfigureTelemetryRequest {
	if m == nil {
		return nil
	}
	cloned := new(ConfigureTelemetryRequest)
	*cloned = *m

	return cloned
}

func init() {
	proto.RegisterType((*ConfigureTelemetryRequest)(nil), "v1.ConfigureTelemetryRequest")
}

func init() { proto.RegisterFile("api/v1/telemetry_service.proto", fileDescriptor_3d29ceed52498e29) }

var fileDescriptor_3d29ceed52498e29 = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xbb, 0x4e, 0xeb, 0x40,
	0x10, 0x86, 0xe3, 0x1c, 0x29, 0x39, 0xd9, 0xea, 0x68, 0x75, 0x20, 0x8e, 0x15, 0x9c, 0xc8, 0x34,
	0x28, 0xc5, 0x5a, 0x06, 0xd1, 0xa4, 0x41, 0x5c, 0xa2, 0x74, 0x14, 0x81, 0x02, 0xd1, 0xa0, 0x8d,
	0x33, 0x98, 0x15, 0xce, 0xae, 0x59, 0x6f, 0xac, 0xd0, 0xf2, 0x0a, 0x34, 0x3c, 0x12, 0xa2, 0x42,
	0xe2, 0x05, 0x50, 0xe0, 0x41, 0x90, 0x6f, 0xb9, 0x91, 0x40, 0x67, 0xcf, 0xcc, 0xff, 0xcd, 0xff,
	0x8f, 0x8d, 0x4c, 0x1a, 0x30, 0x3b, 0x72, 0x6c, 0x05, 0x3e, 0x0c, 0x41, 0xc9, 0xfb, 0xab, 0x10,
	0x64, 0xc4, 0x5c, 0x20, 0x81, 0x14, 0x4a, 0xe0, 0x62, 0xe4, 0x18, 0x75, 0x4f, 0x08, 0xcf, 0x07,
	0x3b, 0x1e, 0xa5, 0x9c, 0x0b, 0x45, 0x15, 0x13, 0x3c, 0x4c, 0x27, 0x8c, 0x6a, 0xa8, 0x84, 0xa4,
	0x1e, 0xcc, 0x10, 0x59, 0x03, 0x67, 0x68, 0x18, 0x06, 0x2a, 0xaf, 0x6d, 0x33, 0xae, 0x40, 0x72,
	0xea, 0xc7, 0x3d, 0x17, 0xb8, 0x92, 0xd4, 0x5f, 0x16, 0x5a, 0xfb, 0xa8, 0x76, 0x2c, 0xf8, 0x35,
	0xf3, 0x46, 0x12, 0xce, 0xf3, 0x5e, 0x0f, 0xee, 0x46, 0x10, 0x2a, 0xac, 0xa3, 0x32, 0x70, 0xda,
	0xf7, 0x61, 0xa0, 0x6b, 0x4d, 0x6d, 0xe7, 0x6f, 0x2f, 0x7f, 0xdd, 0x7d, 0xf9, 0x83, 0xfe, 0x4d,
	0xc7, 0xcf, 0xd2, 0x14, 0xf8, 0x06, 0xd5, 0xba, 0xa0, 0xa6, 0xe5, 0x9c, 0x9b, 0x24, 0xc0, 0x15,
	0x12, 0x39, 0xa4, 0x13, 0xdb, 0x33, 0x1a, 0x24, 0x8b, 0x41, 0x56, 0xcf, 0x5a, 0x8d, 0x87, 0xb7,
	0xcf, 0xc7, 0x62, 0x0d, 0x57, 0x17, 0xae, 0x65, 0xbb, 0xb9, 0x51, 0x3c, 0x46, 0xf8, 0xbb, 0x6b,
	0xbc, 0x15, 0xaf, 0x58, 0x9b, 0xe6, 0xf7, 0xb5, 0x56, 0xb2, 0xb6, 0x6e, 0xac, 0x5b, 0xdb, 0xd6,
	0x5a, 0xf8, 0x14, 0x55, 0xba, 0xa0, 0x52, 0xdd, 0x7c, 0x26, 0x9d, 0x64, 0x17, 0x5e, 0x86, 0x5b,
	0xf5, 0x84, 0xba, 0x89, 0xff, 0xaf, 0xa2, 0xe2, 0x43, 0x54, 0x3e, 0x61, 0x61, 0x7c, 0xd4, 0x79,
	0xda, 0xec, 0xd1, 0x6a, 0x26, 0x72, 0xc3, 0xd8, 0x58, 0x94, 0x0f, 0x52, 0x51, 0x6c, 0xe9, 0x00,
	0x95, 0x3a, 0xfc, 0x07, 0x42, 0x76, 0x4d, 0x63, 0xc9, 0x40, 0xfa, 0x29, 0xdb, 0x5a, 0xeb, 0x88,
	0x3c, 0x4f, 0x4c, 0xed, 0x75, 0x62, 0x6a, 0xef, 0x13, 0x53, 0x7b, 0xfa, 0x30, 0x0b, 0x48, 0x67,
	0x82, 0x84, 0x8a, 0xba, 0xb7, 0x52, 0x8c, 0xd3, 0x1f, 0x85, 0xd0, 0x80, 0x91, 0xc8, 0xb9, 0x2c,
	0x46, 0xce, 0x45, 0xa1, 0x5f, 0x4a, 0x6a, 0x7b, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9e, 0x82,
	0xd1, 0x51, 0xd2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TelemetryServiceClient is the client API for TelemetryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConnInterface.NewStream.
type TelemetryServiceClient interface {
	GetTelemetryConfiguration(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*storage.TelemetryConfiguration, error)
	ConfigureTelemetry(ctx context.Context, in *ConfigureTelemetryRequest, opts ...grpc.CallOption) (*storage.TelemetryConfiguration, error)
	GetConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*central.TelemetryConfig, error)
	Disable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Enable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type telemetryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTelemetryServiceClient(cc grpc.ClientConnInterface) TelemetryServiceClient {
	return &telemetryServiceClient{cc}
}

func (c *telemetryServiceClient) GetTelemetryConfiguration(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*storage.TelemetryConfiguration, error) {
	out := new(storage.TelemetryConfiguration)
	err := c.cc.Invoke(ctx, "/v1.TelemetryService/GetTelemetryConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryServiceClient) ConfigureTelemetry(ctx context.Context, in *ConfigureTelemetryRequest, opts ...grpc.CallOption) (*storage.TelemetryConfiguration, error) {
	out := new(storage.TelemetryConfiguration)
	err := c.cc.Invoke(ctx, "/v1.TelemetryService/ConfigureTelemetry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryServiceClient) GetConfig(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*central.TelemetryConfig, error) {
	out := new(central.TelemetryConfig)
	err := c.cc.Invoke(ctx, "/v1.TelemetryService/GetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryServiceClient) Disable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.TelemetryService/Disable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryServiceClient) Enable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/v1.TelemetryService/Enable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TelemetryServiceServer is the server API for TelemetryService service.
type TelemetryServiceServer interface {
	GetTelemetryConfiguration(context.Context, *Empty) (*storage.TelemetryConfiguration, error)
	ConfigureTelemetry(context.Context, *ConfigureTelemetryRequest) (*storage.TelemetryConfiguration, error)
	GetConfig(context.Context, *Empty) (*central.TelemetryConfig, error)
	Disable(context.Context, *Empty) (*Empty, error)
	Enable(context.Context, *Empty) (*Empty, error)
}

// UnimplementedTelemetryServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTelemetryServiceServer struct {
}

func (*UnimplementedTelemetryServiceServer) GetTelemetryConfiguration(ctx context.Context, req *Empty) (*storage.TelemetryConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTelemetryConfiguration not implemented")
}
func (*UnimplementedTelemetryServiceServer) ConfigureTelemetry(ctx context.Context, req *ConfigureTelemetryRequest) (*storage.TelemetryConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureTelemetry not implemented")
}
func (*UnimplementedTelemetryServiceServer) GetConfig(ctx context.Context, req *Empty) (*central.TelemetryConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (*UnimplementedTelemetryServiceServer) Disable(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disable not implemented")
}
func (*UnimplementedTelemetryServiceServer) Enable(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enable not implemented")
}

func RegisterTelemetryServiceServer(s *grpc.Server, srv TelemetryServiceServer) {
	s.RegisterService(&_TelemetryService_serviceDesc, srv)
}

func _TelemetryService_GetTelemetryConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServiceServer).GetTelemetryConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TelemetryService/GetTelemetryConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServiceServer).GetTelemetryConfiguration(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TelemetryService_ConfigureTelemetry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureTelemetryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServiceServer).ConfigureTelemetry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TelemetryService/ConfigureTelemetry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServiceServer).ConfigureTelemetry(ctx, req.(*ConfigureTelemetryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TelemetryService_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServiceServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TelemetryService/GetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServiceServer).GetConfig(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TelemetryService_Disable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServiceServer).Disable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TelemetryService/Disable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServiceServer).Disable(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TelemetryService_Enable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServiceServer).Enable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.TelemetryService/Enable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServiceServer).Enable(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TelemetryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.TelemetryService",
	HandlerType: (*TelemetryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTelemetryConfiguration",
			Handler:    _TelemetryService_GetTelemetryConfiguration_Handler,
		},
		{
			MethodName: "ConfigureTelemetry",
			Handler:    _TelemetryService_ConfigureTelemetry_Handler,
		},
		{
			MethodName: "GetConfig",
			Handler:    _TelemetryService_GetConfig_Handler,
		},
		{
			MethodName: "Disable",
			Handler:    _TelemetryService_Disable_Handler,
		},
		{
			MethodName: "Enable",
			Handler:    _TelemetryService_Enable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/telemetry_service.proto",
}

func (m *ConfigureTelemetryRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConfigureTelemetryRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConfigureTelemetryRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTelemetryService(dAtA []byte, offset int, v uint64) int {
	offset -= sovTelemetryService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ConfigureTelemetryRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Enabled {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTelemetryService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTelemetryService(x uint64) (n int) {
	return sovTelemetryService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ConfigureTelemetryRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTelemetryService
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
			return fmt.Errorf("proto: ConfigureTelemetryRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConfigureTelemetryRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTelemetryService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Enabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipTelemetryService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTelemetryService
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
func skipTelemetryService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTelemetryService
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
					return 0, ErrIntOverflowTelemetryService
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
					return 0, ErrIntOverflowTelemetryService
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
				return 0, ErrInvalidLengthTelemetryService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTelemetryService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTelemetryService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTelemetryService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTelemetryService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTelemetryService = fmt.Errorf("proto: unexpected end of group")
)
