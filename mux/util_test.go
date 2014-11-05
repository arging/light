// Copyright 2014 li. All rights reserved.

package mux

import (
	"reflect"
	"testing"
)

func TestSplitTrim(t *testing.T) {
	rs1 := reflect.DeepEqual([]string{"abc", "efg"}, splitTrim(" abc/efg ", "/"))
	rs2 := reflect.DeepEqual([]string{"abc", "efg"}, splitTrim("////  abc/ efg/ // /", "/"))
	rs3 := reflect.DeepEqual([]string{"abc", "efg"}, splitTrim("/ abc/ //efg  /", "/"))
	rs4 := reflect.DeepEqual([]string{"a bc", "e fg"}, splitTrim("/ a bc/ //e fg  /", "/"))

	if !(rs1 && rs2 && rs3 && rs4) {
		t.Error("splitTrim not correct.")
	}
}

func TestSortPaths(t *testing.T) {
	p1, _ := initPath("/home/profile1")
	p2, _ := initPath("/home/profile(id:^[1-9]*$)")
	p3, _ := initPath("/home/profile(id)")
	p4, _ := initPath("/home/(all:^[a-z]*$)")
	p5, _ := initPath("/home/(all)")

	paths := []*path{p3, p2, p1, p4, p5}
	sortPaths(paths)
	reflect.DeepEqual(paths, []*path{p1, p2, p3, p4, p5})
}
