// Copyright 2014 li. All rights reserved.

package db

import (
	"bytes"
	"fmt"
	"strings"
)

func panicfIfNil(i interface{}, format string, v ...interface{}) {
	if i == nil {
		panic("light/db: " + fmt.Sprintf(format, v...))
	}
}

func panicfIfErr(err error, format string, v ...interface{}) {
	if err != nil {
		panic(fmt.Sprintf("light/db: %s. Error: %v", fmt.Sprintf(format, v...), err))
	}
}

func panicfIfTrue(t bool, format string, v ...interface{}) {
	if t {
		panic("light/db: " + fmt.Sprintf(format, v...))
	}
}

func panicfIfFalse(t bool, format string, v ...interface{}) {
	if !t {
		panic("light/db: " + fmt.Sprintf(format, v...))
	}
}

func toUpperCase(name string) string {
	if name[0] >= 'a' && name[0] <= 'z' {
		content := []byte(name)
		content[0] = content[0] - 32
		return string(content)
	}
	return name
}

func isExportedField(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}

func getOperation(sql string) Operation {
	sql = strings.TrimSpace(sql)
	if len(sql) < 6 {
		return UNKOWN
	}

	switch op := strings.ToUpper(sql[:6]); op {
	case SELECT, INSERT, UPDATE, DELETE:
		return Operation(op)
	default:
		return UNKOWN
	}
}

// Eg. passWord => pass_word, PassWord => pass_word
//     password => password, Password => password
func toUnderLineCase(name string) string {
	buffer := bytes.NewBufferString("")
	for i, c := range name {
		if c >= 'A' && c <= 'Z' && i > 0 {
			buffer.WriteRune('_')
		}
		buffer.WriteRune(c)
	}
	return strings.ToLower(buffer.String())
}

func toStandardSQL(sql string) (string, []string, error) {
	hasPair := false
	start := 0
	params := make([]string, 0)
	stdSql := bytes.NewBufferString("")

	for i, b := range sql {
		switch {
		case b == '$' && hasPair:
			params = append(params, strings.TrimSpace(string(sql[start+1:i])))
			stdSql.WriteRune('?')
			start = i
			hasPair = false
		case b == '$':
			hasPair = true
			start = i
		case !hasPair:
			stdSql.WriteRune(b)
		}
	}

	if hasPair {
		return "", nil, fmt.Errorf("light/db: invalid sql => %s", sql)
	}

	return stdSql.String(), params, nil
}
