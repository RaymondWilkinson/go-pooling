package grpcpool

import (
	"context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type Conn struct {
	usedAt   int64 // atomic
	grpcConn *grpc.ClientConn

	Inited    bool
	pooled    bool
	createdAt time.Time
}

func NewConn(grpcConn *grpc.ClientConn) *Conn {
	cn := &Conn{
		grpcConn:  grpcConn,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

func (cn *Conn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return cn.grpcConn.Invoke(ctx, method, args, reply, opts...)
}

func (cn *Conn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return cn.grpcConn.NewStream(ctx, desc, method, opts...)
}

func (cn *Conn) Close() error {
	return cn.grpcConn.Close()
}
