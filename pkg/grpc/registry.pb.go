// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: pkg/grpc/registry.proto

package grpc

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

type UDPNetworkType int32

const (
	UDPNetworkType_udp  UDPNetworkType = 0
	UDPNetworkType_udp4 UDPNetworkType = 1
	UDPNetworkType_udp6 UDPNetworkType = 2
)

// Enum value maps for UDPNetworkType.
var (
	UDPNetworkType_name = map[int32]string{
		0: "udp",
		1: "udp4",
		2: "udp6",
	}
	UDPNetworkType_value = map[string]int32{
		"udp":  0,
		"udp4": 1,
		"udp6": 2,
	}
)

func (x UDPNetworkType) Enum() *UDPNetworkType {
	p := new(UDPNetworkType)
	*p = x
	return p
}

func (x UDPNetworkType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UDPNetworkType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_grpc_registry_proto_enumTypes[0].Descriptor()
}

func (UDPNetworkType) Type() protoreflect.EnumType {
	return &file_pkg_grpc_registry_proto_enumTypes[0]
}

func (x UDPNetworkType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UDPNetworkType.Descriptor instead.
func (UDPNetworkType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_grpc_registry_proto_rawDescGZIP(), []int{0}
}

type RegisterPeerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *RegisterPeerResponse) Reset() {
	*x = RegisterPeerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_registry_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterPeerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterPeerResponse) ProtoMessage() {}

func (x *RegisterPeerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_registry_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterPeerResponse.ProtoReflect.Descriptor instead.
func (*RegisterPeerResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_registry_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterPeerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RegisterPeerResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type UnregisterPeerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey string `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *UnregisterPeerRequest) Reset() {
	*x = UnregisterPeerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_registry_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterPeerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterPeerRequest) ProtoMessage() {}

func (x *UnregisterPeerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_registry_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterPeerRequest.ProtoReflect.Descriptor instead.
func (*UnregisterPeerRequest) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_registry_proto_rawDescGZIP(), []int{1}
}

func (x *UnregisterPeerRequest) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

type UnregisterPeerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UnregisterPeerResponse) Reset() {
	*x = UnregisterPeerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_registry_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnregisterPeerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterPeerResponse) ProtoMessage() {}

func (x *UnregisterPeerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_registry_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterPeerResponse.ProtoReflect.Descriptor instead.
func (*UnregisterPeerResponse) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_registry_proto_rawDescGZIP(), []int{2}
}

func (x *UnregisterPeerResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey                          string         `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	AllowedIps                         []string       `protobuf:"bytes,2,rep,name=allowed_ips,json=allowedIps,proto3" json:"allowed_ips,omitempty"`
	EndpointUdpType                    UDPNetworkType `protobuf:"varint,3,opt,name=endpoint_udp_type,json=endpointUdpType,proto3,enum=UDPNetworkType" json:"endpoint_udp_type,omitempty"`
	Endpoint                           string         `protobuf:"bytes,4,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	PresharedKey                       string         `protobuf:"bytes,5,opt,name=preshared_key,json=presharedKey,proto3" json:"preshared_key,omitempty"`
	PersistentKeepaliveIntervalSeconds uint32         `protobuf:"varint,6,opt,name=persistent_keepalive_interval_seconds,json=persistentKeepaliveIntervalSeconds,proto3" json:"persistent_keepalive_interval_seconds,omitempty"`
	ProtocolVersion                    uint64         `protobuf:"varint,7,opt,name=protocol_version,json=protocolVersion,proto3" json:"protocol_version,omitempty"`
	ReceiveBytes                       int64          `protobuf:"varint,8,opt,name=receive_bytes,json=receiveBytes,proto3" json:"receive_bytes,omitempty"`
	TransmitBytes                      int64          `protobuf:"varint,9,opt,name=transmit_bytes,json=transmitBytes,proto3" json:"transmit_bytes,omitempty"`
	LastHandshakeTimeUnixSec           int64          `protobuf:"varint,10,opt,name=last_handshake_time_unix_sec,json=lastHandshakeTimeUnixSec,proto3" json:"last_handshake_time_unix_sec,omitempty"`
	Address                            string         `protobuf:"bytes,11,opt,name=address,proto3" json:"address,omitempty"`
	Zone                               string         `protobuf:"bytes,12,opt,name=zone,proto3" json:"zone,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_registry_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_registry_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_registry_proto_rawDescGZIP(), []int{3}
}

func (x *Peer) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *Peer) GetAllowedIps() []string {
	if x != nil {
		return x.AllowedIps
	}
	return nil
}

func (x *Peer) GetEndpointUdpType() UDPNetworkType {
	if x != nil {
		return x.EndpointUdpType
	}
	return UDPNetworkType_udp
}

func (x *Peer) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Peer) GetPresharedKey() string {
	if x != nil {
		return x.PresharedKey
	}
	return ""
}

func (x *Peer) GetPersistentKeepaliveIntervalSeconds() uint32 {
	if x != nil {
		return x.PersistentKeepaliveIntervalSeconds
	}
	return 0
}

