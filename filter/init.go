// Copyright 2014 li. All rights reserved.

package filter

import (
	"github.com/roverli/light/webcore"
)

func init() {
	webcore.Register(&PanicFilter{})
	webcore.Register(&RouteFilter{})
	webcore.Register(&ParamsFilter{})
	webcore.Register(&SessionFilter{})
}
