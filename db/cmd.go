// Copyright 2014 li. All rights reserved.

package db

import (
	"database/sql"
	"github.com/roverli/light/log"
	"github.com/roverli/utils/errors"
)

// Error Codes.
const (
	NoStatementErr     = iota + 1 // No match statement error
	PreSQLErr                     // Prepare SQL error
	QueryErr                      // PrepareStatement.Query() error
	ExecErr                       // PrepareStatement.Exec() error
	GetColumnsErr                 // Rows.Columns() error
	RowScanErr                    // Rows.Scan() error
	RowsError                     // Rows.Err() error
	PanicErr                      // Panic unkown error
	OperateNotMatchErr            // SQL type is not match function call.
	MultiResultErr                // Return multi results when call queryOne
	UnkownErr                     // Unkown reason error.
)

func selectRaw(id string, pvalue interface{}, executor Executor) (data []interface{}, e errors.Error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("light/db: Panic error. Err: %v", err)
			data = nil
			e = errors.NewfByCode(PanicErr, "Panic error. Error: %v", err)
		}
	}()

	st := executor.statement(id)
	if st == nil {
		return nil, errors.NewfByCode(NoStatementErr, "light/db: No such statement: %s.", id)
	}

	sql, params := st.ProcParam(pvalue)
	preStatement, err := executor.prepare(sql)
	log.Debugf("light/db: Exec statement: %s. SQL:%s. Params: %v.", id, sql, params)

	if err != nil {
		return nil, errors.WrapfByCode(PreSQLErr, err, "light/db: Prepare SQL error. SQL: %s. Statement: %s.", sql, id)
	}

	defer preStatement.Close()
	rows, err := preStatement.Query(params...)

	if err != nil {
		return nil, errors.WrapfByCode(QueryErr, err, "light/db: Query Err. SQL: %s. Params: %v. Statement: %s.", sql, params, id)
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.WrapfByCode(GetColumnsErr, err, "light/db: Get columns error. Statment: %s.", id)
	}

	result := make([]interface{}, 0)

	for rows.Next() {
		sourceHolder, destHolder := st.HoldValue(columns)
		err := rows.Scan(destHolder...)
		if err != nil {
			return nil, errors.WrapfByCode(RowScanErr, err, "light/db: Rows scan error. Statment: %s.", id)
		}
		sv := st.ToValue(sourceHolder)
		result = append(result, sv)
	}

	log.Debugf("light/db: Exec statement: %s. Columns: %v.", id, columns[0])

	err = rows.Err()
	if err != nil {
		return nil, errors.WrapfByCode(RowsError, err, "light/db: Rows error. Statment: %s.", id)
	}
	return result, nil
}

func rawExec(op Operation, id string, pValue interface{}, executor Executor) (data sql.Result, e errors.Error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("light/db: Panic error. Err: %v", err)
			data = nil
			e = errors.NewfByCode(PanicErr, "Panic error. Error: %v", err)
		}
	}()

	statement := executor.statement(id)
	if statement == nil {
		return nil, errors.NewfByCode(NoStatementErr, "light/db: No such statement: %s.", id)
	}

	if op != statement.op && op != UNKOWN {
		return nil, errors.NewfByCode(OperateNotMatchErr, "light/db: Operation not match. Statement: %s.", id)
	}

	sql, params := statement.ProcParam(pValue)
	preStatement, err := executor.prepare(sql)
	log.Debugf("light/db: Exec statement: %s. SQL:%s. Params: %v.", id, sql, params)

	if err != nil {
		return nil, errors.WrapfByCode(PreSQLErr, e, "light/db: Prepare SQL error. SQL: %s. Statement: %s.", sql, id)
	}

	defer preStatement.Close()

	result, err := preStatement.Exec(params...)
	if err != nil {
		return nil, errors.WrapfByCode(ExecErr, e, "light/db: Exec Err. SQL: %s. Params: %v. Statement: %s.", sql, params, id)
	}

	return result, nil
}

func queryOne(id string, param interface{}, executor Executor) (interface{}, errors.Error) {
	results, err := selectRaw(id, param, executor)
	if err != nil {
		return nil, err
	}

	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		return results[0], nil
	default:
		return nil, errors.NewfByCode(MultiResultErr, "light/db: Too many results for queryOne method. Statement: %s.", id)
	}
}

func queryMany(id string, param interface{}, executor Executor) ([]interface{}, errors.Error) {
	return selectRaw(id, param, executor)
}

func insert(id string, param interface{}, executor Executor) (int64, errors.Error) {
	result, err1 := rawExec(INSERT, id, param, executor)
	if err1 != nil {
		return 0, err1
	}

	insertId, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, errors.WrapfByCode(UnkownErr, err2, "light/db: Unkown Error. Statement: %s.", id)
	}
	return insertId, nil
}

func execWithAffectedRows(op Operation, id string, param interface{}, executor Executor) (int64, errors.Error) {
	result, err1 := rawExec(op, id, param, executor)
	if err1 != nil {
		return 0, err1
	}

	affectedRows, err2 := result.RowsAffected()
	if err2 != nil {
		return 0, errors.WrapfByCode(UnkownErr, err2, "light/db: Unkown Error. Statement: %s.", id)
	}
	return affectedRows, nil

}

func update(id string, param interface{}, executor Executor) (int64, errors.Error) {
	return execWithAffectedRows(UPDATE, id, param, executor)
}

func delete(id string, param interface{}, executor Executor) (int64, errors.Error) {
	return execWithAffectedRows(DELETE, id, param, executor)
}

func exec(id string, param interface{}, executor Executor) (sql.Result, errors.Error) {
	return rawExec(UNKOWN, id, param, executor)
}
