package grpcclient

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type ClientPool interface {
	Get() *grpc.ClientConn
	Put(*grpc.ClientConn)
}

type grpcClientPool struct {
	sync.Pool
}

func (p *grpcClientPool) Get() *grpc.ClientConn {
	c := p.Pool.Get().(*grpc.ClientConn)
	if c.GetState() == connectivity.Shutdown || c.GetState() == connectivity.TransientFailure {
		c.Close()
		c = p.Pool.New().(*grpc.ClientConn)
	}
	return c
}

func (p *grpcClientPool) Put(c *grpc.ClientConn) {
	if c.GetState() == connectivity.Shutdown || c.GetState() == connectivity.TransientFailure {
		c.Close()
		return
	}
	p.Pool.Put(c)
}
func NewClientPool(target string, opt ...grpc.DialOption) (ClientPool, error) {
	return &grpcClientPool{
		Pool: sync.Pool{
			New: func() any {
				conn, err := grpc.NewClient(target, opt...)
				if err != nil {
					panic(err)
				}
				return conn
			},
		},
	}, nil
}
