package boat

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Create() {
	url := "postgres://xqxxogxt:DtBnIRdBp82_ts8pjBBdL2WUIWJx5mXc@satao.db.elephantsql.com:5432/xqxxogxt"
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect fail, ", err.Error())
	}
	defer db.Close()
	fmt.Println("Connected")

	queryCreateTable := `
		CREATE TABLE IF NOT EXISTS todos(
			id SERIAL PRIMARY KEY,
			title TEXT,
			status TEXT
		);
	`

	_, err = db.Exec(queryCreateTable)
	if err != nil {
		log.Fatal("Cant create table, ", err.Error())
	}

	fmt.Println("Created")
}
