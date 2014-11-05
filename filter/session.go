// Copyright 2014 li. All rights reserved.

package filter

import (
	"github.com/roverli/light/log"
	"github.com/roverli/light/web"
	"github.com/roverli/light/webcore"
)

type SessionFilter struct {
}

// TODO compatibe not use session case.
func (f *SessionFilter) DoFilter(c *web.Context, chain web.FilterChain) {

	session, err := webcore.SessionManager.Get(c.Req)
	if err != nil {
		log.Errorf("light/filter: Get session error. %v", err)
		//c.Status = web.Code500
		return
	}

	if session == nil {
		session, err = webcore.SessionManager.Create(c.Req, c.Resp)
	}
	if err != nil {
		log.Errorf("light/filter: Create session error. %v", err)
		return
	}

	c.Session = session
	chain.DoFilter(c)
}
