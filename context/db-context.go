package context

import "database/sql"

func Connect() *sql.DB {
	db, err := sql.Open("sqlserver", "server=localhost;database=GoLang;trusted_connection=yes")

	if err != nil {
		panic(err)
	}
	return db
}
