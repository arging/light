// Copyright 2014 li. All rights reserved.

package mux

// Url matching priority for piece.
// Piece has five kinds of priority, from high to low:
// preciseM, pregexM, pparamM, fregexM, fparamM.
type priority byte

const (
	fparamM  priority = iota // Fully param matching: /home/(id)
	fregexM                  // Fully regex matching: /home/(id:^123$)
	pparamM                  // Partial param matching: /home/page(id)
	pregexM                  // Partial regex matching: /home/page(id:^123$)
	preciseM                 // Precise matching: /home/profile
)

const (
	pathSep  = "/" // Url path separator
	lBrace   = "(" // Left brace
	rBrace   = ")" // Right brace
	regexSep = ":" // Seperator for regex key and value.
)

const (
	_LOW   = -1
	_EQUAL = 0
	_HIGH  = 1
)
