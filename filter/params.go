// Copyright 2014 li. All rights reserved.

package filter

import (
	"github.com/roverli/light/log"
	"github.com/roverli/light/web"
	"os"
)

type ParamsFilter struct {
}

func (f *ParamsFilter) DoFilter(c *web.Context, chain web.FilterChain) {

	defer func() {
		if c.Req.MultipartForm != nil {
			err := c.Req.MultipartForm.RemoveAll()
			if err != nil {
				log.Warnf("light/filter: Removing multipartForm err. %v", err)
			}
		}

		for _, tmpFile := range c.Params.TmpFiles {
			err := os.Remove(tmpFile.Name())
			if err != nil {
				log.Warnf("light/filter: Remove tmpFile err. %v", err)
			}
		}
	}()

	web.ParseParams(c.Params, c.Req)
	chain.DoFilter(c)
}
