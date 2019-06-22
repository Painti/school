package boat

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Query() {
	url := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect fail, ", err.Error())
	}
	defer db.Close()
	fmt.Println("Connected")

	query := `SELECT id, title, status FROM todos`
	stmt, err := db.Prepare(query)

	rows, _ := stmt.Query()

	tt := []Todo{}

	for rows.Next() {
		t := Todo{}
		err = rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			log.Fatal("error, ", err.Error())
		}
		tt = append(tt, t)
	}

	fmt.Println("one row", tt)
}
