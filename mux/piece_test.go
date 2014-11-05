// Copyright 2014 li. All rights reserved.

package mux

import (
	"light/errors"
	"reflect"
	"testing"
)

func assertTrue(rs bool, msg string, t *testing.T) {
	if !rs {
		// Track the test error source.
		t.Error(errors.New(msg))
	}
}
func assertFalse(rs bool, msg string, t *testing.T) {
	assertTrue(!rs, msg, t)
}

func TestPieceInit(t *testing.T) {
	//normal case
	p1, _ := initPiece("match")
	rs1 := reflect.DeepEqual(p1, &piece{name: "match", prio: preciseM})
	assertTrue(rs1, "case p1", t)

	p2, _ := initPiece("( id )")
	rs2 := reflect.DeepEqual(p2, &piece{name: "id", prio: fparamM})
	assertTrue(rs2, "case p2", t)

	p3, _ := initPiece("page ( id )")
	rs3 := reflect.DeepEqual(p3, &piece{name: "id", prefix: "page", prio: pparamM})
	assertTrue(rs3, "case p3", t)

	p4, _ := initPiece("page ( id ) Num")
	rs4 := reflect.DeepEqual(p4, &piece{name: "id", prefix: "page", suffix: "Num", prio: pparamM})
	assertTrue(rs4, "case p4", t)

	p5, _ := initPiece(`(:^ab.*c$)`)
	rs5 := p5.prio == fregexM && p5.name == "" &&
		p5.regex.MatchString("abxxc") && !p5.regex.MatchString("123")
	assertTrue(rs5, "case p5", t)

	p6, _ := initPiece(`(name : ^ab.*c$)`)
	rs6 := p6.name == "name" && p6.prio == fregexM &&
		p6.regex.MatchString("abxxc") && !p6.regex.MatchString("123")
	assertTrue(rs6, "case p6", t)

	p7, _ := initPiece(` page (name : ^ab.*c$) num `)
	rs7 := p7.name == "name" && p7.prio == pregexM &&
		p7.regex.MatchString("abxxc") && !p7.regex.MatchString("123") &&
		p7.suffix == "num" && p7.prefix == "page"
	assertTrue(rs7, "case p7", t)

	p8, _ := initPiece(` page (:^ab.*c$) num `)
	rs8 := p8.name == "" && p8.prio == pregexM &&
		p8.regex.MatchString("abxxc") && !p8.regex.MatchString("123") &&
		p8.suffix == "num" && p8.prefix == "page"
	assertTrue(rs8, "case p8", t)

	//exceptional case
	_, err9 := initPiece(`( )`)
	assertTrue(err9 != nil, "case p9", t)
}

func TestPieceMatch(t *testing.T) {

	//normal case
	p1, _ := initPiece("match")
	assertTrue(p1.match("match"), "case p1", t)
	assertFalse(p1.match("mat ch"), "case p1", t)
	assertFalse(p1.match("match "), "case p1", t)

	p2, _ := initPiece("( id )")
	assertTrue(p2.match("123"), "case p2", t)
	assertTrue(p2.match(" 123x"), "case p2", t)

	p3, _ := initPiece("page ( id )")
	assertTrue(p3.match("page123"), "case p3", t)
	assertTrue(p3.match("page 123x"), "case p3", t)
	assertFalse(p3.match("abc_page_123x"), "case p3", t)

	p4, _ := initPiece("page ( id ) Num")
	assertTrue(p4.match("page123Num"), "case p3", t)
	assertTrue(p4.match("page 123xNum"), "case p3", t)
	assertFalse(p4.match("page123N"), "case p3", t)
	assertFalse(p4.match("pag123Num"), "case p3", t)

	p5, _ := initPiece(`(:^ab.*c$)`)
	assertTrue(p5.match("abxxc") && !p5.match("123"), "case p5", t)

	p6, _ := initPiece(`(name : ^ab.*c$)`)
	assertTrue(p6.match("abxxc") && !p6.match("123"), "case p6", t)

	p7, _ := initPiece(` page (name : ^ab.*c$) num `)
	assertTrue(p7.match("pageabxxcnum"), "case p7", t)
	assertFalse(p7.match("pagabxxcnum"), "case p7", t)
	assertFalse(p7.match("pageabxxcnu"), "case p7", t)
	assertFalse(p7.match("page123num"), "case p7", t)

	p8, _ := initPiece(` page (:^ab.*c$) num `)
	assertTrue(p8.match("pageabxxcnum"), "case p7", t)
	assertFalse(p8.match("pagabxxcnum"), "case p7", t)
	assertFalse(p8.match("pageabxxcnu"), "case p7", t)
	assertFalse(p8.match("page123num"), "case p7", t)
}

func TestPieceParseParam(t *testing.T) {
	p1, _ := initPiece("match")
	assertFalse(p1.isParseParam(), "case p1", t)

	p2, _ := initPiece("( id )")
	k2, v2 := p2.parseParam("123")
	assertTrue(k2 == "id", "case p2", t)
	assertTrue(v2 == "123", "case p2", t)

	p3, _ := initPiece("page ( id )")
	k3, v3 := p3.parseParam("page123")
	assertTrue(k3 == "id", "case p3", t)
	assertTrue(v3 == "123", "case p3", t)

	p4, _ := initPiece("page ( id ) Num")
	k4, v4 := p4.parseParam("page123Num")
	assertTrue(k4 == "id", "case p4", t)
	assertTrue(v4 == "123", "case p4", t)

	p5, _ := initPiece(`(:^ab.*c$)`)
	assertFalse(p5.isParseParam(), "case p5", t)

	p6, _ := initPiece(`(name : ^ab.*c$)`)
	k6, v6 := p6.parseParam("abxxc")
	assertTrue(k6 == "name", "case p6", t)
	assertTrue(v6 == "abxxc", "case p6", t)

	p7, _ := initPiece(` page (name : ^ab.*c$) num `)
	k7, v7 := p7.parseParam("pageabxxcnum")
	assertTrue(k7 == "name", "case p7", t)
	assertTrue(v7 == "abxxc", "case p7", t)

	p8, _ := initPiece(` page (:^ab.*c$) num `)
	assertFalse(p8.isParseParam(), "case p8", t)
}
