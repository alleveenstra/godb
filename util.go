package godb

import (
	"reflect"
	"strings"
)

func sqlName(sv reflect.Value, columnName string) reflect.Value {
	st := sv.Type()
	for index := 0; index < sv.NumField(); index++ {
		sf := st.Field(index)
		tag := sf.Tag.Get("sql")
		// option 1: ignore the field
		if tag == "-" || sf.Anonymous {
			continue
		}
		// option 2: the tag matches the columnname
		if parseTag(tag) == columnName {
			return sv.Field(index)
		}
		// option 3: the field matches the columnname
		if tag == "" && sf.Name == columnName {
			return sv.Field(index)
		}
	}
	return reflect.Value{}
}

func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
