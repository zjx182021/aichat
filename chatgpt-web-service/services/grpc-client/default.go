package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DefaultClient struct{}
type ServiceClient interface {
	GetPool(addr string) ClientPool
}

func (c *DefaultClient) GetPool(target string) ClientPool {
	pool, err := NewClientPool(target, c.getopt()...)
	if err != nil {
		panic(err)
	}
	return pool
}

func (c *DefaultClient) getopt() []grpc.DialOption {
	myopt := make([]grpc.DialOption, 0)
	myopt = append(myopt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return myopt
}
