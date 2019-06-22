package boat

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Insert() {
	url := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect fail, ", err.Error())
	}
	defer db.Close()
	fmt.Println("Connected")

	title := "Home Work"
	status := "Inactive"
	query := `INSERT INTO todos(title, status) VALUES($1, $2) RETURNING id`

	var id int
	row := db.QueryRow(query, title, status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("Can't scan id", id)
	}
	fmt.Println("ID: ", id)

}
