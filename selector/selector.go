package selector

import (
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

const trafficSelect = "traffic-select"

var (
	selectorMap       = make(map[string]Selector)
	colorAttributeKey = attribute.Key("selector.color")
)

type (
	Selector interface {
		Select(conns []Conn, info balancer.PickInfo) []Conn
		Name() string
	}
	Conn interface {
		Address() resolver.Address
		SubConn() balancer.SubConn
		Tag() string
	}
)

func Register(selector Selector) {
	selectorMap[selector.Name()] = selector
}

func Get(name string) Selector {
	if b, ok := selectorMap[name]; ok {
		return b
	}
	return nil
}
