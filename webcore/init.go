// Copyright 2014 li. All rights reserved.

package webcore

import (
	"github.com/roverli/light/conf"
	"github.com/roverli/light/log"
	"github.com/roverli/light/mux"
	"github.com/roverli/light/session"
	"github.com/roverli/utils/errors"
	"sync"
)

var (
	Router         mux.Router = mux.New("LightRouter")
	SessionManager *session.Manager
)

var initOnce sync.Once

// Call by light package.
// So clients have chance to register their own SessionStore, Filter...
func Start() errors.Error {
	var err errors.Error
	initOnce.Do(func() {
		config, err := conf.Load("session.conf")
		if err != nil {
			log.Error("light/web: Load session config error. Error: %v", err)
		} else {
			SessionManager, err = session.New(config)
			if err != nil {
				log.Error("light/web: Init session error. Error: %v", err)
			}
		}

		Router.Start()

		// Ensure InvokerFilter is at the last of the chain.
		filters = tmpFilters
		tmpFilters = nil
		filters = append(filters, &InvokeFilter{})
	})

	return err
}
