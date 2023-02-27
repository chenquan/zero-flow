package md

import (
	"context"
	"net/url"
	"strings"

	"github.com/chenquan/zero-flow/internal/xstrings"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/metadata"
)

var _ Carrier = (*Metadata)(nil)

type (
	Metadata    map[string][]string
	metadataKey struct{}
)

func (m Metadata) Append(key string, values ...string) {
	if len(values) == 0 {
		return
	}

	key = strings.ToLower(key)
	m[key] = append(m[key], values...)
}

func (m Metadata) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, strings.ToLower(k))
	}

	return keys
}

func (m Metadata) Range(f func(key string, values ...string) bool) {
	for key, value := range m {
		key := strings.ToLower(key)
		if !f(key, value...) {
			break
		}
	}
}

func (m Metadata) Set(key string, values ...string) {
	key = strings.ToLower(key)
	m[key] = values
}

func (m Metadata) Get(key string) []string {
	key = strings.ToLower(key)
	return m[key]
}

func (m Metadata) Delete(key string) {
	key = strings.ToLower(key)
	delete(m, key)
}

func (m Metadata) String() string {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for k, values := range m {
		k = strings.ToLower(k)
		builder.WriteString(k)
		builder.WriteRune('=')
		builder.WriteRune('[')
		if len(values) != 0 {
			builder.WriteString(values[0])
			for _, value := range values[1:] {
				builder.WriteString(", ")
				builder.WriteString(value)
			}
		}

		builder.WriteRune(']')
	}
	builder.WriteRune('}')
	return builder.String()
}

func (m Metadata) Clone() Metadata {
	metadata := make(Metadata, len(m))
	for k, v := range m {
		k = strings.ToLower(k)
		metadata[k] = copyOf(v)
	}

	return metadata
}

func (m Metadata) Merge(metadata Metadata) {
	metadata.Range(func(key string, values ...string) bool {
		m.Append(strings.ToLower(key), copyOf(xstrings.Distinct(values))...)
		return true
	})
}

func (m Metadata) Distinct() {
	for k, v := range m {
		m[k] = xstrings.Distinct(v)
	}
}

func FromContext(ctx context.Context) (Metadata, bool) {
	value := ctx.Value(metadataKey{})
	if value == nil {
		return nil, false
	}

	return value.(Metadata).Clone(), true
}

func NewContext(ctx context.Context, carrier Carrier) context.Context {
	md := Metadata{}
	for _, k := range carrier.Keys() {
		md[strings.ToLower(k)] = carrier.Get(k)
	}

	return context.WithValue(ctx, metadataKey{}, md)
}

func ToGrpcMetadata(m Metadata) metadata.MD {
	md := metadata.MD{}
	keys := m.Keys()
	for _, key := range keys {
		md.Set(key, m.Get(key)...)
	}
	m.Set(zeroFlowFieldsKey, keys...)

	return md
}

func NewMetaDataFromContext(ctx context.Context, carrier Carrier) (context.Context, Metadata) {
	metadata, ok := FromContext(ctx)
	if !ok {
		metadata = Metadata{}
	} else {
		metadata = metadata.Clone()
	}
	keys := carrier.Keys()
	for _, key := range keys {
		metadata.Append(strings.ToLower(key), xstrings.Distinct(carrier.Get(key))...)
	}
	metadata.Distinct()
	return context.WithValue(ctx, metadataKey{}, metadata), metadata
}

func copyOf(v []string) []string {
	vals := make([]string, len(v))
	copy(vals, v)
	return vals
}

type metadataAttributesKey struct{}

func NewAttributes(v any) *attributes.Attributes {
	return attributes.New(metadataAttributesKey{}, v)
}

func FromGrpcAttributes(attributes *attributes.Attributes) (Metadata, bool) {
	value := attributes.Value(metadataAttributesKey{})
	if value == nil {
		return nil, false
	}

	m, ok := value.(url.Values)
	if !ok {
		return nil, false
	}

	return Metadata(m).Clone(), true
}
