// Copyright 2014 li. All rights reserved.

package mux

import (
	"github.com/roverli/utils/errors"
	"regexp"
	"strings"
)

// Piece is a url segment.
// Such as url: /home/profile, both the "home" and "profile" are pieces.
type piece struct {
	name string // Name is the parameter name when piece matching type
	// is not precise matching. Otherwise, name is the piece value.

	prefix string // Prefix for partial matching
	suffix string // Suffix for partial matching

	prio  priority       // The matching priority, also identify the matching type.
	regex *regexp.Regexp // Instance for regular piece matching
}

func initPiece(str string) (*piece, errors.Error) {

	l := strings.Index(str, lBrace)
	r := strings.LastIndex(str, rBrace)

	if l == -1 || r == -1 || r < l {
		return &piece{name: str, prio: preciseM}, nil
	}

	p := &piece{}
	fmatched := l == 0 && r == len(str)-1
	content := strings.TrimSpace(str[l+1 : r])
	regexSepIndex := strings.Index(content, regexSep)

	if len(content) == 0 {
		return nil, errors.Newf(`bad url piece: %s, no content between "(" and ")"`, str)
	}

	if !fmatched {
		p.prefix = string(strings.TrimSpace(str[0:l]))
		p.suffix = string(strings.TrimSpace(str[r+1:]))
	}

	// param matching
	if regexSepIndex == -1 {
		p.name = content
		if fmatched {
			p.prio = fparamM
		} else {
			p.prio = pparamM
		}
	} else { // regex matching
		name := strings.TrimSpace(content[:regexSepIndex])
		expr := strings.TrimSpace(content[regexSepIndex+1:])

		regex, err := regexp.Compile(expr)
		if err != nil {
			return nil, errors.Newf(`bad url piece: %s, regex expression compile error`, str)
		}

		p.name = name
		p.regex = regex

		if fmatched {
			p.prio = fregexM
		} else {
			p.prio = pregexM
		}
	}

	return p, nil
}

// Is the piece matching the str.
// Return true, when the str matching this piece.
func (p *piece) match(str string) bool {

	switch p.prio {
	case preciseM:
		return str == p.name
	case fparamM:
		return true
	case fregexM:
		return p.regex.MatchString(str)

	case pparamM, pregexM:
		strLen := len(str)
		preLen := len(p.prefix)
		sufLen := len(p.suffix)

		if strLen < preLen+sufLen {
			return false
		}

		partMatch := str[:preLen] == p.prefix &&
			str[(strLen-sufLen):] == p.suffix

		if p.prio == pparamM {
			return partMatch
		} else {
			return partMatch &&
				p.regex.Match([]byte(str[preLen:strLen-sufLen]))
		}
	}
	panic("Never happen!")
}

// Is need to parse the piece param.
func (p *piece) isParseParam() bool {
	return p.prio != preciseM && p.name != ""
}

// Parse the piece params.
// Return the parameter name and corresponding value.
func (p *piece) parseParam(str string) (string, string) {
	if p.prio == pparamM || p.prio == pregexM {
		return p.name, str[len(p.prefix) : len(str)-len(p.suffix)]
	}

	return p.name, str
}

// Compare the matching priority.
// Return postive value, when current piece has higher priority than other piece.
// Return zero value, when current piece has same priority to other piece.
// Return negative value, when current piece has lower priority to other piece.
func (p *piece) compare(other *piece) int {
	if other == nil || p.prio > other.prio {
		return _HIGH
	} else if p.prio == other.prio {
		return _EQUAL
	} else {
		return _LOW
	}
}
