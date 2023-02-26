package zrpc

import (
	"fmt"
	"log"
	"net/url"

	"github.com/chenquan/zero-flow/internal/p2c"
	"github.com/chenquan/zero-flow/md"
	"github.com/chenquan/zero-flow/selector"
	"github.com/chenquan/zero-flow/zrpc/internal/clientinterceptors"
	"github.com/chenquan/zero-flow/zrpc/internal/serverinterceptors"

	"github.com/chenquan/zero-flow/zrpc/internal/discover"
	_ "github.com/chenquan/zero-flow/zrpc/internal/resolver"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type RpcClientConf struct {
	zrpc.RpcClientConf
	Metadata string `json:",optional,env=FLOW_METADATA"`
}

type RpcServerConf struct {
	zrpc.RpcServerConf
	Metadata string `json:",optional,env=FLOW_METADATA"`
}

func MustNewClient(c RpcClientConf, options ...zrpc.ClientOption) zrpc.Client {
	query, err := url.ParseQuery(c.Metadata)
	if err != nil {
		log.Panicln(err)
	}
	metadata := md.Metadata(query)
	metadata.Merge(selector.DefaultSelectorMd)
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
	options = append(options,
		zrpc.WithUnaryClientInterceptor(clientinterceptors.UnaryMdInterceptor(metadata.Clone())),
		zrpc.WithStreamClientInterceptor(clientinterceptors.StreamMdInterceptor(metadata.Clone())),
		zrpc.WithDialOption(grpc.WithDefaultServiceConfig(svcCfg)),
	)
	return zrpc.MustNewClient(c.RpcClientConf, options...)
}

func MustNewServer(c RpcServerConf, register func(*grpc.Server)) *zrpc.RpcServer {
	etcdConf := c.Etcd
	c.Etcd = discov.EtcdConf{}
	server := zrpc.MustNewServer(c.RpcServerConf, func(server *grpc.Server) {
		register(server)
		discover.MustRegisterRpc(discover.EtcdConf{
			EtcdConf: etcdConf,
			Metadata: c.Metadata,
		}, c.RpcServerConf.ListenOn)
	})
	server.AddUnaryInterceptors(serverinterceptors.UnaryMdInterceptor)
	server.AddStreamInterceptors(serverinterceptors.StreamMdInterceptor)

	return server
}
