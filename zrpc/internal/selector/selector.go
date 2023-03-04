package selector

import (
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var tagAttributeKey = attribute.Key("selector.tag")

type (
	Selector interface {
		Select(conns []Conn, info balancer.PickInfo) []Conn
	}
	Conn interface {
		Address() resolver.Address
		SubConn() balancer.SubConn
		Tag() string
	}
)
