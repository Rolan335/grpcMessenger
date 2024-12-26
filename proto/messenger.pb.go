// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.29.2
// source: messenger.proto

package proto

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

type InitSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InitSessionRequest) Reset() {
	*x = InitSessionRequest{}
	mi := &file_messenger_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InitSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitSessionRequest) ProtoMessage() {}

func (x *InitSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitSessionRequest.ProtoReflect.Descriptor instead.
func (*InitSessionRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{0}
}

type InitSessionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionUuid   string                 `protobuf:"bytes,1,opt,name=session_uuid,json=sessionUuid,proto3" json:"session_uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InitSessionResponse) Reset() {
	*x = InitSessionResponse{}
	mi := &file_messenger_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InitSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitSessionResponse) ProtoMessage() {}

func (x *InitSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitSessionResponse.ProtoReflect.Descriptor instead.
func (*InitSessionResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{1}
}

func (x *InitSessionResponse) GetSessionUuid() string {
	if x != nil {
		return x.SessionUuid
	}
	return ""
}

type Chat struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionUuid   string                 `protobuf:"bytes,1,opt,name=session_uuid,json=sessionUuid,proto3" json:"session_uuid,omitempty"`
	Ttl           int32                  `protobuf:"varint,2,opt,name=ttl,proto3" json:"ttl,omitempty"`
	ReadOnly      bool                   `protobuf:"varint,3,opt,name=read_only,json=readOnly,proto3" json:"read_only,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Chat) Reset() {
	*x = Chat{}
	mi := &file_messenger_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Chat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chat) ProtoMessage() {}

func (x *Chat) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chat.ProtoReflect.Descriptor instead.
func (*Chat) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{2}
}

func (x *Chat) GetSessionUuid() string {
	if x != nil {
		return x.SessionUuid
	}
	return ""
}

func (x *Chat) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *Chat) GetReadOnly() bool {
	if x != nil {
		return x.ReadOnly
	}
	return false
}

type CreateChatRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Chat          *Chat                  `protobuf:"bytes,1,opt,name=chat,proto3" json:"chat,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateChatRequest) Reset() {
	*x = CreateChatRequest{}
	mi := &file_messenger_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateChatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChatRequest) ProtoMessage() {}

func (x *CreateChatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChatRequest.ProtoReflect.Descriptor instead.
func (*CreateChatRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{3}
}

func (x *CreateChatRequest) GetChat() *Chat {
	if x != nil {
		return x.Chat
	}
	return nil
}

type CreateChatResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChatUuid      string                 `protobuf:"bytes,1,opt,name=chat_uuid,json=chatUuid,proto3" json:"chat_uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateChatResponse) Reset() {
	*x = CreateChatResponse{}
	mi := &file_messenger_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateChatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateChatResponse) ProtoMessage() {}

func (x *CreateChatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateChatResponse.ProtoReflect.Descriptor instead.
func (*CreateChatResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{4}
}

func (x *CreateChatResponse) GetChatUuid() string {
	if x != nil {
		return x.ChatUuid
	}
	return ""
}

type SendMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChatUuid      string                 `protobuf:"bytes,1,opt,name=chat_uuid,json=chatUuid,proto3" json:"chat_uuid,omitempty"`
	SessionUuid   string                 `protobuf:"bytes,2,opt,name=session_uuid,json=sessionUuid,proto3" json:"session_uuid,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	mi := &file_messenger_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{5}
}

func (x *SendMessageRequest) GetChatUuid() string {
	if x != nil {
		return x.ChatUuid
	}
	return ""
}

func (x *SendMessageRequest) GetSessionUuid() string {
	if x != nil {
		return x.SessionUuid
	}
	return ""
}

func (x *SendMessageRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type SendMessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageResponse) Reset() {
	*x = SendMessageResponse{}
	mi := &file_messenger_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageResponse) ProtoMessage() {}

func (x *SendMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageResponse.ProtoReflect.Descriptor instead.
func (*SendMessageResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{6}
}

type GetHistoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ChatUuid      string                 `protobuf:"bytes,1,opt,name=chat_uuid,json=chatUuid,proto3" json:"chat_uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetHistoryRequest) Reset() {
	*x = GetHistoryRequest{}
	mi := &file_messenger_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetHistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHistoryRequest) ProtoMessage() {}

