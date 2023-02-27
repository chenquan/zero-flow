package md

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

const zeroFlowFieldsKey = "zero-flow-fields"

var _ Carrier = (*GrpcMetadataCarrier)(nil)

type GrpcMetadataCarrier metadata.MD

func (h GrpcMetadataCarrier) Append(key string, values ...string) {
	if len(values) == 0 {
		return
	}

	key = strings.ToLower(key)
	h[key] = append(h[key], values...)
}

func (h GrpcMetadataCarrier) Get(key string) []string {
	key = strings.ToLower(key)
	return h[key]
}

func (h GrpcMetadataCarrier) Set(key string, value ...string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h GrpcMetadataCarrier) Keys() []string {
	return h[zeroFlowFieldsKey]
}
