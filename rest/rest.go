package rest

import (
	"github.com/chenquan/zero-flow/rest/internal/handler"
	"github.com/zeromicro/go-zero/rest"
)

func MustNewServer(c rest.RestConf, opts ...rest.RunOption) *rest.Server {
	server := rest.MustNewServer(c, opts...)
	server.Use(rest.ToMiddleware(handler.ColorHandler))
	return server
}
