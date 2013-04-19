package example

import "log"

type Person struct {
	Name string
	Age  int
}

func load() {
	result, selectErr := godb.SqlAll("SELECT * FROM `person`")
	fatal(selectErr)
	data := make([]Person, len(result.Rows))
	unmarshalErr := godb.Unmarshal(data, result)
	fatal(unmarshalErr)
	// ...
}

func save() {
	var person Person{"Johnny", 30}

sql := godb.Marshal(person, "person", -1)
godb.Sql(sql)
}


func fatal(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}
