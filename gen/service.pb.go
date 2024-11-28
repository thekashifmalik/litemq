// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: service.proto

package gen

import (
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

type EnqueueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
	Data  []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *EnqueueRequest) Reset() {
	*x = EnqueueRequest{}
	mi := &file_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnqueueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnqueueRequest) ProtoMessage() {}

func (x *EnqueueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnqueueRequest.ProtoReflect.Descriptor instead.
func (*EnqueueRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *EnqueueRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

func (x *EnqueueRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Nothing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Nothing) Reset() {
	*x = Nothing{}
	mi := &file_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Nothing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nothing) ProtoMessage() {}

func (x *Nothing) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nothing.ProtoReflect.Descriptor instead.
func (*Nothing) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

type QueueID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue string `protobuf:"bytes,1,opt,name=queue,proto3" json:"queue,omitempty"`
}

func (x *QueueID) Reset() {
	*x = QueueID{}
	mi := &file_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueID) ProtoMessage() {}

func (x *QueueID) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueID.ProtoReflect.Descriptor instead.
func (*QueueID) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *QueueID) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

type QueueLength struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count int64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *QueueLength) Reset() {
	*x = QueueLength{}
	mi := &file_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueLength) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueLength) ProtoMessage() {}

func (x *QueueLength) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueLength.ProtoReflect.Descriptor instead.
func (*QueueLength) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *QueueLength) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DequeueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *DequeueResponse) Reset() {
	*x = DequeueResponse{}
	mi := &file_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DequeueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DequeueResponse) ProtoMessage() {}

func (x *DequeueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DequeueResponse.ProtoReflect.Descriptor instead.
func (*DequeueResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{4}
}

func (x *DequeueResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x3a, 0x0a, 0x0e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x09, 0x0a, 0x07, 0x4e,
	0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x1f, 0x0a, 0x07, 0x51, 0x75, 0x65, 0x75, 0x65, 0x49,
	0x44, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x0b, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x25, 0x0a, 0x0f,
	0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x32, 0xca, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12,
	0x0f, 0x2e, 0x45, 0x6e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0c, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x00,
	0x12, 0x27, 0x0a, 0x07, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x08, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x49, 0x44, 0x1a, 0x10, 0x2e, 0x44, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x21, 0x0a, 0x05, 0x50, 0x75, 0x72,
	0x67, 0x65, 0x12, 0x08, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44, 0x1a, 0x0c, 0x2e, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x00, 0x12, 0x22, 0x0a, 0x06,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x08, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x49, 0x44,
	0x1a, 0x0c, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x00,
	0x12, 0x1e, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x08, 0x2e, 0x4e, 0x6f, 0x74,
	0x68, 0x69, 0x6e, 0x67, 0x1a, 0x08, 0x2e, 0x4e, 0x6f, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x22, 0x00,
	0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_service_proto_goTypes = []any{
	(*EnqueueRequest)(nil),  // 0: EnqueueRequest
	(*Nothing)(nil),         // 1: Nothing
	(*QueueID)(nil),         // 2: QueueID
	(*QueueLength)(nil),     // 3: QueueLength
	(*DequeueResponse)(nil), // 4: DequeueResponse
}
var file_service_proto_depIdxs = []int32{
	0, // 0: QueueService.Enqueue:input_type -> EnqueueRequest
	2, // 1: QueueService.Dequeue:input_type -> QueueID
	2, // 2: QueueService.Purge:input_type -> QueueID
	2, // 3: QueueService.Length:input_type -> QueueID
	1, // 4: QueueService.Health:input_type -> Nothing
	3, // 5: QueueService.Enqueue:output_type -> QueueLength
	4, // 6: QueueService.Dequeue:output_type -> DequeueResponse
	3, // 7: QueueService.Purge:output_type -> QueueLength
	3, // 8: QueueService.Length:output_type -> QueueLength
	1, // 9: QueueService.Health:output_type -> Nothing
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
