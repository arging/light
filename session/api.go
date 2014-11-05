// Copyright 2014 li. All rights reserved.

package session

import (
	"github.com/roverli/light/conf"
	"github.com/roverli/light/util"
	"github.com/roverli/utils/errors"
)

var stores = make(map[string]Store)

func Register(name string, store Store) {
	util.PanicfIfTrue(stores[name] != nil, "light/session: duplicate store %s.", name)
	stores[name] = store
}

func New(config conf.Config) (*Manager, errors.Error) {
	store := stores[config["store"]]
	util.PanicfIfTrue(store == nil, "light/session: No such store %s.", config["store"])

	innerConfig := convToInnerConfig(config)
	err := store.Init(innerConfig)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		store:  store,
		config: innerConfig}
	manager.start()

	return manager, nil
}

func convToInnerConfig(config conf.Config) Config {
	return Config{
		Store:             config["store"],
		CookieName:        config.String("cookieName", "goSessionId"),
		Domain:            config["domain"],
		EnableCookie:      config.Bool("enableCookie", false),
		MaxAge:            config.Int("cookieMaxAge", 30*60),
		Secure:            config.Bool("secure", false),
		Hash:              config["hash"],
		Sed:               config["sed"],
		HttpOnly:          config.Bool("httpOnly", true),
		Extern:            config["extern"],
		GcInterval:        config.Int("gcInterval", 2*60),
		MaxActiveInterval: config.Int("maxActiveInterval", 60*10),
	}
}
