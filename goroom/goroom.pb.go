// Code generated by protoc-gen-go.
// source: goroom.proto
// DO NOT EDIT!

/*
Package goroom is a generated protocol buffer package.

It is generated from these files:
	goroom.proto

It has these top-level messages:
	CreateRequest
	CreateReply
	JoinRequest
	JoinReply
	InfoRequest
	Room
	InfoReply
*/
package goroom

import proto "github.com/golang/protobuf/proto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type ResponseStatus int32

const (
	ResponseStatus_SUCCESS         ResponseStatus = 0
	ResponseStatus_ROOM_FULL       ResponseStatus = 1
	ResponseStatus_ROOM_NOT_FOUND  ResponseStatus = 2
	ResponseStatus_ALREADY_IN_ROOM ResponseStatus = 3
)

var ResponseStatus_name = map[int32]string{
	0: "SUCCESS",
	1: "ROOM_FULL",
	2: "ROOM_NOT_FOUND",
	3: "ALREADY_IN_ROOM",
}
var ResponseStatus_value = map[string]int32{
	"SUCCESS":         0,
	"ROOM_FULL":       1,
	"ROOM_NOT_FOUND":  2,
	"ALREADY_IN_ROOM": 3,
}

func (x ResponseStatus) String() string {
	return proto.EnumName(ResponseStatus_name, int32(x))
}

type CreateRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id" json:"user_id,omitempty"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}

type CreateReply struct {
	Status ResponseStatus `protobuf:"varint,1,opt,name=status,enum=goroom.ResponseStatus" json:"status,omitempty"`
	RoomId string         `protobuf:"bytes,2,opt,name=room_id" json:"room_id,omitempty"`
}

func (m *CreateReply) Reset()         { *m = CreateReply{} }
func (m *CreateReply) String() string { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()    {}

type JoinRequest struct {
	RoomId string `protobuf:"bytes,1,opt,name=room_id" json:"room_id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id" json:"user_id,omitempty"`
}

func (m *JoinRequest) Reset()         { *m = JoinRequest{} }
func (m *JoinRequest) String() string { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()    {}

type JoinReply struct {
	Status ResponseStatus `protobuf:"varint,1,opt,name=status,enum=goroom.ResponseStatus" json:"status,omitempty"`
}

func (m *JoinReply) Reset()         { *m = JoinReply{} }
func (m *JoinReply) String() string { return proto.CompactTextString(m) }
func (*JoinReply) ProtoMessage()    {}

type InfoRequest struct {
	RoomId string `protobuf:"bytes,1,opt,name=room_id" json:"room_id,omitempty"`
}

func (m *InfoRequest) Reset()         { *m = InfoRequest{} }
func (m *InfoRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()    {}

type Room struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Room) Reset()         { *m = Room{} }
func (m *Room) String() string { return proto.CompactTextString(m) }
func (*Room) ProtoMessage()    {}

type InfoReply struct {
	Status ResponseStatus `protobuf:"varint,1,opt,name=status,enum=goroom.ResponseStatus" json:"status,omitempty"`
	Room   *Room          `protobuf:"bytes,2,opt,name=room" json:"room,omitempty"`
}

func (m *InfoReply) Reset()         { *m = InfoReply{} }
func (m *InfoReply) String() string { return proto.CompactTextString(m) }
func (*InfoReply) ProtoMessage()    {}

func (m *InfoReply) GetRoom() *Room {
	if m != nil {
		return m.Room
	}
	return nil
}

func init() {
	proto.RegisterEnum("goroom.ResponseStatus", ResponseStatus_name, ResponseStatus_value)
}

// Client API for RoomManager service

type RoomManagerClient interface {
	CreateRoom(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
	JoinRoom(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error)
	RoomInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error)
}

type roomManagerClient struct {
	cc *grpc.ClientConn
}

func NewRoomManagerClient(cc *grpc.ClientConn) RoomManagerClient {
	return &roomManagerClient{cc}
}

func (c *roomManagerClient) CreateRoom(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := grpc.Invoke(ctx, "/goroom.RoomManager/CreateRoom", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomManagerClient) JoinRoom(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error) {
	out := new(JoinReply)
	err := grpc.Invoke(ctx, "/goroom.RoomManager/JoinRoom", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomManagerClient) RoomInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error) {
	out := new(InfoReply)
	err := grpc.Invoke(ctx, "/goroom.RoomManager/RoomInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RoomManager service

type RoomManagerServer interface {
	CreateRoom(context.Context, *CreateRequest) (*CreateReply, error)
	JoinRoom(context.Context, *JoinRequest) (*JoinReply, error)
	RoomInfo(context.Context, *InfoRequest) (*InfoReply, error)
}

func RegisterRoomManagerServer(s *grpc.Server, srv RoomManagerServer) {
	s.RegisterService(&_RoomManager_serviceDesc, srv)
}

func _RoomManager_CreateRoom_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(CreateRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(RoomManagerServer).CreateRoom(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _RoomManager_JoinRoom_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(JoinRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(RoomManagerServer).JoinRoom(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _RoomManager_RoomInfo_Handler(srv interface{}, ctx context.Context, buf []byte) (proto.Message, error) {
	in := new(InfoRequest)
	if err := proto.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(RoomManagerServer).RoomInfo(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _RoomManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "goroom.RoomManager",
	HandlerType: (*RoomManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _RoomManager_CreateRoom_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _RoomManager_JoinRoom_Handler,
		},
		{
			MethodName: "RoomInfo",
			Handler:    _RoomManager_RoomInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
