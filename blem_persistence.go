package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	// https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/05.3.html
	// https://siongui.github.io/2016/01/09/go-sqlite-example-basic-usage/
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func InitTables(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS TableMap(
		TableID BIGINT PRIMARY KEY,
		Schema TEXT,
		Table TEXT
	);

    CREATE TABLE IF NOT EXISTS (

    );
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func insertTableMap(t TypeTableName, DB *sql.DB) {
	insert_TableMap := `
    INSERT OR REPLACE INTO TableMap (
        TableID, Schema, Table
    ) VALUES (?,?,?)
    `

	stmt, err := DB.Prepare(insert_TableMap)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

}
