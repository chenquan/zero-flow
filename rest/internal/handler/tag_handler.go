package handler

import (
	"context"
	"net/http"

	"github.com/chenquan/zero-flow/internal/tag"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var httpTagAttributeKey = attribute.Key("http.header.flow.tag")

func TagHandler(tagHeader string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			ctx = newBaggage(ctx, request, tagHeader)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func newBaggage(ctx context.Context, request *http.Request, tagHeader string) context.Context {
	span := trace.SpanFromContext(ctx)
	tagString := request.Header.Get(tagHeader)
	if len(tagString) == 0 {
		return ctx
	}

	logx.WithContext(ctx).Debugw("flow staining...", logx.Field(tag.Key, tagString))

	ctx = tag.ContextWithTag(ctx, tagString)
	span.SetAttributes(httpTagAttributeKey.String(tagString))

	return ctx
}
