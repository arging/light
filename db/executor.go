// Copyright 2014 li. All rights reserved.

package db

import (
	"database/sql"
	"github.com/roverli/utils/errors"
)

// Executor declares all methods involved with executing statements
// and batches for an SQL Map.
// NOTE: This is the main CRUD api for users.
type Executor interface {

	// Executes a mapped SQL SELECT statement that returns data to populate
	// a single object instance.
	// The parameter object is generally used to supply the input
	// data for the WHERE clause parameter(s) of the SELECT statement.
	QueryOne(id string, param interface{}) (interface{}, errors.Error)

	// Executes a mapped SQL SELECT statement that returns data to populate
	// a number of result objects.
	// The parameter object is generally used to supply the input
	// data for the WHERE clause parameter(s) of the SELECT statement.
	QueryMany(id string, param interface{}) ([]interface{}, errors.Error)

	// Executes a mapped SQL INSERT statement.
	// Insert is a bit different from other update methods, as it
	// provides facilities for returning the primary key of the
	// newly inserted row (rather than the effected rows).  This
	// functionality is of course optional.
	// The parameter object is generally used to supply the input
	// data for the INSERT values.
	Insert(id string, param interface{}) (int64, errors.Error)

	// Executes a mapped SQL UPDATE statement.
	// Update can also be used for any other update statement type,
	// such as inserts and deletes. Update returns the number of
	// rows effected.
	// The parameter object is generally used to supply the input
	// data for the UPDATE values as well as the WHERE clause parameter(s).
	Update(id string, param interface{}) (int64, errors.Error)

	// Executes a mapped SQL DELETE statement.
	// Delete returns the number of rows effected.
	// The parameter object is generally used to supply the input
	// data for the WHERE clause parameter(s) of the DELETE statement.
	Delete(id string, param interface{}) (int64, errors.Error)

	// Executes any SQL statement.
	// The parameter object is generally used to supply the input
	// data for the WHERE clause or supply the input data.
	Exec(id string, param interface{}) (sql.Result, errors.Error)

	/** Inner methods. **/

	// Prepare a sql stament
	prepare(sql string) (*sql.Stmt, error)

	// Retrieve a sql statement definition.
	statement(id string) *Statement
}

// DB is the core struct for CRUD operations.
type DB struct {
	db         *sql.DB
	config     Config
	types      map[string]interface{}
	statements map[string]*Statement
	dynamicers map[string]Dynamicer
}

func (db *DB) Typer(name string, v interface{}) {
	panicfIfTrue(db.types[name] != nil, "duplicate typer. Name: %s.", name)
	db.types[name] = v
}

func (db *DB) Dynamicer(name string, dynamicer Dynamicer) {
	panicfIfTrue(db.dynamicers[name] != nil, "Duplicate dynamicer. Name: %s.", name)
	db.dynamicers[name] = dynamicer
}

func (db *DB) QueryOne(id string, param interface{}) (interface{}, errors.Error) {
	return queryOne(id, param, db)
}

func (db *DB) QueryMany(id string, param interface{}) ([]interface{}, errors.Error) {
	return queryMany(id, param, db)
}

func (db *DB) Insert(id string, param interface{}) (int64, errors.Error) {
	return insert(id, param, db)
}

func (db *DB) Update(id string, param interface{}) (int64, errors.Error) {
	return update(id, param, db)
}

func (db *DB) Delete(id string, param interface{}) (int64, errors.Error) {
	return delete(id, param, db)
}

func (db *DB) Exec(id string, param interface{}) (sql.Result, errors.Error) {
	return exec(id, param, db)
}

func (db *DB) prepare(sql string) (*sql.Stmt, error) {
	return db.db.Prepare(sql)
}

func (db *DB) statement(id string) *Statement {
	return db.statements[id]
}

// Transaction struct.
type Transaction struct {
	tx *sql.Tx
	db *DB
}

// Transaction callback interface.
type TxCallBack interface {

	// Return true indicates commit transaction. Otherwise, rollback.
	doTransaction(t Transaction) bool
}

// Execute a transaction with callback.
// If callBack method panic, tx will auto rollback
func (db *DB) DoTransaction(cb TxCallBack) (err errors.Error) {

	tx, e := db.db.Begin()
	if e != nil {
		err = errors.Wrap(e, "light/db: cann't begin a transaction.")
		return
	}

	defer func() {
		if v := recover(); v != nil {
			// Transaction done ?
			err = errors.Newf("%v, %v", v, tx.Rollback())
		}
	}()

	isCommit := cb.doTransaction(Transaction{tx: tx, db: db})
	if isCommit {
		e = tx.Commit()
	} else {
		e = tx.Rollback()
	}
	if e != nil {
		err = errors.Wrap(e, "light/db: handle trasaction fail.")
	}
	return
}

func (transaction *Transaction) QueryOne(id string, param interface{}) (interface{}, errors.Error) {
	return queryOne(id, param, transaction)
}

func (transaction *Transaction) QueryMany(id string, param interface{}) ([]interface{}, errors.Error) {
	return queryMany(id, param, transaction)
}

func (transaction *Transaction) Insert(id string, param interface{}) (int64, errors.Error) {
	return insert(id, param, transaction)
}

func (transaction *Transaction) Update(id string, param interface{}) (int64, errors.Error) {
	return update(id, param, transaction)
}

func (transaction *Transaction) Delete(id string, param interface{}) (int64, errors.Error) {
	return delete(id, param, transaction)
}

func (transaction *Transaction) Exec(id string, param interface{}) (sql.Result, errors.Error) {
	return exec(id, param, transaction)
}

func (transaction *Transaction) prepare(sql string) (*sql.Stmt, error) {
	return transaction.tx.Prepare(sql)
}

func (transaction *Transaction) statement(id string) *Statement {
	return transaction.db.statement(id)
}
