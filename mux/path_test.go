// Copyright 2014 li. All rights reserved.

package mux

import (
	"testing"
)

func TestPathInit(t *testing.T) {
	//normal case
	_, err1 := initPath("/home/profile/(userId)")
	assertTrue(err1 == nil, "case p1", t)

	//exceptional case
	_, err2 := initPath(`/home/profile/a()bc`)
	assertTrue(err2 != nil, "case p2", t)
}

func TestPathCompare(t *testing.T) {
	p0, _ := initPath(`/`)
	p1, _ := initPath(`/home/`)
	p2, _ := initPath(`/home/profile`)
	p3, _ := initPath(`/home/pro(name:fi)le`)
	p4, _ := initPath(`/home/pro(fi)le`)
	p5, _ := initPath(`/home/(name:profile)`)
	p6, _ := initPath(`/home/(profile)`)
	p7, _ := initPath(`/home/profile/userId`)
	p8, _ := initPath(`/home/profile/(userId)`)
	p9, _ := initPath(`/home/profile/(name)`)

	assertTrue(p0.compare(p1) == _LOW, "caes p0 p1", t)
	assertTrue(p1.compare(p2) == _LOW, "caes p1 p2", t)
	assertTrue(p2.compare(p3) == _HIGH, "caes p2 p3", t)
	assertTrue(p3.compare(p4) == _HIGH, "caes p3 p4", t)
	assertTrue(p4.compare(p5) == _HIGH, "caes p4 p5", t)
	assertTrue(p5.compare(p6) == _HIGH, "caes p5 p6", t)
	assertTrue(p6.compare(p7) == _LOW, "caes p6 p7", t)
	assertTrue(p7.compare(p8) == _HIGH, "caes p7 p8", t)
	assertTrue(p8.compare(p9) == _EQUAL, "caes p8 p9", t)
}

func TestPathMath(t *testing.T) {
	p0, _ := initPath(`/`)
	assertTrue(p0.match([]string{}), "case p0", t)
	assertTrue(p0.match([]string{"abc", "efg"}), "case p0", t)
	assertTrue(p0.match([]string{"abc"}), "case p0", t)

	p1, _ := initPath(`/home/`)
	assertTrue(p1.match([]string{"home"}), "case p1", t)
	assertTrue(p1.match([]string{"home", "asfasdf"}), "case p1", t)
	assertFalse(p1.match([]string{"homex"}), "case p1", t)

	p2, _ := initPath(`/home/profile`)
	assertTrue(p2.match([]string{"home", "profile"}), "case p2", t)
	assertFalse(p2.match([]string{"home", "profilex"}), "case p2", t)
	assertFalse(p2.match([]string{"home", "prof"}), "case p2", t)

	p3, _ := initPath(`/home/pro(name:^[0-9]*$)le`)
	assertTrue(p3.match([]string{"home", "pro123le"}), "case p3", t)
	assertTrue(p3.match([]string{"home", "pro123le", "photo"}), "case p3", t)
	assertFalse(p3.match([]string{"home", "proxxle"}), "case p3", t)

	p4, _ := initPath(`/home/pro(fi)le`)
	assertTrue(p4.match([]string{"home", "profile"}), "case p4", t)
	assertTrue(p4.match([]string{"home", "proxxfixxxxle"}), "case p4", t)
	assertFalse(p4.match([]string{"home", "prxxle"}), "case p4", t)

	p5, _ := initPath(`/home/(name:^[0-9]*$)`)
	assertTrue(p5.match([]string{"home", "123"}), "case p5", t)
	assertFalse(p5.match([]string{"home", "pro123le"}), "case p5", t)
	assertFalse(p5.match([]string{"home", "abc"}), "case p5", t)

	p6, _ := initPath(`/home/(profile)`)
	assertTrue(p6.match([]string{"home", "123efg"}), "case p6", t)
	assertTrue(p6.match([]string{"home", "bacxx"}), "case p6", t)
	assertTrue(p6.match([]string{"home", "123", "abc"}), "case p6", t)
	assertTrue(p6.match([]string{"home", "123"}), "case p6", t)
}

func TestPathParseParams(t *testing.T) {
	p0, _ := initPath(`/`)
	assertFalse(p0.parse, "case p0", t)

	p1, _ := initPath(`/home/`)
	assertFalse(p1.parse, "case p1", t)

	p2, _ := initPath(`/home/profile`)
	assertFalse(p2.parse, "case p2", t)

	p3, _ := initPath(`/(id)/pro(name:^[0-9]*$)le`)
	map3 := p3.parseParams([]string{"123", "pro1le"})
	v31 := map3["id"]
	v32 := map3["name"]
	assertTrue(len(v31) == 1 && v31[0] == "123", "case p3", t)
	assertTrue(len(v32) == 1 && v32[0] == "1", "case p3", t)

	p4, _ := initPath(`/(name)/page(id)`)
	map4 := p4.parseParams([]string{"tony", "page123"})
	v41 := map4["id"]
	v42 := map4["name"]
	assertTrue(len(v41) == 1 && v41[0] == "123", "case p4", t)
	assertTrue(len(v42) == 1 && v42[0] == "tony", "case p4", t)

	p5, _ := initPath(`/(id)/page(id)`)
	map5 := p5.parseParams([]string{"tony", "page123"})
	v51 := map5["id"]
	assertTrue(len(v51) == 2 && v51[0] == "tony" && v51[1] == "123", "case p4", t)
}
