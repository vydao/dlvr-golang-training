package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	cors "github.com/rs/cors/wrapper/gin"
)

type TodoItem struct {
	ID        int64      `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Status    string     `json:"status" db:"status"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type createToDoItemRequest struct {
	Title string `json:"title" binding:"required"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=golang dbname=todo_app password=secret sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	runDBMigration("file://db/migrations", "postgresql://golang:secret@localhost:5432/todo_app?sslmode=disable")

	router := gin.New()
	router.Use(cors.Default())
	groupV1 := router.Group("/api/v1")

	groupV1.Handle(http.MethodPost, "/todos", func(c *gin.Context) {
		var request createToDoItemRequest
		var id int64

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.PrepareNamed(`INSERT INTO todo_items (title,status) VALUES (:title,:status) RETURNING id`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = stmt.Get(&id, TodoItem{Title: request.Title, Status: "Doing"})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": id})
	})

	groupV1.Handle(http.MethodGet, "/todos", func(c *gin.Context) {
		items := []TodoItem{}
		err := db.Select(&items, `SELECT * FROM todo_items`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": items})
	})

	groupV1.Handle(http.MethodDelete, "/todos/:id", func(c *gin.Context) {
		// id, err := strconv.Atoi(c.Param("id"))
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		res, err := db.Exec(`DELETE FROM todo_items WHERE id = $1`, c.Param("id"))
		if err == nil {
			rowsAffected, err := res.RowsAffected()
			if err == nil {
				if rowsAffected == 1 {
					c.JSON(http.StatusOK, gin.H{})
				}
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})

	router.Run(":8080")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("Could not create new migrate instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Could not run db migrate:", err)
	}

	log.Println("DB migrated successfully!")
}
