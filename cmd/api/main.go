package main

import (
	"log"
	"work/todo/internal/database"
	"work/todo/internal/handler"
	"work/todo/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {

	dbcfg := database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secretpassword",
		DBName:   "tododb",
		SSLMode:  "disable",
	}

	db, err := database.NewPostgresConnecction(dbcfg)

	if err != nil {
		log.Fatalf("unable to connect db: %v", err)
	}
	defer db.Close()

	todoRepo := repository.NewTodoRepository(db)

	todoHandler := handler.NewTodoHandler(todoRepo)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/todos", todoHandler.Create)
		v1.GET("/todos", todoHandler.GetAll)
		v1.GET("/todos/:id", todoHandler.GetByID)
		v1.PUT("/todos/:id", todoHandler.Update)
		v1.DELETE("/todos/:id", todoHandler.Delete)
	}
	log.Println("Starting on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
