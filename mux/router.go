// Copyright 2014 li. All rights reserved.

// Package route is the router module for light framework.
// It supports restful url matching and param extracting.
package mux

import (
	"github.com/roverli/utils/errors"
	"net/url"
)

var _ Router = &restRouter{}

// Router for url routing.
// Before call Route method, you should start the Router.
type Router interface {

	// Get the Router name
	Name() string

	// Add route url by specified methods.
	Add(methods []string, url string)

	// Start the router.
	Start() errors.Error

	// Route for the corresponding method and url,and resolve the params.
	Route(method string, url string) *Result
}

// The routing result.
// If the url doesn't match any predefined path, "IsMatch" equals false;
// Otherwise,"IsMatch" will be true, and the "Url" string is the matched predefined path.
// Even the "IsMatch" equals true, "Params" can be nil(the path doesn't need to be resloved).
// So before use the params result, check whether params is nil first.
type Result struct {
	IsMatch bool
	Url     string
	pieces  []string
	path    *path
}

// Parse pramas.
func (r Result) Parse() url.Values {
	if r.IsMatch {
		return r.path.parseParams(r.pieces)
	}
	return nil
}

// Create a router by name.
func New(name string) Router {
	return &restRouter{name: name}
}

// Defined for origin url path.
type routeUrl struct {
	methods []string
	url     string
}

// Restful style struct for for Router interface
type restRouter struct {
	name       string
	routeUrls  []routeUrl
	urlMapping map[string]map[int][]*path
}

func (router *restRouter) Name() string {
	return router.name
}

func (router *restRouter) Add(methods []string, url string) {
	router.routeUrls = append(router.routeUrls, routeUrl{methods, url})
}

func (router *restRouter) Start() errors.Error {
	urlMapping := make(map[string]map[int][]*path)
	for _, routeUrl := range router.routeUrls {
		p, err := initPath(routeUrl.url)
		if err != nil {
			return errors.Wrapf(err, "restRouter url error: %s.", routeUrl.url)
		}

		methods := routeUrl.methods
		if len(methods) == 0 {
			methods = []string{""}
		}

		for _, method := range methods {
			pathMap, ok := urlMapping[method]
			if !ok {
				pathMap = make(map[int][]*path)
				urlMapping[method] = pathMap
			}

			pathArr, _ := pathMap[p.depth]
			pathMap[p.depth] = append(pathArr, p)
		}
	}

	for _, v := range urlMapping {
		for _, paths := range v {
			sortPaths(paths)
		}
	}
	router.urlMapping = urlMapping
	return nil
}

func (router *restRouter) Route(method string, url string) *Result {

	routes := router.urlMapping[method]
	if routes == nil {
		return &Result{}
	}

	var target *path
	strs := splitTrim(url, pathSep)
	depth := len(strs)

	for depth >= 0 {
		if paths := routes[depth]; paths != nil {
			for _, p := range paths {
				if p.match(strs) {
					target = p
					break
				}
			}
			if target != nil {
				break
			}
		}
		depth--
	}

	if target == nil {
		return &Result{}
	}

	return &Result{IsMatch: true, Url: target.origin, pieces: strs, path: target}
}
