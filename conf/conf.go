// Copyright 2014 li. All rights reserved.

package conf

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	ConfDir = "/conf/"
)

// Root is the root path of the application.
var ROOT string

// If application is debug mode.
var IsDebug bool

// For global application configuration.
var AppName = "app.conf"
var App Config

// For global session configuration.
var SessionName = "session.conf"
var Session Config

// For global dataSource configuration.
var DataSources []Config

func init() {
	binDir, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	ROOT = path.Dir(binDir)

	App, _ = Load(AppName)
	IsDebug = App.Bool("isDebug", false)

	Session, _ = Load(SessionName)

	dbs := List("/db/", func(fname string) bool {
		return strings.HasPrefix(fname, "db-") &&
			strings.HasSuffix(fname, ".conf")
	})

	DataSources = make([]Config, len(dbs))
	i := 0
	for _, f := range dbs {
		DataSources[i], _ = Load("/db/" + f)
		i++
	}
}

// Load *.conf file relative to the ROOT+ConfDir path.
func Load(fname string) (Config, error) {
	return read(concatPath(fname))

}

// Read the file data.
func Read(fname string) ([]byte, error) {
	f, err := os.Open(concatPath(fname))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// Check if file exists relative to the ROOT+ConfDir path..
func Exists(fname string) bool {
	_, err := os.Stat(concatPath(fname))
	return err == nil
}

// List the files relative to the ROOT+ConfDir path..
func List(fname string, filter func(fname string) bool) []string {
	f, err := os.Open(concatPath(fname))
	if err != nil {
		return []string{}
	}

	defer f.Close()

	flist, err := f.Readdir(-1)
	if err != nil {
		return []string{}
	}

	names := make([]string, 0)
	for _, fstat := range flist {
		if !fstat.IsDir() && filter(fstat.Name()) {
			names = append(names, fstat.Name())
		}
	}
	return names
}

func concatPath(fname string) string {
	return filepath.Clean(ROOT + ConfDir + fname)
}
