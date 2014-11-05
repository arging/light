// Copyright 2014 li. All rights reserved.

package sqlmap

import (
	"testing"
)

func TestReadDBConfig(t *testing.T) {
	var data = `
	<?xml version="1.0" encoding="UTF-8"?>
	<db name="training">
		<property name="driver" value="mysql"/>
		<sqlmap resource="xxx-common-sqlmap.xml"/>
	</db>
	`
	v, err := ReadDB(data)
	if err != nil || v.Name != "training" ||
		len(v.Props) != 1 ||
		v.Props[0].Name != "driver" ||
		v.Props[0].Value != "mysql" ||
		len(v.Locations) != 1 ||
		v.Locations[0].Resource != "xxx-common-sqlmap.xml" {
		t.Fatalf("DBConfig xml unmarshal not match. xml:%s \n struct:%v", data, v)
	}
}
