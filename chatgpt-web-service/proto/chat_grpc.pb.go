// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: chat.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Chat_ChatCompletion_FullMethodName       = "/ai_chat_service.zvoice.com.Chat/ChatCompletion"
	Chat_ChatCompletionStream_FullMethodName = "/ai_chat_service.zvoice.com.Chat/ChatCompletionStream"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	ChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatCompletionStreamResponse], error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) ChatCompletion(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (*ChatCompletionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChatCompletionResponse)
	err := c.cc.Invoke(ctx, Chat_ChatCompletion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) ChatCompletionStream(ctx context.Context, in *ChatCompletionRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatCompletionStreamResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Chat_ServiceDesc.Streams[0], Chat_ChatCompletionStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ChatCompletionRequest, ChatCompletionStreamResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_ChatCompletionStreamClient = grpc.ServerStreamingClient[ChatCompletionStreamResponse]

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility.
type ChatServer interface {
	ChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(*ChatCompletionRequest, grpc.ServerStreamingServer[ChatCompletionStreamResponse]) error
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatServer struct{}

func (UnimplementedChatServer) ChatCompletion(context.Context, *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChatCompletion not implemented")
}
func (UnimplementedChatServer) ChatCompletionStream(*ChatCompletionRequest, grpc.ServerStreamingServer[ChatCompletionStreamResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ChatCompletionStream not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}
func (UnimplementedChatServer) testEmbeddedByValue()              {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	// If the following call pancis, it indicates UnimplementedChatServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_ChatCompletion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatCompletionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).ChatCompletion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_ChatCompletion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).ChatCompletion(ctx, req.(*ChatCompletionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_ChatCompletionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChatCompletionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServer).ChatCompletionStream(m, &grpc.GenericServerStream[ChatCompletionRequest, ChatCompletionStreamResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_ChatCompletionStreamServer = grpc.ServerStreamingServer[ChatCompletionStreamResponse]

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai_chat_service.zvoice.com.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ChatCompletion",
			Handler:    _Chat_ChatCompletion_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ChatCompletionStream",
			Handler:       _Chat_ChatCompletionStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}
