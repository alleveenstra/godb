package godb

import (
	"github.com/ziutek/mymysql/mysql"
	"github.com/ziutek/mymysql/native"
	"fmt")

type ResultSet struct {
	Columns []Column
	Rows    []Row
}

type Column struct {
	Kind string
	Name string
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Column *Column
	Data interface{}
}

func NewResultSet(rows []mysql.Row, mysqlResult mysql.Result) ResultSet {
	var resultSet ResultSet
	if len(rows) == 0 {
		return resultSet
	}
	fields := mysqlResult.Fields()
	// set the columns
	resultSet.Columns = make([]Column, len(fields))
	for fieldIndex, field := range fields {
		resultSet.Columns[fieldIndex].Name = field.Name
		kind := "invalid"
		switch field.Type {
		case native.MYSQL_TYPE_DECIMAL:    // = 0x00
			kind = "float"
		case native.MYSQL_TYPE_TINY:       // = 0x01 // int8, uint8, bool
			kind = "int"
		case native.MYSQL_TYPE_SHORT:      // = 0x02 // int16, uint16
			kind = "int"
		case native.MYSQL_TYPE_LONG:       // = 0x03 // int32, uint32
			kind = "int"
		case native.MYSQL_TYPE_FLOAT:      // = 0x04 // float32
			kind = "float"
		case native.MYSQL_TYPE_DOUBLE:     // = 0x05 // float64
			kind = "float"
		case native.MYSQL_TYPE_NULL:       // = 0x06 // nil
			kind = "null"
		case native.MYSQL_TYPE_TIMESTAMP:  // = 0x07 // Timestamp
			kind = "timestamp"
		case native.MYSQL_TYPE_LONGLONG:   // = 0x08 // int64, uint64
			kind = "int"
		case native.MYSQL_TYPE_INT24:      // = 0x09
			kind = "int"
		case native.MYSQL_TYPE_DATE:       // = 0x0a // Date
			kind = "date"
		case native.MYSQL_TYPE_TIME:       // = 0x0b // Time
			kind = "time"
		case native.MYSQL_TYPE_DATETIME:   // = 0x0c // time.Time
			kind = "datetime"
		case native.MYSQL_TYPE_YEAR:       // = 0x0d
			kind = "int"
		case native.MYSQL_TYPE_NEWDATE:    // = 0x0e
			kind = "date"
		case native.MYSQL_TYPE_VARCHAR:    // = 0x0f
			kind = "string"
		case native.MYSQL_TYPE_BIT:        // = 0x10
			kind = "bool"
		case native.MYSQL_TYPE_NEWDECIMAL: // = 0xf6
			kind = "float"
		case native.MYSQL_TYPE_ENUM:       // = 0xf7
			kind = "string"
		case native.MYSQL_TYPE_SET:        // = 0xf8
			kind = "string"
		case native.MYSQL_TYPE_TINY_BLOB:  // = 0xf9
			kind = "bytes"
		case native.MYSQL_TYPE_MEDIUM_BLOB:// = 0xfa
			kind = "bytes"
		case native.MYSQL_TYPE_LONG_BLOB:  // = 0xfb
			kind = "bytes"
		case native.MYSQL_TYPE_BLOB:       // = 0xfc // Blob
			kind = "bytes"
		case native.MYSQL_TYPE_VAR_STRING: // = 0xfd // []byte
			kind = "string"
		case native.MYSQL_TYPE_STRING:     // = 0xfe // string
			kind = "string"
		case native.MYSQL_TYPE_GEOMETRY:   // = 0xff
			kind = "vector"
		}
		resultSet.Columns[fieldIndex].Kind = kind
	}
	// set the cells
	resultSet.Rows = make([]Row, len(rows))
	for rowIndex := range rows {
		resultSet.Rows[rowIndex].Cells = make([]Cell, len(fields))
		for columnIndex := range resultSet.Columns {
			resultSet.Rows[rowIndex].Cells[columnIndex].Column = &resultSet.Columns[columnIndex]
			resultSet.Rows[rowIndex].Cells[columnIndex].Data = rows[rowIndex][columnIndex]
		}
	}
	return resultSet
}

func (this *Cell) String() string {
	if this.Data == nil {
		return ""
	}
	switch this.Column.Kind {
	case "date":
		return this.Data.(mysql.Date).String()
	case "string":
		return string(this.Data.([]byte))
	case "float":
		return fmt.Sprintf("%f", this.Data)
	case "int":
		return fmt.Sprintf("%d", this.Data)
	case "bool":
		if this.Data.([]uint8)[0] == 1 {
			return "T"
		} else {
			return "F"
		}
	}
	return fmt.Sprintf("%v", this.Column.Kind)
}
