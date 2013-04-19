package godb

import (
	"github.com/ziutek/mymysql/mysql"
	"github.com/ziutek/mymysql/autorc"
	"fmt"
	"regexp"
	"errors"
)

var Database *autorc.Conn

func ShowTables() ([]string, error) {
	var tables []string
	stmt, prepareErr := Database.Prepare("SHOW TABLES")
	if prepareErr != nil {
		return tables, prepareErr
	}
	rows, _, execErr := stmt.Exec()
	if execErr != nil {
		return tables, execErr
	}
	tables = make([]string, len(rows))
	for tableIndex := range tables {
		tables[tableIndex] = rows[tableIndex].Str(0)
	}
	return tables, nil
}

type SqlColumn struct {
	Field string `sql:"Field"`
	Type string `sql:"Type"`
}

func ShowColumns(db string, fieldLike string) ([]SqlColumn, error) {
	var data []SqlColumn

	if !validDatabaseName(db) {
		return data, errors.New(fmt.Sprintf(`"%s" is not a valid database name`, db))
	}

	stmt, prepareErr := Database.Prepare(fmt.Sprintf("SHOW COLUMNS FROM `%s` WHERE `Field` LIKE ?", db)) // sql injection
	if (prepareErr != nil) {
		return data, prepareErr
	}

	rows, res, execErr := stmt.Exec(fieldLike)
	if (execErr != nil) {
		return data, execErr
	}

	result := NewResultSet(rows, res)

	data = make([]SqlColumn, len(result.Rows))

	errUnm := Unmarshal(data, result)
	if (errUnm != nil) {
		return data, errUnm
	}

	return data, nil
}

func SqlAll(sql string) (ResultSet, error) {
	var result ResultSet
	stmt, prepareErr := Database.Prepare(sql)
	if prepareErr != nil {
		return result, prepareErr
	}
	rows, res, execErr := stmt.Exec()
	if execErr != nil {
		return result, execErr
	}
	result = NewResultSet(rows, res)
	return result, nil
}

func SqlOne(sql string) (mysql.Row, error) {
	var result mysql.Row
	stmt, prepareErr := Database.Prepare(sql)
	if prepareErr != nil {
		return result, prepareErr
	}
	rows, _, execErr := stmt.Exec()
	if execErr != nil {
		return result, execErr
	}
	return rows[0], nil
}

func Execute(sql string, parameters ...interface{}) (err error) {
	stmt, prepareErr := Database.Prepare(sql)
	if prepareErr != nil {
		return prepareErr
	}
	_, _, execErr := stmt.Exec(parameters...)
	if execErr != nil {
		return execErr
	}
	return nil
}

var dbRegex, _ = regexp.Compile(`^[a-zA-Z0-9_-]+$`)

func validDatabaseName(db string) bool {
	return dbRegex.MatchString(db)
}
