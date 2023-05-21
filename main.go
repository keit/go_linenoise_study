package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	linenoise "github.com/GeertJohan/go.linenoise"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "keitakashima"
	password = "password" // even you don't need it, you must pass.
	dbname   = "cmr_core_development"
)

func main() {
	db := initDB()
	// close database
	defer db.Close()

	linenoise.SetMultiline(true)

	for {
		str, err := linenoise.Line("gsql> ")
		if err != nil {
			if err == linenoise.KillSignalError {
				quit()
			}
			fmt.Printf("Unexpected error: %s\n", err)
			quit()
		}

		err = linenoise.AddHistory(str)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		if str == "quit" {
			quit()
		}
		doQuery(db, str)

	}

}

func doQuery(db *sql.DB, query string) {
	fmt.Printf("Query: %s\n", query)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error Query: %s\n", err)
	} else {
		defer rows.Close()
		columnTypes, err := rows.ColumnTypes()
		CheckError(err)
		fmt.Println("Column Name, Type, Length, Nullable")
		for _, v := range columnTypes {
			// do something
			fmt.Printf("%s %s %s %s\n", v.Name(), v.DatabaseTypeName(), fmtColLength(v), fmtColNullable(v))
		}
	}
}

func fmtColLength(col *sql.ColumnType) string {
	length, ok := col.Length()
	if ok {
		return strconv.FormatInt(length, 10)
	} else {
		return "-"
	}
}

func fmtColNullable(col *sql.ColumnType) string {
	nullable, ok := col.Nullable()
	if ok {
		return strconv.FormatBool(nullable)
	} else {
		return "-"
	}
}

func initDB() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func quit() {
	fmt.Println("")
	os.Exit(0)
}
