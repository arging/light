// Copyright 2014 li. All rights reserved.

package mux

import (
	"github.com/roverli/utils/errors"
)

// Path is representation for url .
// Such as: /article/page(num) is a path.
type path struct {
	depth  int      // equals to len(pieces)
	pieces []*piece // pieces in order
	parse  bool     // if need parse params
	origin string   // the origin url
}

func initPath(url string) (*path, errors.Error) {

	strs := splitTrim(url, pathSep)
	pieces := make([]*piece, len(strs))

	for i, v := range strs {
		piece, err := initPiece(v)
		if err != nil {
			return nil, errors.Wrapf(err, "init path error, path: %s.", url)
		}
		pieces[i] = piece
	}

	isParse := false
	for _, piece := range pieces {
		if piece.isParseParam() {
			isParse = true
			break
		}
	}
	return &path{len(pieces), pieces, isParse, url}, nil
}

func (p *path) parseParams(strs []string) (params map[string][]string) {
	if p.parse {
		params = make(map[string][]string)
		for i, str := range strs {
			if i < p.depth && p.pieces[i].isParseParam() {
				k, v := p.pieces[i].parseParam(str)
				arr, _ := params[k]
				params[k] = append(arr, v)
			}
		}
	}
	return
}

func (p *path) match(strs []string) bool {

	if p.depth > len(strs) {
		return false
	}

	for i := p.depth - 1; i >= 0; i-- {
		if !p.pieces[i].match(strs[i]) {
			return false
		}
	}
	return true
}

func (p *path) compare(other *path) int {

	if other == nil || p.depth > other.depth {
		return _HIGH
	}
	if p.depth < other.depth {
		return _LOW
	}

	for i := 0; i < p.depth; i++ {
		r := p.pieces[i].compare(other.pieces[i])
		if r != _EQUAL {
			return r
		}
	}
	return _EQUAL
}