func (x *Peer) GetProtocolVersion() uint64 {
	if x != nil {
		return x.ProtocolVersion
	}
	return 0
}

func (x *Peer) GetReceiveBytes() int64 {
	if x != nil {
		return x.ReceiveBytes
	}
	return 0
}

func (x *Peer) GetTransmitBytes() int64 {
	if x != nil {
		return x.TransmitBytes
	}
	return 0
}

func (x *Peer) GetLastHandshakeTimeUnixSec() int64 {
	if x != nil {
		return x.LastHandshakeTimeUnixSec
	}
	return 0
}

func (x *Peer) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Peer) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

var File_pkg_grpc_registry_proto protoreflect.FileDescriptor

var file_pkg_grpc_registry_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x14, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x36, 0x0a, 0x15, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x32, 0x0a,
	0x16, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x22, 0xfc, 0x03, 0x0a, 0x04, 0x50, 0x65, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x5f, 0x69, 0x70, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x49, 0x70, 0x73, 0x12, 0x3b, 0x0a, 0x11, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x75, 0x64, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x55, 0x44, 0x50, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x55, 0x64, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x73,
	0x68, 0x61, 0x72, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x51, 0x0a, 0x25, 0x70, 0x65, 0x72, 0x73,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x65, 0x70, 0x61, 0x6c, 0x69, 0x76, 0x65,
	0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x22, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74,
	0x65, 0x6e, 0x74, 0x4b, 0x65, 0x65, 0x70, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x76, 0x61, 0x6c, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x74, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x6d, 0x69, 0x74, 0x42, 0x79, 0x74,
	0x65, 0x73, 0x12, 0x3e, 0x0a, 0x1c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x73,
	0x65, 0x63, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x18, 0x6c, 0x61, 0x73, 0x74, 0x48, 0x61,
	0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69, 0x78, 0x53,
	0x65, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x7a, 0x6f, 0x6e, 0x65,
	0x2a, 0x2d, 0x0a, 0x0e, 0x55, 0x44, 0x50, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x75, 0x64, 0x70, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x75,
	0x64, 0x70, 0x34, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x75, 0x64, 0x70, 0x36, 0x10, 0x02, 0x32,
	0x78, 0x0a, 0x05, 0x50, 0x65, 0x65, 0x72, 0x73, 0x12, 0x2c, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x12, 0x05, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x1a,
	0x15, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x12, 0x16, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x50, 0x65, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x16, 0x5a, 0x14, 0x70, 0x32, 0x70,
	0x2d, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x61, 0x79, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpc_registry_proto_rawDescOnce sync.Once
	file_pkg_grpc_registry_proto_rawDescData = file_pkg_grpc_registry_proto_rawDesc
)

func file_pkg_grpc_registry_proto_rawDescGZIP() []byte {
	file_pkg_grpc_registry_proto_rawDescOnce.Do(func() {
		file_pkg_grpc_registry_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpc_registry_proto_rawDescData)
	})
	return file_pkg_grpc_registry_proto_rawDescData
}

var file_pkg_grpc_registry_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_grpc_registry_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_grpc_registry_proto_goTypes = []interface{}{
	(UDPNetworkType)(0),            // 0: UDPNetworkType
	(*RegisterPeerResponse)(nil),   // 1: RegisterPeerResponse
	(*UnregisterPeerRequest)(nil),  // 2: UnregisterPeerRequest
	(*UnregisterPeerResponse)(nil), // 3: UnregisterPeerResponse
	(*Peer)(nil),                   // 4: Peer
}
var file_pkg_grpc_registry_proto_depIdxs = []int32{
	0, // 0: Peer.endpoint_udp_type:type_name -> UDPNetworkType
	4, // 1: Peers.RegisterPeer:input_type -> Peer
	2, // 2: Peers.UnregisterPeer:input_type -> UnregisterPeerRequest
	1, // 3: Peers.RegisterPeer:output_type -> RegisterPeerResponse
	3, // 4: Peers.UnregisterPeer:output_type -> UnregisterPeerResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_grpc_registry_proto_init() }
func file_pkg_grpc_registry_proto_init() {
	if File_pkg_grpc_registry_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpc_registry_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterPeerResponse); i {
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
		file_pkg_grpc_registry_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterPeerRequest); i {
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
		file_pkg_grpc_registry_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnregisterPeerResponse); i {
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
		file_pkg_grpc_registry_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
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
			RawDescriptor: file_pkg_grpc_registry_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_grpc_registry_proto_goTypes,
		DependencyIndexes: file_pkg_grpc_registry_proto_depIdxs,
		EnumInfos:         file_pkg_grpc_registry_proto_enumTypes,
		MessageInfos:      file_pkg_grpc_registry_proto_msgTypes,
	}.Build()
	File_pkg_grpc_registry_proto = out.File
	file_pkg_grpc_registry_proto_rawDesc = nil
	file_pkg_grpc_registry_proto_goTypes = nil
	file_pkg_grpc_registry_proto_depIdxs = nil
}
