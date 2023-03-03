package selector

import (
	"github.com/chenquan/zero-flow/tag"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/balancer"
)

const tagKey = "tag"

var (
	DefaultSelector          = defaultSelector{}
	_               Selector = (*defaultSelector)(nil)
)

func init() {
	Register(defaultSelector{})
}

type defaultSelector struct{}

func (d defaultSelector) Select(conns []Conn, info balancer.PickInfo) []Conn {
	tag := tag.FromContext(info.Ctx)
	if len(tag) == 0 {
		return d.getNoColorConns(conns)
	}

	newConns := make([]Conn, 0, len(conns))
	for _, conn := range conns {
		if len(conn.Tag()) == 0 {
			newConns = append(newConns, conn)
			continue
		}

		if tag == conn.Tag() {
			newConns = append(newConns, conn)
		}
	}

	if len(newConns) != 0 {
		spanCtx := trace.SpanFromContext(info.Ctx)
		spanCtx.SetAttributes(colorAttributeKey.String(tag))
		logx.WithContext(info.Ctx).Debugw("flow dyeing", logx.Field(tagKey, tag))
	}

	return newConns
}

func (d defaultSelector) Name() string {
	return "DefaultSelector"
}

func (d defaultSelector) getNoColorConns(conns []Conn) []Conn {
	var newConns []Conn
	for _, conn := range conns {
		tag := conn.Tag()
		if len(tag) == 0 {
			newConns = append(newConns, conn)
		}
	}

	return newConns
}
