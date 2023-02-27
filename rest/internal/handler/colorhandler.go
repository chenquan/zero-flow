package handler

import (
	"net/http"

	"github.com/chenquan/zero-flow/md"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var httpColorAttributeKey = attribute.Key("http.header.color")

func ColorHandler(defaultMd md.Metadata) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			ctx, _ = md.NewMetaDataFromContext(ctx, defaultMd)

			colors := request.Header.Values("color")
			if len(colors) != 0 {
				ctx, _ = md.NewMetaDataFromContext(ctx, md.Metadata{"color": colors})
			}

			span := trace.SpanFromContext(ctx)
			span.SetAttributes(httpColorAttributeKey.StringSlice(colors))
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
