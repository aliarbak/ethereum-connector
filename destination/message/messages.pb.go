// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.3
// source: messages.message

package destmessage

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

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number     string `protobuf:"bytes,1,opt,name=Number,proto3" json:"Number,omitempty"`
	Block_Hash string `protobuf:"bytes,2,opt,name=Block_Hash,json=BlockHash,proto3" json:"Block_Hash,omitempty"`
	Chain_Id   int64  `protobuf:"varint,3,opt,name=Chain_Id,json=ChainId,proto3" json:"Chain_Id,omitempty"`
	Time       int64  `protobuf:"varint,4,opt,name=Time,proto3" json:"Time,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Block) GetBlock_Hash() string {
	if x != nil {
		return x.Block_Hash
	}
	return ""
}

func (x *Block) GetChain_Id() int64 {
	if x != nil {
		return x.Chain_Id
	}
	return 0
}

func (x *Block) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chain_Id         int64  `protobuf:"varint,1,opt,name=Chain_Id,json=ChainId,proto3" json:"Chain_Id,omitempty"`
	Block_Number     string `protobuf:"bytes,2,opt,name=Block_Number,json=BlockNumber,proto3" json:"Block_Number,omitempty"`
	Block_Hash       string `protobuf:"bytes,3,opt,name=Block_Hash,json=BlockHash,proto3" json:"Block_Hash,omitempty"`
	Block_Time       int64  `protobuf:"varint,4,opt,name=Block_Time,json=BlockTime,proto3" json:"Block_Time,omitempty"`
	Contract_Address string `protobuf:"bytes,5,opt,name=Contract_Address,json=ContractAddress,proto3" json:"Contract_Address,omitempty"`
	Tx_To            string `protobuf:"bytes,6,opt,name=Tx_To,json=TxTo,proto3" json:"Tx_To,omitempty"`
	Tx_Hash          string `protobuf:"bytes,7,opt,name=Tx_Hash,json=TxHash,proto3" json:"Tx_Hash,omitempty"`
	Index            int64  `protobuf:"varint,8,opt,name=Index,proto3" json:"Index,omitempty"`
	Value            string `protobuf:"bytes,9,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

func (x *Transaction) GetChain_Id() int64 {
	if x != nil {
		return x.Chain_Id
	}
	return 0
}

func (x *Transaction) GetBlock_Number() string {
	if x != nil {
		return x.Block_Number
	}
	return ""
}

func (x *Transaction) GetBlock_Hash() string {
	if x != nil {
		return x.Block_Hash
	}
	return ""
}

func (x *Transaction) GetBlock_Time() int64 {
	if x != nil {
		return x.Block_Time
	}
	return 0
}

func (x *Transaction) GetContract_Address() string {
	if x != nil {
		return x.Contract_Address
	}
	return ""
}

func (x *Transaction) GetTx_To() string {
	if x != nil {
		return x.Tx_To
	}
	return ""
}

func (x *Transaction) GetTx_Hash() string {
	if x != nil {
		return x.Tx_Hash
	}
	return ""
}

