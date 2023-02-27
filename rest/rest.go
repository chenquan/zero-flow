package rest

import (
	"log"
	"net/url"

	"github.com/chenquan/zero-flow/md"
	"github.com/chenquan/zero-flow/rest/internal/handler"
	"github.com/chenquan/zero-flow/selector"
	"github.com/zeromicro/go-zero/rest"
)

type (
	RunOption = rest.RunOption
	Server    = rest.Server
	RestConf  struct {
		rest.RestConf
		Metadata string `json:",optional,env=FLOW_METADATA"`
	}
)

func MustNewServer(c RestConf, opts ...RunOption) *Server {
	query, err := url.ParseQuery(c.Metadata)
	if err != nil {
		log.Panicln(err)
	}
	metadata := md.Metadata(query)
	metadata.Merge(selector.DefaultSelectorMd)
	metadata.Distinct()

	server := rest.MustNewServer(c.RestConf, opts...)
	server.Use(rest.ToMiddleware(handler.ColorHandler(metadata)))
	return server
}
