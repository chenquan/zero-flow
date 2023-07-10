package zrpc

import (
	"fmt"

	"github.com/chenquan/zero-flow/zrpc/internal/discover"
	"github.com/chenquan/zero-flow/zrpc/internal/p2c"
	"github.com/chenquan/zero-flow/zrpc/internal/resolver"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ClientOption  = zrpc.ClientOption
	Client        = zrpc.Client
	RpcServer     = zrpc.RpcServer
	RpcClientConf struct {
		zrpc.RpcClientConf
	}
	RpcServerConf struct {
		zrpc.RpcServerConf
		Tag string `json:",optional,env=ZERO_FLOW_TAG"`
	}
)

func init() {
	resolver.Register()
}

func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
	options = append([]ClientOption{zrpc.WithDialOption(grpc.WithDefaultServiceConfig(svcCfg))},
		options...,
	)
	return zrpc.MustNewClient(c.RpcClientConf, options...)
}

func MustNewServer(c RpcServerConf, register func(*grpc.Server)) *RpcServer {
	hasEtcd := c.HasEtcd()

	etcdConf := c.Etcd
	c.Etcd = discov.EtcdConf{}
	server := zrpc.MustNewServer(c.RpcServerConf, func(server *grpc.Server) {
		register(server)
		if hasEtcd {
			discover.MustRegisterRpc(discover.EtcdConf{
				EtcdConf: etcdConf,
				Tag:      c.Tag,
			}, c.RpcServerConf.ListenOn)
		}
	})

	return server
}
