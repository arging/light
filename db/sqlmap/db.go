// Copyright 2014 li. All rights reserved.

package sqlmap

import (
	"encoding/xml"
)

type DB struct {
	XMLName   xml.Name   `xml:"db"`
	Name      string     `xml:"name,attr"`
	Props     []DBProp   `xml:"property"`
	Locations []Location `xml:"sqlmap"`
}

type Location struct {
	Resource string `xml:"resource,attr"`
}

type DBProp struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func ReadDB(data string) (*DB, error) {
	v := &DB{}
	return v, xml.Unmarshal([]byte(data), v)
}
