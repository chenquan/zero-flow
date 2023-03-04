package selector

import (
	"github.com/chenquan/zero-flow/internal/tag"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/balancer"
)

var (
	DefaultSelector          = defaultSelector{}
	_               Selector = (*defaultSelector)(nil)
)

type defaultSelector struct{}

func (d defaultSelector) Select(conns []Conn, info balancer.PickInfo) []Conn {
	tagString := tag.FromContext(info.Ctx)
	if len(tagString) == 0 {
		return d.getNoTagConns(conns)
	}

	newConns := make([]Conn, 0, len(conns))
	for _, conn := range conns {
		if tagString == conn.Tag() {
			newConns = append(newConns, conn)
		}
	}

	if len(newConns) != 0 {
		logx.WithContext(info.Ctx).Debugw("flow staining...", logx.Field(tag.Key, tagString))

		spanCtx := trace.SpanFromContext(info.Ctx)
		spanCtx.SetAttributes(tagAttributeKey.String(tagString))
	}

	return newConns
}

func (d defaultSelector) getNoTagConns(conns []Conn) []Conn {
	var newConns []Conn
	for _, conn := range conns {
		if len(conn.Tag()) == 0 {
			newConns = append(newConns, conn)
		}
	}

	return newConns
}
