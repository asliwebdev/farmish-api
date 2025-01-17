package repository

import (
	"database/sql"
	"fmt"

	"farmish/internal/models"

	"github.com/google/uuid"
)

type FarmRepository struct {
	DB *sql.DB
}

func NewFarmRepository(db *sql.DB) *FarmRepository {
	return &FarmRepository{DB: db}
}

func (r *FarmRepository) CreateFarm(farm *models.Farm) error {
	query := `
        INSERT INTO farms (id, name, location, owner_id)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.DB.Exec(query, farm.ID, farm.Name, farm.Location, farm.OwnerID)
	if err != nil {
		return fmt.Errorf("failed to create farm: %v", err)
	}
	return nil
}

func (r *FarmRepository) GetFarmByID(farmID uuid.UUID) (*models.Farm, error) {
	query := `SELECT id, name, location, owner_id, created_at FROM farms WHERE id = $1`
	row := r.DB.QueryRow(query, farmID)

	var farm models.Farm
	if err := row.Scan(&farm.ID, &farm.Name, &farm.Location, &farm.OwnerID, &farm.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get farm: %v", err)
	}

	return &farm, nil
}

func (r *FarmRepository) GetAllFarms() ([]models.Farm, error) {
	query := `SELECT id, name, location, owner_id, created_at FROM farms`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve farms: %v", err)
	}
	defer rows.Close()

	var farms []models.Farm
	for rows.Next() {
		var farm models.Farm
		if err := rows.Scan(&farm.ID, &farm.Name, &farm.Location, &farm.OwnerID, &farm.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan farm: %v", err)
		}
		farms = append(farms, farm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return farms, nil
}

func (r *FarmRepository) UpdateFarm(farm *models.UpdateFarmRequest) error {
	query := `
        UPDATE farms
        SET name = $1, location = $2, owner_id = $3
        WHERE id = $4
    `
	_, err := r.DB.Exec(query, farm.Name, farm.Location, farm.OwnerID, farm.ID)
	if err != nil {
		return fmt.Errorf("failed to update farm: %v", err)
	}
	return nil
}

func (r *FarmRepository) DeleteFarm(farmID uuid.UUID) error {
	query := `DELETE FROM farms WHERE id = $1`
	_, err := r.DB.Exec(query, farmID)
	if err != nil {
		return fmt.Errorf("failed to delete farm: %v", err)
	}
	return nil
}
