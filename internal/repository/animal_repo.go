package repository

import (
	"database/sql"
	"fmt"

	"farmish/internal/models"

	"github.com/google/uuid"
)

type AnimalRepository struct {
	DB *sql.DB
}

func NewAnimalRepository(db *sql.DB) *AnimalRepository {
	return &AnimalRepository{DB: db}
}

func (r *AnimalRepository) CreateAnimal(animal *models.AnimalWithoutTime) error {
	query := `
    INSERT INTO animals (id, farm_id, name, type, weight, health_status, date_of_birth)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
  `
	_, err := r.DB.Exec(query, animal.ID, animal.FarmID, animal.Name, animal.Type, animal.Weight,
		animal.HealthStatus, animal.DateOfBirth)
	if err != nil {
		return fmt.Errorf("failed to create animal: %v", err)
	}
	return nil
}

func (r *AnimalRepository) GetAnimalByID(id uuid.UUID) (*models.Animal, error) {
	query := `SELECT id, farm_id, name, type, weight, health_status, date_of_birth, last_fed, last_watered, created_at, updated_at FROM animals WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var animal models.Animal
	if err := row.Scan(&animal.ID, &animal.FarmID, &animal.Name, &animal.Type, &animal.Weight,
		&animal.HealthStatus, &animal.DateOfBirth, &animal.LastFed, &animal.LastWatered, &animal.CreatedAt, &animal.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get animal: %v", err)
	}

	return &animal, nil
}

func (r *AnimalRepository) GetAnimalsByFarmID(farmID uuid.UUID) ([]*models.Animal, error) {
	query := `SELECT id, farm_id, name, type, weight, health_status, date_of_birth, last_fed, last_watered, created_at, updated_at FROM animals WHERE farm_id = $1`
	rows, err := r.DB.Query(query, farmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get animals by farm ID: %v", err)
	}
	defer rows.Close()

	var animals []*models.Animal
	for rows.Next() {
		var animal models.Animal
		if err := rows.Scan(&animal.ID, &animal.FarmID, &animal.Name, &animal.Type, &animal.Weight,
			&animal.HealthStatus, &animal.DateOfBirth, &animal.LastFed, &animal.LastWatered, &animal.CreatedAt, &animal.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan animal: %v", err)
		}
		animals = append(animals, &animal)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return animals, nil
}

func (r *AnimalRepository) UpdateAnimal(animal *models.UpdateAnimalReq) error {
	query := `
    UPDATE animals
    SET name = $1, type = $2, weight = $3, health_status = $4, date_of_birth = $5, last_fed = $6, last_watered = $7
    WHERE id = $8
  `
	_, err := r.DB.Exec(query, animal.Name, animal.Type, animal.Weight, animal.HealthStatus, animal.DateOfBirth,
		animal.LastFed, animal.LastWatered, animal.ID)
	if err != nil {
		return fmt.Errorf("failed to update animal: %v", err)
	}
	return nil
}

func (r *AnimalRepository) DeleteAnimal(id uuid.UUID) error {
	query := `DELETE FROM animals WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete animal: %v", err)
	}
	return nil
}