func (x *Transaction) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Transaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type TransactionLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Log_Hash            string `protobuf:"bytes,1,opt,name=Log_Hash,json=LogHash,proto3" json:"Log_Hash,omitempty"`
	Chain_Id            int64  `protobuf:"varint,2,opt,name=Chain_Id,json=ChainId,proto3" json:"Chain_Id,omitempty"`
	Tx_Hash             string `protobuf:"bytes,3,opt,name=Tx_Hash,json=TxHash,proto3" json:"Tx_Hash,omitempty"`
	Index               int64  `protobuf:"varint,4,opt,name=Index,proto3" json:"Index,omitempty"`
	Event_Name          string `protobuf:"bytes,5,opt,name=Event_Name,json=EventName,proto3" json:"Event_Name,omitempty"`
	Event_Alias         string `protobuf:"bytes,6,opt,name=Event_Alias,json=EventAlias,proto3" json:"Event_Alias,omitempty"`
	Contract_Address    string `protobuf:"bytes,7,opt,name=Contract_Address,json=ContractAddress,proto3" json:"Contract_Address,omitempty"`
	Block_Number        string `protobuf:"bytes,8,opt,name=Block_Number,json=BlockNumber,proto3" json:"Block_Number,omitempty"`
	Block_Hash          string `protobuf:"bytes,9,opt,name=Block_Hash,json=BlockHash,proto3" json:"Block_Hash,omitempty"`
	Block_Time          int64  `protobuf:"varint,10,opt,name=Block_Time,json=BlockTime,proto3" json:"Block_Time,omitempty"`
	Tx_Index            int64  `protobuf:"varint,11,opt,name=Tx_Index,json=TxIndex,proto3" json:"Tx_Index,omitempty"`
	Tx_Contract_Address string `protobuf:"bytes,12,opt,name=Tx_Contract_Address,json=TxContractAddress,proto3" json:"Tx_Contract_Address,omitempty"`
	Tx_To               string `protobuf:"bytes,13,opt,name=Tx_To,json=TxTo,proto3" json:"Tx_To,omitempty"`
	Tx_Value            string `protobuf:"bytes,14,opt,name=Tx_Value,json=TxValue,proto3" json:"Tx_Value,omitempty"`
	Data                string `protobuf:"bytes,15,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *TransactionLog) Reset() {
	*x = TransactionLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionLog) ProtoMessage() {}

func (x *TransactionLog) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionLog.ProtoReflect.Descriptor instead.
func (*TransactionLog) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *TransactionLog) GetLog_Hash() string {
	if x != nil {
		return x.Log_Hash
	}
	return ""
}

func (x *TransactionLog) GetChain_Id() int64 {
	if x != nil {
		return x.Chain_Id
	}
	return 0
}

func (x *TransactionLog) GetTx_Hash() string {
	if x != nil {
		return x.Tx_Hash
	}
	return ""
}

func (x *TransactionLog) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *TransactionLog) GetEvent_Name() string {
	if x != nil {
		return x.Event_Name
	}
	return ""
}

func (x *TransactionLog) GetEvent_Alias() string {
	if x != nil {
		return x.Event_Alias
	}
	return ""
}

func (x *TransactionLog) GetContract_Address() string {
	if x != nil {
		return x.Contract_Address
	}
	return ""
}

func (x *TransactionLog) GetBlock_Number() string {
	if x != nil {
		return x.Block_Number
	}
	return ""
}

func (x *TransactionLog) GetBlock_Hash() string {
	if x != nil {
		return x.Block_Hash
	}
	return ""
}

func (x *TransactionLog) GetBlock_Time() int64 {
	if x != nil {
		return x.Block_Time
	}
	return 0
}

func (x *TransactionLog) GetTx_Index() int64 {
	if x != nil {
		return x.Tx_Index
	}
	return 0
}

func (x *TransactionLog) GetTx_Contract_Address() string {
	if x != nil {
		return x.Tx_Contract_Address
	}
	return ""
}

func (x *TransactionLog) GetTx_To() string {
	if x != nil {
		return x.Tx_To
	}
	return ""
}

func (x *TransactionLog) GetTx_Value() string {
	if x != nil {
		return x.Tx_Value
	}
	return ""
}

func (x *TransactionLog) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type RawTransactionLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Log_Hash            string `protobuf:"bytes,1,opt,name=Log_Hash,json=LogHash,proto3" json:"Log_Hash,omitempty"`
	Chain_Id            int64  `protobuf:"varint,2,opt,name=Chain_Id,json=ChainId,proto3" json:"Chain_Id,omitempty"`
	Tx_Hash             string `protobuf:"bytes,3,opt,name=Tx_Hash,json=TxHash,proto3" json:"Tx_Hash,omitempty"`
	Index               int64  `protobuf:"varint,4,opt,name=Index,proto3" json:"Index,omitempty"`
	Contract_Address    string `protobuf:"bytes,5,opt,name=Contract_Address,json=ContractAddress,proto3" json:"Contract_Address,omitempty"`
	Block_Number        string `protobuf:"bytes,6,opt,name=Block_Number,json=BlockNumber,proto3" json:"Block_Number,omitempty"`
	Block_Hash          string `protobuf:"bytes,7,opt,name=Block_Hash,json=BlockHash,proto3" json:"Block_Hash,omitempty"`
	Block_Time          int64  `protobuf:"varint,8,opt,name=Block_Time,json=BlockTime,proto3" json:"Block_Time,omitempty"`
	Tx_Index            int64  `protobuf:"varint,9,opt,name=Tx_Index,json=TxIndex,proto3" json:"Tx_Index,omitempty"`
	Tx_Contract_Address string `protobuf:"bytes,10,opt,name=Tx_Contract_Address,json=TxContractAddress,proto3" json:"Tx_Contract_Address,omitempty"`
	Tx_To               string `protobuf:"bytes,11,opt,name=Tx_To,json=TxTo,proto3" json:"Tx_To,omitempty"`
	Tx_Value            string `protobuf:"bytes,12,opt,name=Tx_Value,json=TxValue,proto3" json:"Tx_Value,omitempty"`
	Data                string `protobuf:"bytes,13,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *RawTransactionLog) Reset() {
	*x = RawTransactionLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawTransactionLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawTransactionLog) ProtoMessage() {}

