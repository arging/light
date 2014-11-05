// Copyright 2014 li. All rights reserved.

package filter

import (
	"github.com/roverli/light/web"
	"github.com/roverli/light/webcore"
)

type RouteFilter struct{}

func (f *RouteFilter) DoFilter(c *web.Context, chain web.FilterChain) {
	c.RouteResult = webcore.Router.Route(c.Req.Method, c.Req.URL.Path)

	switch c.RouteResult.IsMatch {
	case true:
		c.Params.Route = c.RouteResult.Parse()
		chain.DoFilter(c)

	case false:
		//c.Status = web.Code404
	}
}