func (x *GetHistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHistoryRequest.ProtoReflect.Descriptor instead.
func (*GetHistoryRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{7}
}

func (x *GetHistoryRequest) GetChatUuid() string {
	if x != nil {
		return x.ChatUuid
	}
	return ""
}

type ChatMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionUuid   string                 `protobuf:"bytes,1,opt,name=session_uuid,json=sessionUuid,proto3" json:"session_uuid,omitempty"`
	MessageUuid   string                 `protobuf:"bytes,2,opt,name=message_uuid,json=messageUuid,proto3" json:"message_uuid,omitempty"`
	Text          string                 `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_messenger_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{8}
}

func (x *ChatMessage) GetSessionUuid() string {
	if x != nil {
		return x.SessionUuid
	}
	return ""
}

func (x *ChatMessage) GetMessageUuid() string {
	if x != nil {
		return x.MessageUuid
	}
	return ""
}

func (x *ChatMessage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type GetHistoryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Messages      []*ChatMessage         `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetHistoryResponse) Reset() {
	*x = GetHistoryResponse{}
	mi := &file_messenger_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetHistoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHistoryResponse) ProtoMessage() {}

func (x *GetHistoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHistoryResponse.ProtoReflect.Descriptor instead.
func (*GetHistoryResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{9}
}

func (x *GetHistoryResponse) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type GetActiveChatsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetActiveChatsRequest) Reset() {
	*x = GetActiveChatsRequest{}
	mi := &file_messenger_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetActiveChatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveChatsRequest) ProtoMessage() {}

func (x *GetActiveChatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveChatsRequest.ProtoReflect.Descriptor instead.
func (*GetActiveChatsRequest) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{10}
}

type GetActiveChatsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Chats         []*Chat                `protobuf:"bytes,1,rep,name=chats,proto3" json:"chats,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetActiveChatsResponse) Reset() {
	*x = GetActiveChatsResponse{}
	mi := &file_messenger_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetActiveChatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetActiveChatsResponse) ProtoMessage() {}

func (x *GetActiveChatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_messenger_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetActiveChatsResponse.ProtoReflect.Descriptor instead.
func (*GetActiveChatsResponse) Descriptor() ([]byte, []int) {
	return file_messenger_proto_rawDescGZIP(), []int{11}
}

func (x *GetActiveChatsResponse) GetChats() []*Chat {
	if x != nil {
		return x.Chats
	}
	return nil
}

var File_messenger_proto protoreflect.FileDescriptor

var file_messenger_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x22, 0x14, 0x0a, 0x12,
	0x49, 0x6e, 0x69, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x38, 0x0a, 0x13, 0x49, 0x6e, 0x69, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x22, 0x58, 0x0a, 0x04,
	0x43, 0x68, 0x61, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x61,
	0x64, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65,
	0x61, 0x64, 0x4f, 0x6e, 0x6c, 0x79, 0x22, 0x38, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x04, 0x63,
	0x68, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73,
	0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x04, 0x63, 0x68, 0x61, 0x74,
	0x22, 0x31, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x61, 0x74, 0x55,
	0x75, 0x69, 0x64, 0x22, 0x6e, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x68, 0x61,
	0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68,
	0x61, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x30, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x61, 0x74, 0x55, 0x75, 0x69, 0x64, 0x22, 0x67, 0x0a, 0x0b,
	0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x55, 0x75, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x48, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22,
	0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x68, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3f, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x68,
	0x61, 0x74, 0x52, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x32, 0x9b, 0x03, 0x0a, 0x10, 0x4d, 0x65,
	0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4c,
	0x0a, 0x0b, 0x49, 0x6e, 0x69, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x12, 0x1c, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65,
	0x6e, 0x67, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67,
	0x65, 0x72, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65,
	0x72, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74,
	0x6f, 0x72, 0x79, 0x12, 0x1c, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x47, 0x65,
	0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x55, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x68, 0x61,
	0x74, 0x73, 0x12, 0x20, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x68, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x6e, 0x67, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x43, 0x68, 0x61, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messenger_proto_rawDescOnce sync.Once
	file_messenger_proto_rawDescData = file_messenger_proto_rawDesc
)

func file_messenger_proto_rawDescGZIP() []byte {
	file_messenger_proto_rawDescOnce.Do(func() {
		file_messenger_proto_rawDescData = protoimpl.X.CompressGZIP(file_messenger_proto_rawDescData)
	})
	return file_messenger_proto_rawDescData
}

var file_messenger_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_messenger_proto_goTypes = []any{
	(*InitSessionRequest)(nil),     // 0: messenger.InitSessionRequest
	(*InitSessionResponse)(nil),    // 1: messenger.InitSessionResponse
	(*Chat)(nil),                   // 2: messenger.Chat
	(*CreateChatRequest)(nil),      // 3: messenger.CreateChatRequest
	(*CreateChatResponse)(nil),     // 4: messenger.CreateChatResponse
	(*SendMessageRequest)(nil),     // 5: messenger.SendMessageRequest
	(*SendMessageResponse)(nil),    // 6: messenger.SendMessageResponse
	(*GetHistoryRequest)(nil),      // 7: messenger.GetHistoryRequest
	(*ChatMessage)(nil),            // 8: messenger.ChatMessage
	(*GetHistoryResponse)(nil),     // 9: messenger.GetHistoryResponse
	(*GetActiveChatsRequest)(nil),  // 10: messenger.GetActiveChatsRequest
	(*GetActiveChatsResponse)(nil), // 11: messenger.GetActiveChatsResponse
}
var file_messenger_proto_depIdxs = []int32{
	2,  // 0: messenger.CreateChatRequest.chat:type_name -> messenger.Chat
	8,  // 1: messenger.GetHistoryResponse.messages:type_name -> messenger.ChatMessage
	2,  // 2: messenger.GetActiveChatsResponse.chats:type_name -> messenger.Chat
	0,  // 3: messenger.MessengerService.InitSession:input_type -> messenger.InitSessionRequest
	3,  // 4: messenger.MessengerService.CreateChat:input_type -> messenger.CreateChatRequest
	5,  // 5: messenger.MessengerService.SendMessage:input_type -> messenger.SendMessageRequest
	7,  // 6: messenger.MessengerService.GetHistory:input_type -> messenger.GetHistoryRequest
	10, // 7: messenger.MessengerService.GetActiveChats:input_type -> messenger.GetActiveChatsRequest
	1,  // 8: messenger.MessengerService.InitSession:output_type -> messenger.InitSessionResponse
	4,  // 9: messenger.MessengerService.CreateChat:output_type -> messenger.CreateChatResponse
	6,  // 10: messenger.MessengerService.SendMessage:output_type -> messenger.SendMessageResponse
	9,  // 11: messenger.MessengerService.GetHistory:output_type -> messenger.GetHistoryResponse
	11, // 12: messenger.MessengerService.GetActiveChats:output_type -> messenger.GetActiveChatsResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_messenger_proto_init() }
func file_messenger_proto_init() {
	if File_messenger_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messenger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_messenger_proto_goTypes,
		DependencyIndexes: file_messenger_proto_depIdxs,
		MessageInfos:      file_messenger_proto_msgTypes,
	}.Build()
	File_messenger_proto = out.File
	file_messenger_proto_rawDesc = nil
	file_messenger_proto_goTypes = nil
	file_messenger_proto_depIdxs = nil
}
