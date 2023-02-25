package serverinterceptors

import (
	"context"

	"github.com/chenquan/zero-flow/md"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryMdInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extractMd(ctx)

	return handler(ctx, req)
}

func StreamMdInterceptor(svr interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := extractMd(ss.Context())

	return handler(svr, &wrappedServerStream{ss: ss, ctx: ctx})
}

func extractMd(ctx context.Context) context.Context {
	incomingMd, b := metadata.FromIncomingContext(ctx)
	if b {
		ctx = md.NewContext(ctx, md.GrpcMetadataCarrier(incomingMd))
	}
	m, ok := md.FromContext(ctx)
	if ok {
		logx.WithContext(ctx).Debug("metadata:", m.String())
	}

	return ctx
}

var _ grpc.ServerStream = (*wrappedServerStream)(nil)

type wrappedServerStream struct {
	ss  grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) SetHeader(m metadata.MD) error {
	return w.ss.SetHeader(m)
}

func (w *wrappedServerStream) SendHeader(m metadata.MD) error {
	return w.ss.SendHeader(m)
}

func (w *wrappedServerStream) SetTrailer(m metadata.MD) {
	w.ss.SetTrailer(m)
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	return w.ss.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	return w.ss.RecvMsg(m)
}
