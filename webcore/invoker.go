// Copyright 2014 li. All rights reserved.

package webcore

import (
	"github.com/roverli/light/bind"
	_ "github.com/roverli/light/log"
	"github.com/roverli/light/session"
	"github.com/roverli/light/web"
	"net/http"
	"reflect"
)

var (
	httpRequestType  = reflect.TypeOf(http.Request{})
	httpResponseType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	httpSessionType  = reflect.TypeOf((*session.Session)(nil)).Elem()
	bindResultType   = reflect.TypeOf(BindResult{})
)

type BindResult struct {
}

type Invoker struct {
	Args []*InvokeArg
	Func reflect.Value
}

type InvokeResult struct {
	Values []reflect.Value
	Model  map[string]interface{} //.For view rendering
}

func (invoker *Invoker) Invoke(c *web.Context) *InvokeResult {

	invokeResult := &InvokeResult{}

	in := make([]reflect.Value, len(invoker.Args))
	for _, arg := range invoker.Args {
		v := reflect.New(arg.Type).Elem()
		switch arg.Type {
		case reflect.TypeOf(http.Request{}):
			v.Set(reflect.ValueOf(c.Req).Elem())
		case httpResponseType:
			v.Set(reflect.ValueOf(c.Resp))
		case httpSessionType:
			v.Set(reflect.ValueOf(c.Session))
		case bindResultType:
			v.Set(reflect.New(arg.Type))
		default:
			switch arg.Type.Kind() {
			case reflect.Map:
				mapArg := reflect.MakeMap(arg.Type)
				invokeResult.Model = mapArg.Interface().(map[string]interface{})
				v.Set(mapArg)

			case reflect.Struct:
				for name, index := range arg.ExportFields {
					r := bind.Bind(c.Params, name, arg.Type.FieldByIndex(index).Type)
					v.FieldByIndex(index).Set(r)
				}
			}
		}

		if arg.IsPtr {
			v = v.Addr()
		}
		in[arg.Index] = v
	}

	invokeResult.Values = invoker.Func.Call(in)

	return invokeResult
}

type InvokeArg struct {
	Index        int
	Type         reflect.Type
	IsPtr        bool
	ExportFields map[string][]int
}
