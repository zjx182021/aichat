package services

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func PusBearToToken(ctx context.Context, token string) context.Context {
	md := metadata.Pairs("Authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}
