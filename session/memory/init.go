// Copyright 2014 li. All rights reserved.

package memory

import (
	"github.com/roverli/light/session"
)

func init() {
	session.Register("memory", &MemoryStore{})
}
