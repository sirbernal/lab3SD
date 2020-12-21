// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.1
// source: proto/AdminService.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type RegAdmRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reg string `protobuf:"bytes,1,opt,name=reg,proto3" json:"reg,omitempty"`
}

func (x *RegAdmRequest) Reset() {
	*x = RegAdmRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegAdmRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegAdmRequest) ProtoMessage() {}

func (x *RegAdmRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegAdmRequest.ProtoReflect.Descriptor instead.
func (*RegAdmRequest) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{0}
}

func (x *RegAdmRequest) GetReg() string {
	if x != nil {
		return x.Reg
	}
	return ""
}

type RegAdmResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RegAdmResponse) Reset() {
	*x = RegAdmResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegAdmResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegAdmResponse) ProtoMessage() {}

func (x *RegAdmResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegAdmResponse.ProtoReflect.Descriptor instead.
func (*RegAdmResponse) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{1}
}

func (x *RegAdmResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type BrokerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action string `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
}

func (x *BrokerRequest) Reset() {
	*x = BrokerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BrokerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BrokerRequest) ProtoMessage() {}

func (x *BrokerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BrokerRequest.ProtoReflect.Descriptor instead.
func (*BrokerRequest) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{2}
}

func (x *BrokerRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type BrokerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip    string  `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Clock []int64 `protobuf:"varint,2,rep,packed,name=clock,proto3" json:"clock,omitempty"`
}

func (x *BrokerResponse) Reset() {
	*x = BrokerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BrokerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BrokerResponse) ProtoMessage() {}

func (x *BrokerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BrokerResponse.ProtoReflect.Descriptor instead.
func (*BrokerResponse) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{3}
}

func (x *BrokerResponse) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *BrokerResponse) GetClock() []int64 {
	if x != nil {
		return x.Clock
	}
	return nil
}

type DnsCommandRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command []string `protobuf:"bytes,2,rep,name=command,proto3" json:"command,omitempty"`
}

func (x *DnsCommandRequest) Reset() {
	*x = DnsCommandRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DnsCommandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DnsCommandRequest) ProtoMessage() {}

func (x *DnsCommandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DnsCommandRequest.ProtoReflect.Descriptor instead.
func (*DnsCommandRequest) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{4}
}

func (x *DnsCommandRequest) GetCommand() []string {
	if x != nil {
		return x.Command
	}
	return nil
}

type DnsCommandResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clock []int64 `protobuf:"varint,2,rep,packed,name=clock,proto3" json:"clock,omitempty"`
}

func (x *DnsCommandResponse) Reset() {
	*x = DnsCommandResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_AdminService_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DnsCommandResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DnsCommandResponse) ProtoMessage() {}

func (x *DnsCommandResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_AdminService_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DnsCommandResponse.ProtoReflect.Descriptor instead.
func (*DnsCommandResponse) Descriptor() ([]byte, []int) {
	return file_proto_AdminService_proto_rawDescGZIP(), []int{5}
}

func (x *DnsCommandResponse) GetClock() []int64 {
	if x != nil {
		return x.Clock
	}
	return nil
}

var File_proto_AdminService_proto protoreflect.FileDescriptor

var file_proto_AdminService_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x21, 0x0a, 0x0d, 0x52, 0x65, 0x67,
	0x41, 0x64, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x67, 0x22, 0x20, 0x0a, 0x0e,
	0x52, 0x65, 0x67, 0x41, 0x64, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x27,
	0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x36, 0x0a, 0x0e, 0x42, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f,
	0x63, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x22,
	0x2d, 0x0a, 0x11, 0x44, 0x6e, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x2a,
	0x0a, 0x12, 0x44, 0x6e, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x32, 0xf5, 0x01, 0x0a, 0x0c, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x06, 0x52,
	0x65, 0x67, 0x41, 0x64, 0x6d, 0x12, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x41, 0x64, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x41, 0x64, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x06, 0x42, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x12, 0x1c,
	0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42,
	0x72, 0x6f, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x72, 0x6f,
	0x6b, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x53, 0x0a,
	0x0a, 0x44, 0x6e, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x20, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x6e, 0x73, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x6e,
	0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_AdminService_proto_rawDescOnce sync.Once
	file_proto_AdminService_proto_rawDescData = file_proto_AdminService_proto_rawDesc
)

func file_proto_AdminService_proto_rawDescGZIP() []byte {
	file_proto_AdminService_proto_rawDescOnce.Do(func() {
		file_proto_AdminService_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_AdminService_proto_rawDescData)
	})
	return file_proto_AdminService_proto_rawDescData
}

