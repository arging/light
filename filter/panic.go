// Copyright 2014 li. All rights reserved.

package filter

import (
	"github.com/roverli/light/log"
	"github.com/roverli/light/web"
	"github.com/roverli/utils/errors"
)

type PanicFilter struct{}

// PanicFilter protective panics. In the first order of the chain.
func (f *PanicFilter) DoFilter(c *web.Context, chain web.FilterChain) {

	defer func() {
		if err := recover(); err != nil {
			f.handlePanic(c, errors.Newf("%v", err))
		}
	}()
	chain.DoFilter(c)
}

func (f *PanicFilter) handlePanic(c *web.Context, err interface{}) {
	log.Errorf("light/filter: server inner panic. %v", err)
}
