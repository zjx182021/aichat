package server

import (
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"chatgpt-web-service/proto"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	proto.UnimplementedChatServer
	config *config.Config
	log    log.ILogger
}

func NewChatService(config *config.Config,
	log log.ILogger) proto.ChatServer {
	return &ChatServer{
		config: config,
		log:    log,
	}
}
func (c *ChatServer) ChatCompletion(_ context.Context, p *proto.ChatCompletionRequest) (*proto.ChatCompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChatCompletion not implemented")
}
func (c *ChatServer) ChatCompletionStream(p *proto.ChatCompletionRequest, s grpc.ServerStreamingServer[proto.ChatCompletionStreamResponse]) error {

	return status.Errorf(codes.Unimplemented, "method ChatCompletionStream not implemented")
}
