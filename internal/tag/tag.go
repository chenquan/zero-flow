package tag

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/baggage"
	"google.golang.org/grpc/attributes"
)

const Key = "x-zero-flow-tag"

func ContextWithTag(ctx context.Context, tag string) context.Context {
	bg := baggage.FromContext(ctx)
	member, err := baggage.NewMember(Key, tag)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return ctx
	}

	bg, err = bg.SetMember(member)
	if err != nil {
		logx.WithContext(ctx).Error(err)
		return ctx
	}

	ctx = baggage.ContextWithBaggage(ctx, bg)

	return ctx
}

func FromContext(ctx context.Context) string {
	bg := baggage.FromContext(ctx)
	member := bg.Member(Key)

	return member.Value()
}

type tagAttributesKey struct{}

func NewAttributes(tag string) *attributes.Attributes {
	return attributes.New(tagAttributesKey{}, tag)
}

func FromGrpcAttributes(attributes *attributes.Attributes) (string, bool) {
	value := attributes.Value(tagAttributesKey{})
	if value == nil {
		return "", false
	}

	m, ok := value.(string)
	if !ok {
		return "", false
	}

	return m, true
}
