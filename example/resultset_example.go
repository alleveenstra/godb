package example

var config struct {
	SqlHostname string
	SqlUser     string
	SqlPassword string
	SqlDatabase string
}

func resultset() {

	godb.Database = autorc.New("tcp", "", config.SqlHostname, config.SqlUser, config.SqlPassword, config.SqlDatabase)
	godb.Database.Register("SET NAMES utf8")

	stmt2, prepareErr2 := godb.Database.Prepare("SELECT * FROM `person` WHERE `id` = LAST_INSERT_ID()")
	rows, res, _ := stmt2.Exec()
	result := godb.NewResultSet(rows, res)
}
