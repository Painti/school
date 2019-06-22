package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int
	Title  string
	Status string
}

func getTodos(c *gin.Context) {
	url := "postgres://xqxxogxt:DtBnIRdBp82_ts8pjBBdL2WUIWJx5mXc@satao.db.elephantsql.com:5432/xqxxogxt"
	db, err := sql.Open("postgres", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	fmt.Println("DB Connected")

	query := `SELECT id, title, status FROM todos`
	stmt, err := db.Prepare(query)

	rows, _ := stmt.Query()

	tt := []Todo{}

	for rows.Next() {
		t := Todo{}
		err = rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tt = append(tt, t)
	}

	c.JSON(http.StatusOK, tt)
}

func main() {
	r := gin.Default()

	// r.GET("/students", getStudentHandler)
	// r.POST("/students", postStudentHandler)

	r.GET("/api/todos", getTodos)

	r.Run(":1234")
}
