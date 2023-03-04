package handler

import (
	"context"
	"net/http"

	"github.com/chenquan/zero-flow/tag"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var httpColorAttributeKey = attribute.Key("http.header.color")

func TagHandler(headerTag string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			ctx = newBaggage(ctx, request, headerTag)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func newBaggage(ctx context.Context, request *http.Request, headerTag string) context.Context {
	span := trace.SpanFromContext(ctx)
	tagString := request.Header.Get(headerTag)
	if len(tagString) == 0 {
		return ctx
	}
	logx.WithContext(ctx).Debugw("flow staining...", logx.Field("tag", tagString))

	ctx = tag.ContextWithTag(ctx, tagString)
	span.SetAttributes(httpColorAttributeKey.String(tagString))

	return ctx
}
