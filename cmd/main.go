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

	animalRepo := repository.NewAnimalRepository(db)
	foodRepo := repository.NewFoodRepository(db)

	userService := services.NewUserService(repository.NewUserRepository(db))
	farmService := services.NewFarmService(repository.NewFarmRepository(db))
	animalService := services.NewAnimalService(animalRepo)
	foodService := services.NewFoodService(foodRepo)
	medicineService := services.NewMedicineService(repository.NewMedicineRepository(db))
	feedingRecordService := services.NewFeedingRecordService(repository.NewFeedingRecordRepository(db), animalRepo, foodRepo)

	h := handlers.NewHandler(userService, farmService, animalService, foodService, medicineService, feedingRecordService)

	r := handlers.Run(h)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
