// Copyright 2014 li. All rights reserved.

package web

import (
	"github.com/roverli/light/log"
	"net/http"
)

func ParseParams(params *Params, r *http.Request) {

	params.Query = r.URL.Query()

	switch ResolveContentType(r) {
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			log.Warn("light/web: Error parsing request body.", err)
		} else {
			params.Form = r.Form
		}

	case "multipart/form-data":
		if err := r.ParseMultipartForm(1024 * 1024 * 10); err != nil {
			log.Warn("light/web: : Error parsing request body.", err)
		} else {
			params.Form = r.MultipartForm.Value
			params.Files = r.MultipartForm.File
		}
	}

	params.Values = params.merge()
}
