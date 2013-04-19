package godb

import (
	"reflect"
	"fmt"
	"strings"
	"errors"
)

func Marshal(data interface{}, tableName string, id int) (string, []interface{}, error) {
	var sql string
	var values []interface{}
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Struct {
		columns, cells, err := marhsalRow(dataValue)
		if err != nil {
			return sql, values, err
		}
		if id >= 0 {
			set := make([]string, len(columns))
			values = make([]interface{}, len(columns) + 1)
			for index, column := range columns {
				set[index] = fmt.Sprintf(" %s = ? ", column)
				values[index] = cells[index]
			}
			values[len(values) - 1] = fmt.Sprintf("%d", id)
			sql = fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(set, ", "))
		} else {
			set := make([]string, len(columns))
			values = make([]interface{}, len(columns))
			for index, _ := range columns {
				set[index] = fmt.Sprintf("?")
				values[index] = cells[index]
			}
			sql = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", tableName, strings.Join(columns, ", "), strings.Join(set, ", "))
		}
		return sql, values, nil
	}
	return sql, values, errors.New("marshal only accepts structs")
}

func marhsalRow(myValue reflect.Value) ([]string, []string, error) {
	var columns []string
	var cells []string
	if myValue.Kind() == reflect.Struct {
		myType := myValue.Type()
		for i := 0; i < myType.NumField(); i++ {
			field := myType.Field(i)
			switch field.Type.Kind() {
			case reflect.String:
				columns = append(columns, fmt.Sprintf("`%s`", field.Name))
				cells = append(cells, fmt.Sprintf("%s", myValue.Field(i).String()))
			case reflect.Bool:
				columns = append(columns, fmt.Sprintf("`%s`", field.Name))
				cells = append(cells, fmt.Sprintf("%d", myValue.Field(i).Bool()))
			case reflect.Int32, reflect.Int64:
				columns = append(columns, fmt.Sprintf("`%s`", field.Name))
				cells = append(cells, fmt.Sprintf("%d", myValue.Field(i).Int()))
			case reflect.Float32, reflect.Float64:
				columns = append(columns, fmt.Sprintf("`%s`", field.Name))
				cells = append(cells, fmt.Sprintf("%f", myValue.Field(i).Float()))
			}
		}
		return columns, cells, nil
	}
	return columns, cells, errors.New("marhsalRow expects a struct")
}
