package godb

import (
	"reflect"
	"errors"
)

func Unmarshal(data interface{}, resultSet ResultSet) error {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Slice {
		if dataValue.Len() != len(resultSet.Rows) {
			return errors.New("unmarshal: data and resultSet slices must have the same length")
		}
		for index := range resultSet.Rows {
			err := unmarshalRow(dataValue.Index(index), resultSet.Rows[index])
			if err != nil {
				return err
			}
		}
		return nil
	}
	if dataValue.Kind() == reflect.Ptr {
		if len(resultSet.Rows) == 0 {
			return errors.New("unmarshal requires at least one row")
		}
		return unmarshalRow(dataValue.Elem(), resultSet.Rows[0])
	}
	return errors.New("unmarshal only accepts slices and pointers to structs")
}

func unmarshalRow(myValue reflect.Value, row Row) error {
	if myValue.Kind() == reflect.Struct {
		for _, cell := range row.Cells {
			fieldValue := sqlName(myValue, cell.Column.Name)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				switch fieldValue.Type().Kind() {
				case reflect.String:
					fieldValue.SetString(string(cell.Data.([]byte)))
				case reflect.Bool:
					fieldValue.SetBool(cell.Data.(bool))
				case reflect.Int32:
					fieldValue.SetInt(int64(cell.Data.(int32)))
				case reflect.Int64:
					fieldValue.SetInt(cell.Data.(int64))
				case reflect.Float32:
					fieldValue.SetFloat(float64(cell.Data.(float32)))
				case reflect.Float64:
					fieldValue.SetFloat(cell.Data.(float64))
				}
			}
		}
		return nil
	}
	return errors.New("unmarshalRow expects a struct")
}
