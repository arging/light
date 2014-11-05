// Copyright (c) 2014 li. All rights reserved.

package util

import (
	"fmt"
)

func PanicIfErr(err error, v ...interface{}) {
	if err != nil {
		panic(fmt.Sprintf("%s\nError:%v", fmt.Sprintln(v...), err))
	}

}

func PanicfIfErr(err error, format string, v ...interface{}) {
	if err != nil {
		panic(fmt.Sprintf("%s\nError:%v", fmt.Sprintf(format, v...), err))
	}
}

func PanicIfNil(i interface{}, v ...interface{}) {
	if i == nil {
		panic(fmt.Sprintln(v...))
	}
}

func PanicfIfNil(i interface{}, format string, v ...interface{}) {
	if i == nil {
		panic(fmt.Sprintf(format, v...))
	}
}

func PanicIfTrue(t bool, v ...interface{}) {
	if t {
		panic(fmt.Sprintln(v...))
	}
}

func PanicfIfTrue(t bool, format string, v ...interface{}) {
	if t {
		panic(fmt.Sprintf(format, v...))
	}
}
