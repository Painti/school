package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func getTodosFromDB(db *sql.DB) (todos []Todo, err error) {
	query := `SELECT id, title, status FROM todos`
	stmt, pErr := db.Prepare(query)
	if pErr != nil {
		err = pErr
		return
	}
	rows, qErr := stmt.Query()
	if qErr != nil {
		err = qErr
		return
	}
	for rows.Next() {
		t := Todo{}
		sErr := rows.Scan(&t.ID, &t.Title, &t.Status)
		if sErr != nil {
			err = sErr
			return
		}
		todos = append(todos, t)
	}
	return
}

func getTodoByIDFromDB(db *sql.DB, id int) (t Todo, err error) {
	query := `SELECT id, title, status FROM todos WHERE id=$1`
	stmt, pErr := db.Prepare(query)
	if pErr != nil {
		err = pErr
		return
	}
	row := stmt.QueryRow(id)
	sErr := row.Scan(&t.ID, &t.Title, &t.Status)
	if sErr != nil {
		err = sErr
		return
	}
	return
}

func createTodo(db *sql.DB, todo Todo) (id int, err error) {
	query := `INSERT INTO todos(title, status) VALUES($1, $2) RETURNING id`
	row := db.QueryRow(query, todo.Title, todo.Status)
	sErr := row.Scan(&id)
	if sErr != nil {
		err = sErr
	}
	return
}

func deleteTodo(db *sql.DB, id int) (err error) {
	query := `DELETE FROM todos WHERE id=$1`
	stmt, pErr := db.Prepare(query)
	if pErr != nil {
		err = pErr
		return
	}
	result, eErr := stmt.Exec(id)
	if eErr != nil {
		err = eErr
	}

	row, rErr := result.RowsAffected()
	if rErr != nil {
		err = rErr
	}
	if row == 0 {
		err = errors.New(fmt.Sprintf("The todo ID:%d is not exists.", id))
	}
	return
}

func getTodosHandler(c *gin.Context, db *sql.DB) {
	todos, err := getTodosFromDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func getTodoByIDHandler(c *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := getTodoByIDFromDB(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func createTodoHandler(c *gin.Context, db *sql.DB) {
	var t Todo
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := createTodo(db, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	t.ID = id
	c.JSON(http.StatusOK, t)
}

func deleteTodoHandler(c *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteTodo(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func main() {
	url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Println("DB connect failed")
		return
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/api/todos", func(c *gin.Context) {
		getTodosHandler(c, db)
	})
	r.GET("/api/todos/:id", func(c *gin.Context) {
		getTodoByIDHandler(c, db)
	})
	r.POST("/api/todos", func(c *gin.Context) {
		createTodoHandler(c, db)
	})
	r.DELETE("/api/todos/:id", func(c *gin.Context) {
		deleteTodoHandler(c, db)
	})

	r.Run(":1234")
	fmt.Println("Server started")

}
