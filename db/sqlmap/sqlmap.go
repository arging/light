// Copyright 2014 li. All rights reserved.

package sqlmap

import (
	"encoding/xml"
)

type SQLMap struct {
	XMLName xml.Name `xml:"sqlMap"`

	Inserts  []ActionInsert  `xml:"insert"`
	Selects  []ActionSelect  `xml:"select"`
	Updates  []ActionUpdate  `xml:"update"`
	Deletes  []ActionDelete  `xml:"delete"`
	Operates []ActionOperate `xml:"operate"`

	ResultMaps []ResultMap `xml:"resultMap"`
	// ParamMaps  []SQLParamMap `xml:"paramMap"`
}

type ResultMap struct {
	XMLName    xml.Name         `xml:"resultMap"`
	Id         string           `xml:"id,attr"`
	Struct     string           `xml:"struct,attr"`
	Properties []ResultProperty `xml:"result"`
}

type ResultProperty struct {
	XMLName  xml.Name `xml:"result"`
	Property string   `xml:"property,attr"`
	Column   string   `xml:"column,attr"`
	GoType   string   `xml:"gotype,attr"`
	// DBType   string   `xml:"dbtype,attr"`
	NilValue string `xml:"nil,attr"`
}

type SQLParamMap struct {
	XMLName xml.Name `xml:"resultMap"`
}

// Common for ActionOperate/ActionInsert/ActionSelect/ActionUpdate/ActionDelete
type Action struct {
	Id           string `xml:"id,attr"`
	ResultStruct string `xml:"resultStruct,attr"`
	ResultMap    string `xml:"resultMap,attr"`
	//ParamStruct  string `xml:"paramStruct,attr"`
	//ParamMap     string `xml:"paramMap,attr"`
	Dynamicer string `xml:"dynamicer,attr"`
	SQL       string `xml:",chardata"`
}

type ActionOperate struct {
	Action
	XMLName xml.Name `xml:"operate"`
}

type ActionInsert struct {
	Action
	XMLName xml.Name `xml:"insert"`
}

type ActionSelect struct {
	Action
	XMLName xml.Name `xml:"select"`
}

type ActionUpdate struct {
	Action
	XMLName xml.Name `xml:"update"`
}

type ActionDelete struct {
	Action
	XMLName xml.Name `xml:"delete"`
}

func ReadSQLMap(data string) (*SQLMap, error) {
	v := &SQLMap{}
	return v, xml.Unmarshal([]byte(data), v)
}
