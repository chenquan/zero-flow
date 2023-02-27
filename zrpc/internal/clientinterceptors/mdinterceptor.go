package clientinterceptors

import (
	"context"

	"github.com/chenquan/zero-flow/md"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryMdInterceptor(defaultMd md.Metadata) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = injectionMd(ctx, defaultMd)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamMdInterceptor(defaultMd md.Metadata) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = injectionMd(ctx, defaultMd)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func injectionMd(ctx context.Context, defaultMd md.Metadata) context.Context {
	ctx, m := md.NewMetaDataFromContext(ctx, defaultMd)

	outgoingMd, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		outgoingMd = metadata.MD{}
	}

	grpcMetadata := md.ToGrpcMetadata(m)
	for key, values := range grpcMetadata {
		outgoingMd.Append(key, values...)
	}
	ctx = metadata.NewOutgoingContext(ctx, outgoingMd)

	m, ok = md.FromContext(ctx)
	if ok {
		logx.WithContext(ctx).Debug("metadata:", m.String())
	}

	return ctx
}
