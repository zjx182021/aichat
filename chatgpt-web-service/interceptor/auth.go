package interceptor

import (
	"chatgpt-web-service/pkg/config"
	"chatgpt-web-service/pkg/log"
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CustomError struct {
	Message string
}

// 实现 error 接口的 Error 方法
func (e *CustomError) Error() string {
	return e.Message
}
func UnaryHandler(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if info.FullMethod != "google.golang.org/grpc/health/grpc_health_v1" {
		err := token_check(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}
func token_check(ctx context.Context) error {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := md["authorization"]
		if len(token) > 0 {

			noprefix_token := strings.TrimPrefix(token[0], "Bearer ")
			cfg := config.GetConfig()
			if noprefix_token == cfg.Server.AccessToken {
				return nil
			}
			return errors.New("token is invalid")
		}
		return errors.New("token is missing ")
	}
	log.My_log.Error("token is invalid")
	return errors.New("token is invalid")
}
