// Copyright 2014 li. All rights reserved.

package mux

import (
	"sort"
	"strings"
)

// Split the string by the separator and trim the splited strings.
// If the splited string is trimed to empty, it will not add to the result.
func splitTrim(str string, sep string) []string {
	strArr := strings.Split(str, sep)
	trimStrs := make([]string, 0, len(strArr))

	for _, s := range strArr {
		trimStr := strings.TrimSpace(s)
		if trimStr != "" {
			trimStrs = append(trimStrs, trimStr)
		}
	}

	return trimStrs
}

// Struct for sort the paths.
type sortedPaths []*path

func (a sortedPaths) Len() int {
	return len(a)
}

func (a sortedPaths) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortedPaths) Less(i, j int) bool {
	return a[i].compare(a[j]) > 0
}

func sortPaths(paths []*path) {
	sort.Sort(sortedPaths(paths))
}
