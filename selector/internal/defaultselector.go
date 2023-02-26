package internal

import (
	"sort"
	"strings"

	"github.com/chenquan/zero-flow/md"
	"github.com/chenquan/zero-flow/selector"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/balancer"
)

const DefaultSelector = "defaultSelector"
const colorKey = "color"

var (
	_ selector.Selector = (*defaultSelector)(nil)
)

func init() {
	selector.Register(defaultSelector{})
}

type defaultSelector struct{}

func (d defaultSelector) Select(conns []selector.Conn, info balancer.PickInfo) []selector.Conn {
	m, ok := md.FromContext(info.Ctx)
	if !ok {
		return d.getNoColorConns(conns)
	}

	clientColors := m.Get(colorKey)
	if len(clientColors) == 0 {
		return d.getNoColorConns(conns)
	}

	newConns := make([]selector.Conn, 0, len(conns))
	sort.Strings(clientColors)
	for i := len(clientColors) - 1; i >= 0; i-- {
		clientColor := clientColors[i]
		for _, conn := range conns {
			metadataFromGrpcAttributes := conn.Metadata()
			colors := metadataFromGrpcAttributes.Get(colorKey)

			if len(colors) == 0 {
				newConns = append(newConns, conn)
				continue
			}

			for _, color := range colors {
				if clientColor == color {
					newConns = append(newConns, conn)
				}
			}
		}

		if len(newConns) != 0 {
			spanCtx := trace.SpanFromContext(info.Ctx)
			spanCtx.SetAttributes(selector.ColorAttributeKey.String(clientColor))
			logx.WithContext(info.Ctx).Debugw("flow dyeing", logx.Field(colorKey, clientColor), logx.Field("candidateColors", "["+strings.Join(clientColors, ", ")+"]"))

			break
		}
	}

	return newConns
}

func (d defaultSelector) Name() string {
	return DefaultSelector
}

func (d defaultSelector) getNoColorConns(conns []selector.Conn) []selector.Conn {
	var newConns []selector.Conn
	for _, conn := range conns {
		metadataFromGrpcAttributes := conn.Metadata()
		colors := metadataFromGrpcAttributes.Get(colorKey)
		if len(colors) == 0 {
			newConns = append(newConns, conn)
		}
	}

	return newConns
}
