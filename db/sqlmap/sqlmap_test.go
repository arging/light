// Copyright 2014 li. All rights reserved.

package sqlmap

import (
	"testing"
)

func TestReadSQLMap(t *testing.T) {
	var data = `
	<?xml version="1.0" encoding="UTF-8"?>
	<sqlMap>
		<insert id="1">			
		</insert>
		<select id="2" paramMap="paramMap1" resultMap="resultMap1"  paramStruct="struct1" resultStruct="resutStruct1">
			 select * from user
		</select>	
		<update id="3"/>
		<delete id="4"/>
		<operate id="5" />
		<resultMap id="r1" struct="model/user">
			<result property="p1" column="c1" gotype="int" dbtype="int" nil="3"/> 
		</resultMap>
	</sqlMap>
	`

	v, err := ReadSQLMap(data)
	if err != nil || v.XMLName.Local != "sqlMap" ||
		len(v.Inserts) != 1 || v.Inserts[0].Id != "1" ||
		len(v.Selects) != 1 || v.Selects[0].Id != "2" ||
		//v.Selects[0].ParamMap != "paramMap1" ||
		v.Selects[0].ResultMap != "resultMap1" ||
		//v.Selects[0].ParamStruct != "struct1" ||
		v.Selects[0].ResultStruct != "resutStruct1" ||
		len(v.Updates) != 1 || v.Updates[0].Id != "3" ||
		len(v.Deletes) != 1 || v.Deletes[0].Id != "4" ||
		len(v.Operates) != 1 || v.Operates[0].Id != "5" ||
		len(v.ResultMaps) != 1 || v.ResultMaps[0].Id != "r1" ||
		v.ResultMaps[0].Struct != "model/user" ||
		len(v.ResultMaps[0].Properties) != 1 ||
		v.ResultMaps[0].Properties[0].Property != "p1" ||
		v.ResultMaps[0].Properties[0].Column != "c1" ||
		v.ResultMaps[0].Properties[0].GoType != "int" ||
		//v.ResultMaps[0].Properties[0].DBType != "int" ||
		v.ResultMaps[0].Properties[0].NilValue != "3" {
		t.Fatalf("SqlMapConfig xml unmarshal not match. result:%v", v)
	}
}
