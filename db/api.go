// Copyright 2014 li. All rights reserved.

package db

// Interface shields behaviors on different databases.
type Dialect interface {
}

// MySql Dialect
type Mysql interface {
}

// Handle sql dynamic
type Dynamicer func(param interface{}) (string, []interface{})

// Config is the representation of DataSource settings.
type Config struct {
	Name         string  // DataSource name
	Driver       string  // The DB Driver, like "mysql" ...
	Dsn          string  // The DataSource Url, like "userx:passwordx@tcp(www.mysql1.com:3306)/db1"
	MaxIdleConns int     // Same to db.SetMaxIdleConns
	MaxOpenConns int     // Same to db.SetMaxOpenConns
	Dial         Dialect // Used to shield the differ of underlying database.
}

// Get a DB instance.
func Get(name string) *DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := dataSources[name]
	if db == nil {
		db = &DB{
			types:      make(map[string]interface{}),
			dynamicers: make(map[string]Dynamicer),
		}
		dataSources[name] = db
	}
	return db
}
