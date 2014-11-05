// Copyright 2014 li. All rights reserved.

package db

import (
	"reflect"
)

// Result Processor
type RProcessor interface {
	GetValueHolder(columns []string) (reflect.Value, []interface{})
	ToValue(source reflect.Value) interface{}
}

type PProcessor interface {
	Proc(v interface{}) (string, []interface{})
}

type StdPProcessor struct {
	stdSQL string
	params []string
}

// Return SQL and params
func (p *StdPProcessor) Proc(v interface{}) (string, []interface{}) {
	return p.stdSQL, procParam(p.params, v)
}

type FieldIndex struct {
	index       []int        // Field index in struct. support nested struct field.
	typ         reflect.Type // Field type.
	hasNilValue bool         // If need hanle zero value.
	value       interface{}  // Default value when field is zero value.
}

type SQLStructResult struct {
	typ     reflect.Type          // Struct type.
	isPtr   bool                  // If the result is ptr
	indexes map[string]FieldIndex // For struct.key:
}

// Processor for resultMap case.
type RMProcessor struct {
	*SQLStructResult
	id            string
	columnMapping map[string]FieldIndex
}

func (p *RMProcessor) GetValueHolder(columns []string) (reflect.Value, []interface{}) {
	source := reflect.New(p.typ)
	dest := make([]interface{}, len(columns))

	for i, column := range columns {
		fieldIndex, isExist := p.columnMapping[column]
		panicfIfFalse(isExist, "Cann't find field for column %s. Statement ID: %s.", column, p.id)
		dest[i] = source.Elem().FieldByIndex(fieldIndex.index).Addr().Interface()
	}

	return source, dest
}

func (p *RMProcessor) ToValue(source reflect.Value) interface{} {
	if p.isPtr {
		source = source.Elem().Addr()
	}
	return source.Interface()
}

// Processor for resultStruct case.
type RSProcessor struct {
	*SQLStructResult
	id string
}

func (p *RSProcessor) GetValueHolder(columns []string) (reflect.Value, []interface{}) {
	source := reflect.New(p.typ)
	dest := make([]interface{}, len(columns))

	for i, column := range columns {
		fieldIndex, isExist := p.indexes[column]
		panicfIfFalse(isExist, "Cann't find field for column %s. Statement: %s.", column, p.id)
		dest[i] = source.Elem().FieldByIndex(fieldIndex.index).Addr().Interface()
	}

	return source, dest
}

func (p *RSProcessor) ToValue(source reflect.Value) interface{} {
	if p.isPtr {
		source = source.Elem().Addr()
	}
	return source.Interface()
}

type OneColumnProcessor struct {
	id  string
	typ reflect.Type
}

func (p *OneColumnProcessor) GetValueHolder(columns []string) (reflect.Value, []interface{}) {
	source := reflect.New(p.typ)
	panicfIfTrue(len(columns) != 1,
		"Columns length expected to be 1. But it was %d. Statement ID: %s.", len(columns), p.id)
	return source, []interface{}{source.Interface()}
}

func (p *OneColumnProcessor) ToValue(source reflect.Value) interface{} {
	return source.Elem().Interface()
}
