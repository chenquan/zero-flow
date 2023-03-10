package rest

import (
	"github.com/chenquan/zero-flow/rest/internal/handler"
	"github.com/zeromicro/go-zero/rest"
)

type (
	RunOption = rest.RunOption
	Server    = rest.Server
	RestConf  struct {
		rest.RestConf
		TagHeader string `json:",optional"`
	}
)

func MustNewServer(c RestConf, opts ...RunOption) *Server {
	server := rest.MustNewServer(c.RestConf, opts...)
	server.Use(rest.ToMiddleware(handler.TagHandler(c.TagHeader)))
	return server
}
