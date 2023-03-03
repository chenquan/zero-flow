package zrpc

import (
	"fmt"

	"github.com/chenquan/zero-flow/internal/p2c"
	"github.com/chenquan/zero-flow/zrpc/internal/discover"
	_ "github.com/chenquan/zero-flow/zrpc/internal/resolver"
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
		Metadata string `json:",optional,env=FLOW_METADATA"`
	}
)

func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
	options = append(options,
		zrpc.WithDialOption(grpc.WithDefaultServiceConfig(svcCfg)),
	)
	return zrpc.MustNewClient(c.RpcClientConf, options...)
}

func MustNewServer(c RpcServerConf, register func(*grpc.Server)) *RpcServer {
	etcdConf := c.Etcd
	c.Etcd = discov.EtcdConf{}
	server := zrpc.MustNewServer(c.RpcServerConf, func(server *grpc.Server) {
		register(server)
		discover.MustRegisterRpc(discover.EtcdConf{
			EtcdConf: etcdConf,
			Metadata: c.Metadata,
		}, c.RpcServerConf.ListenOn)
	})

	return server
}
