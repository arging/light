// Copyright 2014 li. All rights reserved.

package web

// Filter is the pre handler for http requests.
type Filter interface {
	DoFilter(c *Context, chain FilterChain)
}

type FilterChain interface {
	DoFilter(c *Context)
}
