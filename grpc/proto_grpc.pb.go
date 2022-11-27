// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: grpc/proto.proto

package Assignment5_git

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RegisterClient is the client API for Register service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterClient interface {
	JoinServer(ctx context.Context, in *Request, opts ...grpc.CallOption) (Register_JoinServerClient, error)
	PlaceBid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (*Conformation, error)
	Result(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*Auctionresult, error)
}

type registerClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterClient(cc grpc.ClientConnInterface) RegisterClient {
	return &registerClient{cc}
}

func (c *registerClient) JoinServer(ctx context.Context, in *Request, opts ...grpc.CallOption) (Register_JoinServerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Register_ServiceDesc.Streams[0], "/proto.register/joinServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &registerJoinServerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Register_JoinServerClient interface {
	Recv() (*Reply, error)
	grpc.ClientStream
}

type registerJoinServerClient struct {
	grpc.ClientStream
}

func (x *registerJoinServerClient) Recv() (*Reply, error) {
	m := new(Reply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *registerClient) PlaceBid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (*Conformation, error) {
	out := new(Conformation)
	err := c.cc.Invoke(ctx, "/proto.register/placeBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registerClient) Result(ctx context.Context, in *ResultRequest, opts ...grpc.CallOption) (*Auctionresult, error) {
	out := new(Auctionresult)
	err := c.cc.Invoke(ctx, "/proto.register/result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterServer is the server API for Register service.
// All implementations must embed UnimplementedRegisterServer
// for forward compatibility
type RegisterServer interface {
	JoinServer(*Request, Register_JoinServerServer) error
	PlaceBid(context.Context, *Bid) (*Conformation, error)
	Result(context.Context, *ResultRequest) (*Auctionresult, error)
	mustEmbedUnimplementedRegisterServer()
}

// UnimplementedRegisterServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterServer struct {
}

func (UnimplementedRegisterServer) JoinServer(*Request, Register_JoinServerServer) error {
	return status.Errorf(codes.Unimplemented, "method JoinServer not implemented")
}
func (UnimplementedRegisterServer) PlaceBid(context.Context, *Bid) (*Conformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceBid not implemented")
}
func (UnimplementedRegisterServer) Result(context.Context, *ResultRequest) (*Auctionresult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedRegisterServer) mustEmbedUnimplementedRegisterServer() {}

// UnsafeRegisterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterServer will
// result in compilation errors.
type UnsafeRegisterServer interface {
	mustEmbedUnimplementedRegisterServer()
}

func RegisterRegisterServer(s grpc.ServiceRegistrar, srv RegisterServer) {
	s.RegisterService(&Register_ServiceDesc, srv)
}

func _Register_JoinServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RegisterServer).JoinServer(m, &registerJoinServerServer{stream})
}

type Register_JoinServerServer interface {
	Send(*Reply) error
	grpc.ServerStream
}

type registerJoinServerServer struct {
	grpc.ServerStream
}

func (x *registerJoinServerServer) Send(m *Reply) error {
	return x.ServerStream.SendMsg(m)
}

func _Register_PlaceBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Bid)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterServer).PlaceBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.register/placeBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterServer).PlaceBid(ctx, req.(*Bid))
	}
	return interceptor(ctx, in, info, handler)
}

func _Register_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.register/result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterServer).Result(ctx, req.(*ResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Register_ServiceDesc is the grpc.ServiceDesc for Register service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Register_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.register",
	HandlerType: (*RegisterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "placeBid",
			Handler:    _Register_PlaceBid_Handler,
		},
		{
			MethodName: "result",
			Handler:    _Register_Result_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "joinServer",
			Handler:       _Register_JoinServer_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/proto.proto",
}
