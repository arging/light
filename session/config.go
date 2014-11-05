// Copyright 2014 li. All rights reserved.

package session

import (
	"encoding/json"
)

// Session config
type Config struct {
	Store             string
	CookieName        string
	Domain            string
	EnableCookie      bool
	MaxAge            int // for cookie
	Secure            bool
	Hash              string
	Sed               string
	HttpOnly          bool
	GcInterval        int    // in second
	MaxActiveInterval int    // in second
	Extern            string // Json string for customer store
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return "Config[" + string(bytes) + "]"
}
