// Copyright 2014 li. All rights reserved.

package webcore

import (
	"github.com/roverli/light/view"
	"github.com/roverli/light/web"
)

var (
	filters    []web.Filter
	tmpFilters []web.Filter
)

func Register(f web.Filter) {
	tmpFilters = append(tmpFilters, f)
}

type filterChain struct {
	filters []web.Filter
	index   int
}

func (chain *filterChain) DoFilter(c *web.Context) {
	if chain.index < len(chain.filters) {
		index := chain.index
		chain.index++
		chain.filters[index].DoFilter(c, chain)
	}
}

func newChain(filters []web.Filter) web.FilterChain {
	return &filterChain{filters: filters}
}

type InvokeFilter struct {
}

func (f *InvokeFilter) DoFilter(c *web.Context, chain web.FilterChain) {
	invoker := invokers[c.Req.Method+"-"+c.RouteResult.Url]
	r := invoker.Invoke(c)

	// TODO support mutlti returns
	view.Render(r.Values[0].String(), r.Model, c.Resp)
}
