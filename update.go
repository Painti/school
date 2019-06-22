package boat

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Update() {
	url := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("Connect fail, ", err.Error())
	}
	defer db.Close()
	fmt.Println("Connected")

	query := `UPDATE todos SET status=$2 WHERE id=$1`

	var id int
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("prepare error", id)
	}

	if _, err := stmt.Exec(1, "active"); err != nil {
		log.Fatal("exec error, ", err.Error())
	}
	fmt.Println("update success")

}
