package resolver

import (
	"net/url"
	"strings"

	"github.com/chenquan/zero-flow/internal/tag"
	"github.com/chenquan/zero-flow/zrpc/internal/targets"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/resolver"
)

const (
	EndpointSepChar = ','
	subsetSize      = 32
)

func Register() {
	resolver.Register(&etcdFlowBuilder{})
}

type etcdFlowBuilder struct{}

func (b *etcdFlowBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(targets.GetAuthority(target), func(r rune) bool {
		return r == EndpointSepChar
	})

	sub, err := discov.NewSubscriber(hosts, targets.GetEndpoints(target))
	if err != nil {
		return nil, err
	}

	update := func() {
		addresses, err := parserAddr(sub)
		if err != nil {
			logx.Error(err)
		}

		if err := cc.UpdateState(resolver.State{
			Addresses: addresses,
		}); err != nil {
			logx.Error(err)
		}
	}
	sub.AddListener(update)
	update()

	return &nopResolver{cc: cc}, nil
}

func (b *etcdFlowBuilder) Scheme() string {
	return "etcd"
}

func parserAddr(sub *discov.Subscriber) ([]resolver.Address, error) {
	var addrs []resolver.Address
	for _, val := range subset(sub.Values(), subsetSize) {
		u, err := url.Parse("rpc://" + val)
		if err != nil {
			logx.Error(err)
			continue
		}

		t := u.Query().Get("tag")
		attr := tag.NewAttributes(t)
		addr := u.Host

		addrs = append(addrs, resolver.Address{
			Addr:               addr,
			BalancerAttributes: attr,
		})
	}

	return addrs, nil
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close()                                        {}
func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {}
