// Copyright 2014 li. All rights reserved.

package webcore

import (
	_ "github.com/roverli/light/log"
	"github.com/roverli/light/web"
	//"go.net/websocket"
	"net/http"
)

func Invoke(resp http.ResponseWriter, req *http.Request) {

	c := &web.Context{
		Req:    req,
		Resp:   resp,
		Params: &web.Params{},
	}

	newChain(filters).DoFilter(c)
}