var file_proto_AdminService_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_AdminService_proto_goTypes = []interface{}{
	(*RegAdmRequest)(nil),      // 0: admin_service.RegAdmRequest
	(*RegAdmResponse)(nil),     // 1: admin_service.RegAdmResponse
	(*BrokerRequest)(nil),      // 2: admin_service.BrokerRequest
	(*BrokerResponse)(nil),     // 3: admin_service.BrokerResponse
	(*DnsCommandRequest)(nil),  // 4: admin_service.DnsCommandRequest
	(*DnsCommandResponse)(nil), // 5: admin_service.DnsCommandResponse
}
var file_proto_AdminService_proto_depIdxs = []int32{
	0, // 0: admin_service.AdminService.RegAdm:input_type -> admin_service.RegAdmRequest
	2, // 1: admin_service.AdminService.Broker:input_type -> admin_service.BrokerRequest
	4, // 2: admin_service.AdminService.DnsCommand:input_type -> admin_service.DnsCommandRequest
	1, // 3: admin_service.AdminService.RegAdm:output_type -> admin_service.RegAdmResponse
	3, // 4: admin_service.AdminService.Broker:output_type -> admin_service.BrokerResponse
	5, // 5: admin_service.AdminService.DnsCommand:output_type -> admin_service.DnsCommandResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_AdminService_proto_init() }
func file_proto_AdminService_proto_init() {
	if File_proto_AdminService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_AdminService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegAdmRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_AdminService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegAdmResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_AdminService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BrokerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_AdminService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BrokerResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_AdminService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DnsCommandRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_AdminService_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DnsCommandResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_AdminService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_AdminService_proto_goTypes,
		DependencyIndexes: file_proto_AdminService_proto_depIdxs,
		MessageInfos:      file_proto_AdminService_proto_msgTypes,
	}.Build()
	File_proto_AdminService_proto = out.File
	file_proto_AdminService_proto_rawDesc = nil
	file_proto_AdminService_proto_goTypes = nil
	file_proto_AdminService_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdminServiceClient interface {
	RegAdm(ctx context.Context, in *RegAdmRequest, opts ...grpc.CallOption) (*RegAdmResponse, error)
	Broker(ctx context.Context, in *BrokerRequest, opts ...grpc.CallOption) (*BrokerResponse, error)
	DnsCommand(ctx context.Context, in *DnsCommandRequest, opts ...grpc.CallOption) (*DnsCommandResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) RegAdm(ctx context.Context, in *RegAdmRequest, opts ...grpc.CallOption) (*RegAdmResponse, error) {
	out := new(RegAdmResponse)
	err := c.cc.Invoke(ctx, "/admin_service.AdminService/RegAdm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) Broker(ctx context.Context, in *BrokerRequest, opts ...grpc.CallOption) (*BrokerResponse, error) {
	out := new(BrokerResponse)
	err := c.cc.Invoke(ctx, "/admin_service.AdminService/Broker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) DnsCommand(ctx context.Context, in *DnsCommandRequest, opts ...grpc.CallOption) (*DnsCommandResponse, error) {
	out := new(DnsCommandResponse)
	err := c.cc.Invoke(ctx, "/admin_service.AdminService/DnsCommand", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
type AdminServiceServer interface {
	RegAdm(context.Context, *RegAdmRequest) (*RegAdmResponse, error)
	Broker(context.Context, *BrokerRequest) (*BrokerResponse, error)
	DnsCommand(context.Context, *DnsCommandRequest) (*DnsCommandResponse, error)
}

// UnimplementedAdminServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (*UnimplementedAdminServiceServer) RegAdm(context.Context, *RegAdmRequest) (*RegAdmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegAdm not implemented")
}
func (*UnimplementedAdminServiceServer) Broker(context.Context, *BrokerRequest) (*BrokerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Broker not implemented")
}
func (*UnimplementedAdminServiceServer) DnsCommand(context.Context, *DnsCommandRequest) (*DnsCommandResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DnsCommand not implemented")
}

func RegisterAdminServiceServer(s *grpc.Server, srv AdminServiceServer) {
	s.RegisterService(&_AdminService_serviceDesc, srv)
}

func _AdminService_RegAdm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegAdmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).RegAdm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_service.AdminService/RegAdm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).RegAdm(ctx, req.(*RegAdmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_Broker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrokerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).Broker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_service.AdminService/Broker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).Broker(ctx, req.(*BrokerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_DnsCommand_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DnsCommandRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).DnsCommand(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_service.AdminService/DnsCommand",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).DnsCommand(ctx, req.(*DnsCommandRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdminService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "admin_service.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegAdm",
			Handler:    _AdminService_RegAdm_Handler,
		},
		{
			MethodName: "Broker",
			Handler:    _AdminService_Broker_Handler,
		},
		{
			MethodName: "DnsCommand",
			Handler:    _AdminService_DnsCommand_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/AdminService.proto",
}
