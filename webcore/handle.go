// Copyright 2014 li. All rights reserved.

package webcore

import (
	"github.com/roverli/light/log"
	"github.com/roverli/utils/slice"
	"reflect"
	"strings"
)

var (
	invokers = make(map[string]*Invoker)
)

func Handle(url string, handler interface{}) {
	i := strings.Index(url, "/")
	if i < 0 {
		log.Errorf("light/web: bad restful httpUrl, url: %s.", url)
		return
	}

	methods := slice.MapString(strings.Split(url[:i], "|"), func(s string) string { return strings.TrimSpace(s) })
	slice.Foreach(methods, func(method string) {
		key := strings.TrimSpace(method) + "-" + url[i:]
		if _, dup := invokers[key]; dup {
			log.Warnf("light/web: duplicate httpUrl, url: %s.", url)
		} else {
			invokers[key] = toInvoker(handler)
		}
	})
	Router.Add(methods, url[i:])
}

func toInvoker(handler interface{}) *Invoker {
	v := reflect.ValueOf(handler)
	t := v.Type()

	if v.Kind() != reflect.Func {
		panic("light/web: Handler must be function kind.")
	}

	invoker := &Invoker{Func: v}
	numIn := t.NumIn()
	invoker.Args = make([]*InvokeArg, numIn)

	for i := 0; i < numIn; i++ {
		arg := &InvokeArg{Index: i}
		argType := t.In(i)

		if argType.Kind() == reflect.Ptr {
			argType = argType.Elem()
			arg.IsPtr = true
		}
		arg.Type = argType

		switch argType {
		// Excluded
		case httpRequestType, httpResponseType, httpSessionType, bindResultType:

		default:
			if argType.Kind() != reflect.Struct {
				break
			}

			for i, num := 0, argType.NumField(); i < num; i++ {

				field := argType.Field(i)
				// Not exported field.
				if char := field.Name[0]; char < 'A' || char > 'Z' {
					continue
				}

				if arg.ExportFields == nil {
					arg.ExportFields = make(map[string][]int)
				}

				// Tag wins field name.
				switch name := field.Tag.Get("$"); name {
				case "":
					arg.ExportFields[field.Name] = field.Index
				case "-": //ignore
				default:
					arg.ExportFields[name] = field.Index
				}

				// TODO: support validate
				// For Validate
				// if tag := field.Tag.Get("@"); tag != "" {
				// 	var _ tag
				// }
			}
		}

		invoker.Args[i] = arg
	}

	// TODO: support multi?
	// for i, numOut := 0, t.NumOut(); i < numOut; i++ {
	// 	switch t.Out(i).Kind() {
	// 	case reflect.String:
	// 	case reflect.
	// 	}
	// }
	return invoker
}