func (x *RawTransactionLog) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawTransactionLog.ProtoReflect.Descriptor instead.
func (*RawTransactionLog) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

func (x *RawTransactionLog) GetLog_Hash() string {
	if x != nil {
		return x.Log_Hash
	}
	return ""
}

func (x *RawTransactionLog) GetChain_Id() int64 {
	if x != nil {
		return x.Chain_Id
	}
	return 0
}

func (x *RawTransactionLog) GetTx_Hash() string {
	if x != nil {
		return x.Tx_Hash
	}
	return ""
}

func (x *RawTransactionLog) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *RawTransactionLog) GetContract_Address() string {
	if x != nil {
		return x.Contract_Address
	}
	return ""
}

func (x *RawTransactionLog) GetBlock_Number() string {
	if x != nil {
		return x.Block_Number
	}
	return ""
}

func (x *RawTransactionLog) GetBlock_Hash() string {
	if x != nil {
		return x.Block_Hash
	}
	return ""
}

func (x *RawTransactionLog) GetBlock_Time() int64 {
	if x != nil {
		return x.Block_Time
	}
	return 0
}

func (x *RawTransactionLog) GetTx_Index() int64 {
	if x != nil {
		return x.Tx_Index
	}
	return 0
}

func (x *RawTransactionLog) GetTx_Contract_Address() string {
	if x != nil {
		return x.Tx_Contract_Address
	}
	return ""
}

func (x *RawTransactionLog) GetTx_To() string {
	if x != nil {
		return x.Tx_To
	}
	return ""
}

func (x *RawTransactionLog) GetTx_Value() string {
	if x != nil {
		return x.Tx_Value
	}
	return ""
}

func (x *RawTransactionLog) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x6d, 0x0a,
	0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d,
	0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x19, 0x0a,
	0x08, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x8e, 0x02, 0x0a,
	0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x5f, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x5f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x13, 0x0a, 0x05, 0x54, 0x78, 0x5f, 0x54, 0x6f, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x54, 0x78, 0x54, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x54, 0x78, 0x5f, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x78, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xd0, 0x03,
	0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67,
	0x12, 0x19, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x5f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x4c, 0x6f, 0x67, 0x48, 0x61, 0x73, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x43,
	0x68, 0x61, 0x69, 0x6e, 0x5f, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x43,
	0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x54, 0x78, 0x5f, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54, 0x78, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1d, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x41, 0x6c,
	0x69, 0x61, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x41, 0x6c, 0x69, 0x61, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x5f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x21, 0x0a, 0x0c, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x19, 0x0a, 0x08, 0x54, 0x78, 0x5f, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x54, 0x78, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x2e, 0x0a, 0x13,
	0x54, 0x78, 0x5f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x54, 0x78, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x13, 0x0a, 0x05,
	0x54, 0x78, 0x5f, 0x54, 0x6f, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x78, 0x54,
	0x6f, 0x12, 0x19, 0x0a, 0x08, 0x54, 0x78, 0x5f, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x54, 0x78, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x22, 0x93, 0x03, 0x0a, 0x11, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x4c, 0x6f, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x5f, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4c, 0x6f, 0x67, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x19, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x54, 0x78, 0x5f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x54,
	0x78, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x29, 0x0a, 0x10, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x5f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x48, 0x61, 0x73, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x54, 0x78, 0x5f, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x54, 0x78, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x2e, 0x0a, 0x13, 0x54, 0x78, 0x5f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x5f, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x11, 0x54, 0x78, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x13, 0x0a, 0x05, 0x54, 0x78, 0x5f, 0x54, 0x6f, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x54, 0x78, 0x54, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x54, 0x78, 0x5f, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x54, 0x78, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_messages_proto_goTypes = []interface{}{
	(*Block)(nil),             // 0: destination.Block
	(*Transaction)(nil),       // 1: destination.Transaction
	(*TransactionLog)(nil),    // 2: destination.TransactionLog
	(*RawTransactionLog)(nil), // 3: destination.RawTransactionLog
}
var file_messages_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionLog); i {
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
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawTransactionLog); i {
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
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}