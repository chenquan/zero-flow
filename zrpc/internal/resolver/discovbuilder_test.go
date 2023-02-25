package internal

import (
	"fmt"
	"net/url"
	"testing"
)

func TestName(t *testing.T) {
	parse, err := url.Parse("rpc://192.168.1.9:8080?color=1")
	fmt.Println(err)
	fmt.Println(parse.EscapedFragment())
	fmt.Println(parse.EscapedPath())
	fmt.Println(parse.String())
	fmt.Println(parse.Scheme + ":" + parse.Opaque)
	fmt.Println(parse, err)
}
