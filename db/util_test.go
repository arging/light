// Copyright 2014 li. All rights reserved.

package db

import (
	"reflect"
	"testing"
)

func TestGetOperation(t *testing.T) {
	selectSQL := " select * FROM table1 "
	insertSQL := " insert INTO table1 values($xx$,$yy$) "
	deleteSQL := " delete FROM table1 "
	updateSQL := " update table1 SET xx=$xx$ "
	showTablesSQL := " show tables;"

	if getOperation(selectSQL) != SELECT {
		t.Fatalf("light/db: %v,should be %v type", selectSQL, SELECT)
	}
	if getOperation(deleteSQL) != DELETE {
		t.Fatalf("light/db: %v,should be %v type", deleteSQL, DELETE)
	}
	if getOperation(insertSQL) != INSERT {
		t.Fatalf("light/db: %v,should be %v type", insertSQL, INSERT)
	}
	if getOperation(updateSQL) != UPDATE {
		t.Fatalf("light/db: %v,should be %v type", updateSQL, UPDATE)
	}
	if getOperation(showTablesSQL) != UNKOWN {
		t.Fatalf("light/db: %v,should be %v type", showTablesSQL, UNKOWN)
	}
}

func TestToStandardSQL(t *testing.T) {
	sql1 := `SELECT * FROM department WHERE name=$name$ AND count>$Count$`
	stdSql1 := `SELECT * FROM department WHERE name=? AND count>?`
	stdSql, params, err := toStandardSQL(sql1)
	if err != nil {
		t.Fatalf("light/db: %v, should not be error", sql1)
	}
	if !reflect.DeepEqual(params, []string{"name", "Count"}) {
		t.Fatalf(`light/db: %v, params should be error {"name", "Count"}`, sql1)
	}
	if stdSql != stdSql1 {
		t.Fatalf(`light/db: %v, conver to stdSql should be:%v`, sql1, stdSql1)
	}

	sql2 := `SELECT * FROM department WHERE name=$name$ AND count>$Count`
	_, _, err = toStandardSQL(sql2)
	if err == nil {
		t.Fatalf("light/db: %v,should be bad sql", sql2)
	}

	sql3 := `DELETE FROM department`
	stdSql, params, err = toStandardSQL(sql3)
	if err != nil && stdSql != sql3 && len(params) != 0 {
		t.Fatalf("light/db: %v,not match expected.", stdSql)
	}
}
