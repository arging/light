// Copyright 2014 li. All rights reserved.

package hook

import (
	"github.com/roverli/light/log"
	"github.com/roverli/utils/slice"
)

type Hooks []func()

var (
	startHooks    = make(Hooks, 0)
	shutDownHooks = make(Hooks, 0)
)

func OnAppStart(f func()) {
	startHooks = append(startHooks, f)
}

func OnAppShutDown(f func()) {
	shutDownHooks = append(shutDownHooks, f)
}

// For cross package, should only be called by framework.
func Start() {
	slice.Foreach(startHooks, func(f func()) {
		deferCall(f)
	})
}

// For cross package, should only be called by framework.
func ShutDown() {
	slice.Foreach(shutDownHooks, func(f func()) {
		deferCall(f)
	})
}

func deferCall(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	f()
}
