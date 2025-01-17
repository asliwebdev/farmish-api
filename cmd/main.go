package main

import (
	"farmish/internal/handlers"
	"farmish/internal/repository"
	"farmish/internal/services"
	"farmish/pkg/config"
	"log"
)

func main() {
	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userService := services.NewUserService(repository.NewUserRepository(db))

	h := handlers.NewHandler(userService)

	r := handlers.Run(h)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
