// source: api/v2/compliance_integration_service.proto
package v2

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Next Tag: 6
type ComplianceIntegration struct {
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Version     string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	ClusterId   string `protobuf:"bytes,3,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	ClusterName string `protobuf:"bytes,4,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	Namespace   string `protobuf:"bytes,5,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Collection of errors that occurred while trying to obtain comliance operator health info.
	StatusErrors         []string `protobuf:"bytes,6,rep,name=status_errors,json=statusErrors,proto3" json:"status_errors,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComplianceIntegration) Reset()         { *m = ComplianceIntegration{} }
func (m *ComplianceIntegration) String() string { return proto.CompactTextString(m) }
func (*ComplianceIntegration) ProtoMessage()    {}
func (*ComplianceIntegration) Descriptor() ([]byte, []int) {
	return fileDescriptor_406ee33181fc48e0, []int{0}
}
func (m *ComplianceIntegration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ComplianceIntegration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ComplianceIntegration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ComplianceIntegration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComplianceIntegration.Merge(m, src)
}
func (m *ComplianceIntegration) XXX_Size() int {
	return m.Size()
}
func (m *ComplianceIntegration) XXX_DiscardUnknown() {
	xxx_messageInfo_ComplianceIntegration.DiscardUnknown(m)
}

var xxx_messageInfo_ComplianceIntegration proto.InternalMessageInfo

func (m *ComplianceIntegration) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ComplianceIntegration) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ComplianceIntegration) GetClusterId() string {
	if m != nil {
		return m.ClusterId
	}
	return ""
}

func (m *ComplianceIntegration) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func (m *ComplianceIntegration) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *ComplianceIntegration) GetStatusErrors() []string {
	if m != nil {
		return m.StatusErrors
	}
	return nil
}

func (m *ComplianceIntegration) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ComplianceIntegration) Clone() *ComplianceIntegration {
	if m == nil {
		return nil
	}
	cloned := new(ComplianceIntegration)
	*cloned = *m

	if m.StatusErrors != nil {
		cloned.StatusErrors = make([]string, len(m.StatusErrors))
		copy(cloned.StatusErrors, m.StatusErrors)
	}
	return cloned
}

type ComplianceIntegrationStatusRequest struct {
	ClusterId            string   `protobuf:"bytes,1,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComplianceIntegrationStatusRequest) Reset()         { *m = ComplianceIntegrationStatusRequest{} }
func (m *ComplianceIntegrationStatusRequest) String() string { return proto.CompactTextString(m) }
func (*ComplianceIntegrationStatusRequest) ProtoMessage()    {}
func (*ComplianceIntegrationStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_406ee33181fc48e0, []int{1}
}
func (m *ComplianceIntegrationStatusRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ComplianceIntegrationStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ComplianceIntegrationStatusRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ComplianceIntegrationStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComplianceIntegrationStatusRequest.Merge(m, src)
}
func (m *ComplianceIntegrationStatusRequest) XXX_Size() int {
	return m.Size()
}
func (m *ComplianceIntegrationStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ComplianceIntegrationStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ComplianceIntegrationStatusRequest proto.InternalMessageInfo

func (m *ComplianceIntegrationStatusRequest) GetClusterId() string {
	if m != nil {
		return m.ClusterId
	}
	return ""
}

func (m *ComplianceIntegrationStatusRequest) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ComplianceIntegrationStatusRequest) Clone() *ComplianceIntegrationStatusRequest {
	if m == nil {
		return nil
	}
	cloned := new(ComplianceIntegrationStatusRequest)
	*cloned = *m

	return cloned
}

type ListComplianceIntegrationsResponse struct {
	Integrations         []*ComplianceIntegration `protobuf:"bytes,1,rep,name=integrations,proto3" json:"integrations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ListComplianceIntegrationsResponse) Reset()         { *m = ListComplianceIntegrationsResponse{} }
func (m *ListComplianceIntegrationsResponse) String() string { return proto.CompactTextString(m) }
func (*ListComplianceIntegrationsResponse) ProtoMessage()    {}
func (*ListComplianceIntegrationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_406ee33181fc48e0, []int{2}
}
func (m *ListComplianceIntegrationsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ListComplianceIntegrationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ListComplianceIntegrationsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ListComplianceIntegrationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListComplianceIntegrationsResponse.Merge(m, src)
}
func (m *ListComplianceIntegrationsResponse) XXX_Size() int {
	return m.Size()
}
func (m *ListComplianceIntegrationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListComplianceIntegrationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListComplianceIntegrationsResponse proto.InternalMessageInfo

func (m *ListComplianceIntegrationsResponse) GetIntegrations() []*ComplianceIntegration {
	if m != nil {
		return m.Integrations
	}
	return nil
}

func (m *ListComplianceIntegrationsResponse) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ListComplianceIntegrationsResponse) Clone() *ListComplianceIntegrationsResponse {
	if m == nil {
		return nil
	}
	cloned := new(ListComplianceIntegrationsResponse)
	*cloned = *m

	if m.Integrations != nil {
		cloned.Integrations = make([]*ComplianceIntegration, len(m.Integrations))
		for idx, v := range m.Integrations {
			cloned.Integrations[idx] = v.Clone()
		}
	}
	return cloned
}

func init() {
	proto.RegisterType((*ComplianceIntegration)(nil), "v2.ComplianceIntegration")
	proto.RegisterType((*ComplianceIntegrationStatusRequest)(nil), "v2.ComplianceIntegrationStatusRequest")
	proto.RegisterType((*ListComplianceIntegrationsResponse)(nil), "v2.ListComplianceIntegrationsResponse")
}

func init() {
	proto.RegisterFile("api/v2/compliance_integration_service.proto", fileDescriptor_406ee33181fc48e0)
}

var fileDescriptor_406ee33181fc48e0 = []byte{
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> bc25ec9ebb (address comments)
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x8a, 0xd4, 0x40,
	0x10, 0xc6, 0xb7, 0x33, 0xba, 0x92, 0xde, 0xe8, 0xa1, 0x45, 0xe8, 0x8d, 0xb3, 0x61, 0xcc, 0x82,
	0x0c, 0x08, 0x09, 0xc4, 0xb3, 0x17, 0x17, 0x0f, 0x0b, 0x22, 0x18, 0x2f, 0xe2, 0x25, 0xb4, 0x9d,
	0x62, 0x6c, 0x9c, 0x74, 0x67, 0xbb, 0x3a, 0x59, 0xbd, 0xfa, 0x0a, 0x1e, 0xf4, 0x89, 0xc4, 0xa3,
	0xe0, 0x0b, 0xc8, 0xe8, 0x83, 0x48, 0xfe, 0xec, 0x9f, 0x19, 0x22, 0x1e, 0xfb, 0xf7, 0x7d, 0x55,
	0xf5, 0x35, 0x55, 0xf4, 0x91, 0xa8, 0x55, 0xda, 0x66, 0xa9, 0x34, 0x55, 0xbd, 0x56, 0x42, 0x4b,
	0x28, 0x94, 0x76, 0xb0, 0xb2, 0xc2, 0x29, 0xa3, 0x0b, 0x04, 0xdb, 0x2a, 0x09, 0x49, 0x6d, 0x8d,
	0x33, 0xcc, 0x6b, 0xb3, 0x70, 0xbe, 0x32, 0x66, 0xb5, 0x86, 0xb4, 0xab, 0x13, 0x5a, 0x1b, 0xd7,
	0x1b, 0x71, 0x70, 0x84, 0x77, 0xaf, 0xda, 0x55, 0x46, 0x8f, 0xf0, 0x70, 0x84, 0x08, 0xc2, 0xca,
	0x77, 0xc5, 0x59, 0x03, 0xf6, 0xe3, 0x20, 0xc5, 0xdf, 0x08, 0xbd, 0x77, 0x72, 0x39, 0xfa, 0xf4,
	0x6a, 0x32, 0xbb, 0x43, 0x3d, 0x55, 0x72, 0xb2, 0x20, 0x4b, 0x3f, 0xf7, 0x54, 0xc9, 0x38, 0xbd,
	0xd5, 0x82, 0x45, 0x65, 0x34, 0xf7, 0x7a, 0x78, 0xf1, 0x64, 0x47, 0x94, 0xca, 0x75, 0x83, 0x0e,
	0x6c, 0xa1, 0x4a, 0x3e, 0xeb, 0x45, 0x7f, 0x24, 0xa7, 0x25, 0x7b, 0x40, 0x83, 0x0b, 0x59, 0x8b,
	0x0a, 0xf8, 0x8d, 0xde, 0x70, 0x30, 0xb2, 0x17, 0xa2, 0x02, 0x36, 0xa7, 0x7e, 0x27, 0x61, 0x2d,
	0x24, 0xf0, 0x9b, 0x43, 0x83, 0x4b, 0xc0, 0x8e, 0xe9, 0x6d, 0x74, 0xc2, 0x35, 0x58, 0x80, 0xb5,
	0xc6, 0x22, 0xdf, 0x5f, 0xcc, 0x96, 0x7e, 0x1e, 0x0c, 0xf0, 0x59, 0xcf, 0xe2, 0x13, 0x1a, 0x4f,
	0xfe, 0xe3, 0x55, 0x6f, 0xca, 0xe1, 0xac, 0x01, 0x74, 0x3b, 0x51, 0xc9, 0x4e, 0xd4, 0x58, 0xd2,
	0xf8, 0xb9, 0x42, 0x37, 0xd9, 0x08, 0x73, 0xc0, 0xda, 0x68, 0x04, 0xf6, 0x84, 0x06, 0xd7, 0x56,
	0x84, 0x9c, 0x2c, 0x66, 0xcb, 0x83, 0xec, 0x30, 0x69, 0xb3, 0x64, 0xb2, 0x32, 0xdf, 0xb2, 0x67,
	0x5f, 0x08, 0x9d, 0x4f, 0x47, 0x1d, 0x76, 0xcd, 0xce, 0x69, 0xf8, 0xef, 0x14, 0x2c, 0xe8, 0xe6,
	0xe4, 0xe2, 0xfc, 0x65, 0xb7, 0xc5, 0xf0, 0x61, 0xf7, 0xfa, 0x7f, 0xe6, 0xf8, 0xf8, 0xd3, 0xcf,
	0x3f, 0x9f, 0xbd, 0x23, 0x76, 0x7f, 0xfb, 0xd6, 0xd2, 0xeb, 0xc9, 0x9e, 0x26, 0xdf, 0x37, 0x11,
	0xf9, 0xb1, 0x89, 0xc8, 0xaf, 0x4d, 0x44, 0xbe, 0xfe, 0x8e, 0xf6, 0x28, 0x57, 0x26, 0x41, 0x27,
	0xe4, 0x7b, 0x6b, 0x3e, 0x0c, 0x17, 0x93, 0x88, 0x5a, 0x25, 0x6d, 0xf6, 0xc6, 0x6b, 0xb3, 0xd7,
	0x7b, 0x6f, 0xf7, 0x7b, 0xf6, 0xf8, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x33, 0x22, 0x0b, 0xc5,
	0xc6, 0x02, 0x00, 0x00,
<<<<<<< HEAD
=======
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x8a, 0x13, 0x41,
	0x10, 0x86, 0xb7, 0x27, 0xba, 0x92, 0xda, 0xe8, 0xa1, 0x45, 0xe8, 0x1d, 0xb3, 0x21, 0xce, 0xc2,
	0x1a, 0x11, 0x66, 0x60, 0x3c, 0x7b, 0x71, 0x11, 0x59, 0x10, 0xc1, 0xf1, 0x22, 0x5e, 0x86, 0x76,
	0x52, 0xc4, 0xc6, 0xa4, 0x7b, 0xb6, 0xab, 0x33, 0xab, 0x88, 0x17, 0x1f, 0xc0, 0x8b, 0x17, 0x1f,
	0xc1, 0x27, 0x11, 0x8f, 0x82, 0x2f, 0x20, 0xd1, 0x07, 0x91, 0xe9, 0xc9, 0x6e, 0x36, 0x61, 0x82,
	0xc7, 0xfe, 0xea, 0xef, 0xaa, 0xbf, 0xfb, 0x2f, 0xb8, 0x2f, 0x4b, 0x95, 0x54, 0x69, 0x52, 0x98,
	0x59, 0x39, 0x55, 0x52, 0x17, 0x98, 0x2b, 0xed, 0x70, 0x62, 0xa5, 0x53, 0x46, 0xe7, 0x84, 0xb6,
	0x52, 0x05, 0xc6, 0xa5, 0x35, 0xce, 0xf0, 0xa0, 0x4a, 0xc3, 0xfe, 0xc4, 0x98, 0xc9, 0x14, 0x93,
	0xfa, 0x9e, 0xd4, 0xda, 0x38, 0x2f, 0xa4, 0x46, 0x11, 0xde, 0x5c, 0xb5, 0x9b, 0x19, 0xbd, 0x84,
	0xfb, 0x4b, 0x48, 0x28, 0x6d, 0xf1, 0x26, 0x3f, 0x9d, 0xa3, 0x7d, 0xdf, 0x94, 0xa2, 0xef, 0x0c,
	0x6e, 0x1d, 0x5f, 0x8c, 0x3e, 0x59, 0x4d, 0xe6, 0x37, 0x20, 0x50, 0x63, 0xc1, 0x86, 0x6c, 0xd4,
	0xcd, 0x02, 0x35, 0xe6, 0x02, 0xae, 0x55, 0x68, 0x49, 0x19, 0x2d, 0x02, 0x0f, 0xcf, 0x8f, 0xfc,
	0x00, 0xa0, 0x98, 0xce, 0xc9, 0xa1, 0xcd, 0xd5, 0x58, 0x74, 0x7c, 0xb1, 0xbb, 0x24, 0x27, 0x63,
	0x7e, 0x07, 0x7a, 0xe7, 0x65, 0x2d, 0x67, 0x28, 0xae, 0x78, 0xc1, 0xde, 0x92, 0x3d, 0x93, 0x33,
	0xe4, 0x7d, 0xe8, 0xd6, 0x25, 0x2a, 0x65, 0x81, 0xe2, 0x6a, 0xd3, 0xe0, 0x02, 0xf0, 0x43, 0xb8,
	0x4e, 0x4e, 0xba, 0x39, 0xe5, 0x68, 0xad, 0xb1, 0x24, 0x76, 0x87, 0x9d, 0x51, 0x37, 0xeb, 0x35,
	0xf0, 0xb1, 0x67, 0xd1, 0x31, 0x44, 0xad, 0xef, 0x78, 0xe1, 0x45, 0x19, 0x9e, 0xce, 0x91, 0xdc,
	0x86, 0x55, 0xb6, 0x61, 0x35, 0x2a, 0x20, 0x7a, 0xaa, 0xc8, 0xb5, 0x36, 0xa2, 0x0c, 0xa9, 0x34,
	0x9a, 0x90, 0x3f, 0x84, 0xde, 0xa5, 0x88, 0x48, 0xb0, 0x61, 0x67, 0xb4, 0x97, 0xee, 0xc7, 0x55,
	0x1a, 0xb7, 0xde, 0xcc, 0xd6, 0xe4, 0xe9, 0xb7, 0x00, 0xfa, 0xed, 0x56, 0x9b, 0xac, 0xf9, 0x19,
	0x84, 0xdb, 0x5d, 0xf0, 0x5e, 0x3d, 0x27, 0x93, 0x67, 0xcf, 0xeb, 0x14, 0xc3, 0xa3, 0xfa, 0xf4,
	0x7f, 0xcf, 0xd1, 0xe1, 0xa7, 0x5f, 0x7f, 0xbf, 0x04, 0x07, 0xfc, 0xf6, 0xfa, 0xae, 0x25, 0x97,
	0x9d, 0xf1, 0xcf, 0x0c, 0xc4, 0x13, 0x6c, 0x6f, 0xc5, 0x8f, 0xb6, 0xbe, 0x6f, 0xed, 0x8b, 0xc3,
	0xed, 0xff, 0x10, 0x25, 0xde, 0xc4, 0x3d, 0x7e, 0x77, 0xbb, 0x89, 0xe4, 0xc3, 0x2a, 0x9f, 0x8f,
	0x8f, 0xe2, 0x1f, 0x8b, 0x01, 0xfb, 0xb9, 0x18, 0xb0, 0xdf, 0x8b, 0x01, 0xfb, 0xfa, 0x67, 0xb0,
	0x03, 0x42, 0x99, 0x98, 0x9c, 0x2c, 0xde, 0x5a, 0xf3, 0xae, 0x59, 0xe1, 0x58, 0x96, 0x2a, 0xae,
	0xd2, 0x57, 0x41, 0x95, 0xbe, 0xdc, 0x79, 0xbd, 0xeb, 0xd9, 0x83, 0x7f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x44, 0x21, 0xfe, 0x60, 0x57, 0x03, 0x00, 0x00,
>>>>>>> 4eb94b6cbe (update generated)
=======
>>>>>>> bc25ec9ebb (address comments)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ComplianceIntegrationServiceClient is the client API for ComplianceIntegrationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConnInterface.NewStream.
type ComplianceIntegrationServiceClient interface {
	// ListComplianceIntegrations lists all the compliance operator metadata for the secured clusters
	ListComplianceIntegrations(ctx context.Context, in *RawQuery, opts ...grpc.CallOption) (*ListComplianceIntegrationsResponse, error)
}

type complianceIntegrationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewComplianceIntegrationServiceClient(cc grpc.ClientConnInterface) ComplianceIntegrationServiceClient {
	return &complianceIntegrationServiceClient{cc}
}

func (c *complianceIntegrationServiceClient) ListComplianceIntegrations(ctx context.Context, in *RawQuery, opts ...grpc.CallOption) (*ListComplianceIntegrationsResponse, error) {
	out := new(ListComplianceIntegrationsResponse)
	err := c.cc.Invoke(ctx, "/v2.ComplianceIntegrationService/ListComplianceIntegrations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ComplianceIntegrationServiceServer is the server API for ComplianceIntegrationService service.
type ComplianceIntegrationServiceServer interface {
	// ListComplianceIntegrations lists all the compliance operator metadata for the secured clusters
	ListComplianceIntegrations(context.Context, *RawQuery) (*ListComplianceIntegrationsResponse, error)
}

// UnimplementedComplianceIntegrationServiceServer can be embedded to have forward compatible implementations.
type UnimplementedComplianceIntegrationServiceServer struct {
}

func (*UnimplementedComplianceIntegrationServiceServer) ListComplianceIntegrations(ctx context.Context, req *RawQuery) (*ListComplianceIntegrationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListComplianceIntegrations not implemented")
}

func RegisterComplianceIntegrationServiceServer(s *grpc.Server, srv ComplianceIntegrationServiceServer) {
	s.RegisterService(&_ComplianceIntegrationService_serviceDesc, srv)
}

func _ComplianceIntegrationService_ListComplianceIntegrations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComplianceIntegrationServiceServer).ListComplianceIntegrations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2.ComplianceIntegrationService/ListComplianceIntegrations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComplianceIntegrationServiceServer).ListComplianceIntegrations(ctx, req.(*RawQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var _ComplianceIntegrationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2.ComplianceIntegrationService",
	HandlerType: (*ComplianceIntegrationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListComplianceIntegrations",
			Handler:    _ComplianceIntegrationService_ListComplianceIntegrations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v2/compliance_integration_service.proto",
}

func (m *ComplianceIntegration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ComplianceIntegration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComplianceIntegration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.StatusErrors) > 0 {
		for iNdEx := len(m.StatusErrors) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.StatusErrors[iNdEx])
			copy(dAtA[i:], m.StatusErrors[iNdEx])
			i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.StatusErrors[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Namespace) > 0 {
		i -= len(m.Namespace)
		copy(dAtA[i:], m.Namespace)
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.Namespace)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ClusterName) > 0 {
		i -= len(m.ClusterName)
		copy(dAtA[i:], m.ClusterName)
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.ClusterName)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ClusterId) > 0 {
		i -= len(m.ClusterId)
		copy(dAtA[i:], m.ClusterId)
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.ClusterId)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Version) > 0 {
		i -= len(m.Version)
		copy(dAtA[i:], m.Version)
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.Version)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ComplianceIntegrationStatusRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ComplianceIntegrationStatusRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComplianceIntegrationStatusRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(len(m.ClusterId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ListComplianceIntegrationsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ListComplianceIntegrationsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ListComplianceIntegrationsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Integrations) > 0 {
		for iNdEx := len(m.Integrations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Integrations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintComplianceIntegrationService(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintComplianceIntegrationService(dAtA []byte, offset int, v uint64) int {
	offset -= sovComplianceIntegrationService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ComplianceIntegration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	l = len(m.ClusterId)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	l = len(m.ClusterName)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	l = len(m.Namespace)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	if len(m.StatusErrors) > 0 {
		for _, s := range m.StatusErrors {
			l = len(s)
			n += 1 + l + sovComplianceIntegrationService(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ComplianceIntegrationStatusRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ClusterId)
	if l > 0 {
		n += 1 + l + sovComplianceIntegrationService(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ListComplianceIntegrationsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Integrations) > 0 {
		for _, e := range m.Integrations {
			l = e.Size()
			n += 1 + l + sovComplianceIntegrationService(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovComplianceIntegrationService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozComplianceIntegrationService(x uint64) (n int) {
	return sovComplianceIntegrationService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ComplianceIntegration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowComplianceIntegrationService
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
			return fmt.Errorf("proto: ComplianceIntegration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ComplianceIntegration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatusErrors", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatusErrors = append(m.StatusErrors, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipComplianceIntegrationService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthComplianceIntegrationService
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
func (m *ComplianceIntegrationStatusRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowComplianceIntegrationService
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
			return fmt.Errorf("proto: ComplianceIntegrationStatusRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ComplianceIntegrationStatusRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipComplianceIntegrationService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthComplianceIntegrationService
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
func (m *ListComplianceIntegrationsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowComplianceIntegrationService
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
			return fmt.Errorf("proto: ListComplianceIntegrationsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ListComplianceIntegrationsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Integrations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowComplianceIntegrationService
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
				return ErrInvalidLengthComplianceIntegrationService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthComplianceIntegrationService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Integrations = append(m.Integrations, &ComplianceIntegration{})
			if err := m.Integrations[len(m.Integrations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipComplianceIntegrationService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthComplianceIntegrationService
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
func skipComplianceIntegrationService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowComplianceIntegrationService
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
					return 0, ErrIntOverflowComplianceIntegrationService
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
					return 0, ErrIntOverflowComplianceIntegrationService
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
				return 0, ErrInvalidLengthComplianceIntegrationService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupComplianceIntegrationService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthComplianceIntegrationService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthComplianceIntegrationService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowComplianceIntegrationService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupComplianceIntegrationService = fmt.Errorf("proto: unexpected end of group")
)
