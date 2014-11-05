// Copyright 2014 li. All rights reserved.

package db

import (
	"database/sql"
	"encoding/xml"
	"github.com/roverli/light/conf"
	"github.com/roverli/light/db/sqlmap"
	"github.com/roverli/light/log"
	"github.com/roverli/utils/slice"
	"strconv"
	"strings"
	"sync"
)

var once sync.Once // For config init

func Start() {
	once.Do(func() {
		log.Info("light/db: DB initial begin ......")
		initDB(readDB())
		log.Info("light/db: DB initial end   ......")
	})

}

func readDB() []*sqlmap.DB {
	fNames := conf.List(DBLocation, func(fname string) bool {
		return strings.HasPrefix(fname, "db-") &&
			strings.HasSuffix(fname, ".xml")
	})

	log.Infof("light/db: find db files: %v.", fNames)

	configs := make([]*sqlmap.DB, len(fNames))
	for i, fName := range fNames {
		data, err := conf.Read(DBLocation + fName)
		panicfIfErr(err, "read db config error. DB: %s.", fName)

		config, err := sqlmap.ReadDB(string(data))
		panicfIfErr(err, "unmarshal db config error. DB: %s.", fName)

		configs[i] = config
	}

	return configs
}

func initDB(dbConfigs []*sqlmap.DB) {
	for _, config := range dbConfigs {
		db, err := initDB2(config)
		panicfIfErr(err, "init db error. DB: %s.", config.Name)

		// In DB Scope.
		resultMaps := make(map[string]*SQLResutMap)
		actionCmds := make(map[string]*SQLActionCmd)

		sqlMaps := toSQLMaps(config.Locations)

		slice.Foreach(sqlMaps, func(sQLmap *sqlmap.SQLMap) {

			convResultMap(db, resultMaps, sQLmap.ResultMaps)

			slice.Foreach(sQLmap.Inserts, func(action sqlmap.ActionInsert) {
				convActionCmd(actionCmds, action.XMLName, action.Action)
			})
			slice.Foreach(sQLmap.Deletes, func(action sqlmap.ActionDelete) {
				convActionCmd(actionCmds, action.XMLName, action.Action)
			})
			slice.Foreach(sQLmap.Selects, func(action sqlmap.ActionSelect) {
				convActionCmd(actionCmds, action.XMLName, action.Action)
			})
			slice.Foreach(sQLmap.Updates, func(action sqlmap.ActionUpdate) {
				convActionCmd(actionCmds, action.XMLName, action.Action)
			})
			slice.Foreach(sQLmap.Operates, func(action sqlmap.ActionOperate) {
				convActionCmd(actionCmds, action.XMLName, action.Action)
			})
		})

		initStatements(db, actionCmds, resultMaps)
	}
}

func convActionCmd(actionCmds map[string]*SQLActionCmd,
	name xml.Name, action sqlmap.Action) {

	cmd := &SQLActionCmd{
		id:           action.Id,
		resultMap:    action.ResultMap,
		resultStruct: action.ResultStruct,
		sql:          action.SQL,
	}

	switch name.Local {
	case "insert":
		cmd.operation = INSERT
	case "select":
		cmd.operation = SELECT
	case "update":
		cmd.operation = UPDATE
	case "delete":
		cmd.operation = DELETE
	case "operate":
		cmd.operation = UNKOWN
	default:
		panic("light/db: Never happen.")
	}

	panicfIfTrue(actionCmds[cmd.id] != nil, "Duplicate statement ID: %s.", cmd.id)
	actionCmds[cmd.id] = cmd

	/*------------------ validate --------------*/
	// SQL type not match with xml tag.
	if cmd.operation != getOperation(cmd.sql) {
		log.Warnf("light/db: SQL cmd operation not match with sql. operation id: %s.", cmd.id)
	}

	panicfIfTrue(cmd.id == "", "Statement Id must not be empty.")
	if cmd.operation == SELECT {
		panicfIfTrue(cmd.resultMap == "" && cmd.resultStruct == "",
			"Statment %s must have resultMap or resultStruct.", cmd.id)
	}
}

func convResultMap(db *DB, resultMaps map[string]*SQLResutMap, rmaps []sqlmap.ResultMap) {

	for _, rmap := range rmaps {

		srps := make([]*SQLResultProperty, len(rmap.Properties))
		for i, p := range rmap.Properties {
			srp := &SQLResultProperty{
				column:   p.Column,
				property: p.Property,
				goType:   p.GoType,
				nilValue: p.NilValue,
			}
			srps[i] = srp
		}

		r := &SQLResutMap{
			id:         rmap.Id,
			structName: rmap.Struct,
			properties: srps,
		}

		panicfIfTrue(resultMaps[r.id] != nil, "Duplicate sqlResultMap ID: %s.", r.id)
		resultMaps[r.id] = r
	}
}

func initDB2(dbConfig *sqlmap.DB) (*DB, error) {
	c := Config{Name: dbConfig.Name}
	slice.Foreach(dbConfig.Props, func(p sqlmap.DBProp) {
		switch p.Name {
		case "driver":
			c.Driver = p.Value
		case "dsn":
			c.Dsn = p.Value
		case "maxIdleConns":
			c.MaxIdleConns, _ = strconv.Atoi(p.Value)
		case "maxOpenConns":
			c.MaxOpenConns, _ = strconv.Atoi(p.Value)
		default:
			log.Warnf("light/db: Unknown db config.[name:%s, value:%s]", p.Name, p.Value)
		}
	})

	db, err := sql.Open(c.Driver, c.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	lightDB := &DB{
		db:         db,
		config:     c,
		types:      make(map[string]interface{}),
		statements: make(map[string]*Statement),
		dynamicers: make(map[string]Dynamicer),
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if origin := dataSources[c.Name]; origin == nil {
		dataSources[c.Name] = lightDB
	} else {
		lightDB.types = origin.types
		lightDB.dynamicers = origin.dynamicers
		*origin = *lightDB
	}

	return lightDB, nil
}

func toSQLMaps(locations []sqlmap.Location) []*sqlmap.SQLMap {

	sqlMaps := make([]*sqlmap.SQLMap, len(locations))
	for i, location := range locations {
		data, err := conf.Read(DBLocation + location.Resource)
		panicfIfErr(err, "Read sqlmap config error. File: %s.", location.Resource)

		sqlMap, err := sqlmap.ReadSQLMap(string(data))
		panicfIfErr(err, "Unmarshal sqlmap config error: File: %s.", location.Resource)
		sqlMaps[i] = sqlMap
	}
	return sqlMaps
}
