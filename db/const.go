// Copyright 2014 li. All rights reserved.

package db

import (
	"reflect"
	"sync"
	"time"
)

// Datasources by mappaing
var (
	dataSources = make(map[string]*DB)
	dbMutex     sync.Mutex // Mutex for datasources
)

// Oringin go type mapping for resultStrut
var goTypeMapping map[string]reflect.Type
var (
	TimeType = reflect.TypeOf(time.Time{})
)

// SQL operation type
type Operation string

const (
	UNKOWN Operation = "UNKNOW"
	INSERT           = "INSERT"
	SELECT           = "SELECT"
	UPDATE           = "UPDATE"
	DELETE           = "DELETE"
)

// Directory db configs exists.
const (
	DBLocation = "/db/"
)

func init() {
	goTypeMapping = make(map[string]reflect.Type)
	goTypeMapping["int8"] = reflect.TypeOf(int8(0))
	goTypeMapping["int16"] = reflect.TypeOf(int16(0))
	goTypeMapping["int32"] = reflect.TypeOf(int32(0))
	goTypeMapping["int64"] = reflect.TypeOf(int64(0))
	goTypeMapping["int"] = reflect.TypeOf(int(0))
	goTypeMapping["uint8"] = reflect.TypeOf(uint8(0))
	goTypeMapping["uint16"] = reflect.TypeOf(uint16(0))
	goTypeMapping["uint32"] = reflect.TypeOf(uint32(0))
	goTypeMapping["uint64"] = reflect.TypeOf(uint64(0))
	goTypeMapping["uint"] = reflect.TypeOf(uint(0))
	goTypeMapping["string"] = reflect.TypeOf("")
}
