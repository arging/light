// Copyright 2014 li. All rights reserved.

package light

import (
	"github.com/roverli/light/conf"
	"github.com/roverli/light/db"
	_ "github.com/roverli/light/filter"
	"github.com/roverli/light/hook"
	"github.com/roverli/light/log"
	_ "github.com/roverli/light/session/memory"
	"github.com/roverli/light/webcore"
	"net/http"
)

const Version = "0.1.0"

var (
	AppName     = conf.App["app"]
	HttpAddr    = conf.App["httpAddr"]
	HttpPort    = conf.App["httpPort"]
	HttpUrl     = conf.App.String("httpUrl", "/")
	HttpSsl     = conf.App.Bool("httpSsl", false)
	HttpSslCert = conf.App["httpSslCert"]
	HttpSslKey  = conf.App["HttpSslKey"]
	ServeStatic = conf.App.Bool("serveStatic", false)
	StaticUrl   = conf.App["staticUrl"]
	StaticDir   = conf.App["staticDir"]
)

// Handle registers the handler for the given restful pattern.
func Handle(url string, handler interface{}) {
	webcore.Handle(url, handler)
}

func StartHttp() {
	defer func() {
		hook.ShutDown()
	}()

	log.Infof("Start application %s, light framework version %s.\n", AppName, Version)

	db.Start()
	webcore.Start()
	hook.Start()

	if ServeStatic {
		http.Handle(StaticUrl, http.StripPrefix(StaticUrl, http.FileServer(http.Dir(conf.ROOT+StaticUrl))))
	}
	http.HandleFunc(HttpUrl, func(w http.ResponseWriter, r *http.Request) {
		webcore.Invoke(w, r)
	})

	err := http.ListenAndServe(HttpAddr+":"+HttpPort, nil)
	if err != nil {
		log.Errorf("Start application %s fail.\n", AppName, err)
	}
}
