// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package server

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
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

type RefreshCtagsRequest struct {
	Directory            string   `protobuf:"bytes,1,opt,name=directory,proto3" json:"directory,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RefreshCtagsRequest) Reset()         { *m = RefreshCtagsRequest{} }
func (m *RefreshCtagsRequest) String() string { return proto.CompactTextString(m) }
func (*RefreshCtagsRequest) ProtoMessage()    {}
func (*RefreshCtagsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{0}
}

func (m *RefreshCtagsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RefreshCtagsRequest.Unmarshal(m, b)
}
func (m *RefreshCtagsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RefreshCtagsRequest.Marshal(b, m, deterministic)
}
func (m *RefreshCtagsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RefreshCtagsRequest.Merge(m, src)
}
func (m *RefreshCtagsRequest) XXX_Size() int {
	return xxx_messageInfo_RefreshCtagsRequest.Size(m)
}
func (m *RefreshCtagsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RefreshCtagsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RefreshCtagsRequest proto.InternalMessageInfo

func (m *RefreshCtagsRequest) GetDirectory() string {
	if m != nil {
		return m.Directory
	}
	return ""
}

type RefreshCtagsReply struct {
	File                 string   `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RefreshCtagsReply) Reset()         { *m = RefreshCtagsReply{} }
func (m *RefreshCtagsReply) String() string { return proto.CompactTextString(m) }
func (*RefreshCtagsReply) ProtoMessage()    {}
func (*RefreshCtagsReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{1}
}

func (m *RefreshCtagsReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RefreshCtagsReply.Unmarshal(m, b)
}
func (m *RefreshCtagsReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RefreshCtagsReply.Marshal(b, m, deterministic)
}
func (m *RefreshCtagsReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RefreshCtagsReply.Merge(m, src)
}
func (m *RefreshCtagsReply) XXX_Size() int {
	return xxx_messageInfo_RefreshCtagsReply.Size(m)
}
func (m *RefreshCtagsReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RefreshCtagsReply.DiscardUnknown(m)
}

var xxx_messageInfo_RefreshCtagsReply proto.InternalMessageInfo

func (m *RefreshCtagsReply) GetFile() string {
	if m != nil {
		return m.File
	}
	return ""
}

func init() {
	proto.RegisterType((*RefreshCtagsRequest)(nil), "server.RefreshCtagsRequest")
	proto.RegisterType((*RefreshCtagsReply)(nil), "server.RefreshCtagsReply")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x8c, 0xb9, 0x84,
	0x83, 0x52, 0xd3, 0x8a, 0x52, 0x8b, 0x33, 0x9c, 0x4b, 0x12, 0xd3, 0x8b, 0x83, 0x52, 0x0b, 0x4b,
	0x53, 0x8b, 0x4b, 0x84, 0x64, 0xb8, 0x38, 0x53, 0x32, 0x8b, 0x52, 0x93, 0x4b, 0xf2, 0x8b, 0x2a,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x10, 0x02, 0x4a, 0xea, 0x5c, 0x82, 0xa8, 0x9a, 0x0a,
	0x72, 0x2a, 0x85, 0x84, 0xb8, 0x58, 0xd2, 0x32, 0x73, 0x52, 0xa1, 0xaa, 0xc1, 0x6c, 0xa3, 0x20,
	0x2e, 0xb6, 0x60, 0xb0, 0x3d, 0x42, 0x1e, 0x5c, 0x3c, 0xc8, 0x5a, 0x84, 0xa4, 0xf5, 0xa0, 0xce,
	0xc1, 0x62, 0xbb, 0x94, 0x24, 0x76, 0xc9, 0x82, 0x9c, 0x4a, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x07,
	0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x79, 0xcc, 0xc3, 0xd0, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServerClient is the client API for Server service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServerClient interface {
	RefreshCtags(ctx context.Context, in *RefreshCtagsRequest, opts ...grpc.CallOption) (*RefreshCtagsReply, error)
}

type serverClient struct {
	cc *grpc.ClientConn
}

func NewServerClient(cc *grpc.ClientConn) ServerClient {
	return &serverClient{cc}
}

func (c *serverClient) RefreshCtags(ctx context.Context, in *RefreshCtagsRequest, opts ...grpc.CallOption) (*RefreshCtagsReply, error) {
	out := new(RefreshCtagsReply)
	err := c.cc.Invoke(ctx, "/server.Server/RefreshCtags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServer is the server API for Server service.
type ServerServer interface {
	RefreshCtags(context.Context, *RefreshCtagsRequest) (*RefreshCtagsReply, error)
}

func RegisterServerServer(s *grpc.Server, srv ServerServer) {
	s.RegisterService(&_Server_serviceDesc, srv)
}

func _Server_RefreshCtags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshCtagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServer).RefreshCtags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Server/RefreshCtags",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServer).RefreshCtags(ctx, req.(*RefreshCtagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Server_serviceDesc = grpc.ServiceDesc{
	ServiceName: "server.Server",
	HandlerType: (*ServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RefreshCtags",
			Handler:    _Server_RefreshCtags_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
