package grpcpool

import (
	"context"
	"google.golang.org/grpc"
)

type PooledGrpcClient struct {
	pool *ConnPool
}

func NewPooledGrpcClient(opt *Options) *PooledGrpcClient {
	return &PooledGrpcClient{NewConnPool(opt)}
}

func (p *PooledGrpcClient) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	conn, err := p.pool.Get(ctx)
	if err != nil {
		return err
	}

	defer func() {
		p.pool.Put(ctx, conn)
	}()

	return conn.Invoke(ctx, method, args, reply, opts...)
}

func (p *PooledGrpcClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn, err := p.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		p.pool.Put(ctx, conn)
	}()

	return conn.NewStream(ctx, desc, method, opts...)
}

func (p *PooledGrpcClient) Close() error {
	return p.pool.Close()
}
