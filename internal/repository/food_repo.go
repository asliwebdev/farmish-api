package repository

import (
	"database/sql"
	"fmt"

	"farmish/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FoodRepository struct {
	DB *sql.DB
}

func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{DB: db}
}

func (r *FoodRepository) CreateFood(food *models.FoodWithoutTime) error {
	query := `
        INSERT INTO foods 
        (id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.DB.Exec(query, food.ID, food.FarmID, food.Name, pq.Array(food.SuitableFor), food.UnitOfMeasure, food.Quantity, food.MinThreshold)
	if err != nil {
		return fmt.Errorf("failed to create warehouse food: %v", err)
	}
	return nil
}

func (r *FoodRepository) GetAllFoods(farmID uuid.UUID) ([]models.Food, error) {
	query := `
        SELECT id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold, created_at, updated_at
        FROM foods
        WHERE farm_id = $1
    `
	rows, err := r.DB.Query(query, farmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get foods: %v", err)
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var food models.Food
		err := rows.Scan(
			&food.ID,
			&food.FarmID,
			&food.Name,
			pq.Array(&food.SuitableFor),
			&food.UnitOfMeasure,
			&food.Quantity,
			&food.MinThreshold,
			&food.CreatedAt,
			&food.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan food: %v", err)
		}
		foods = append(foods, food)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return foods, nil
}

func (r *FoodRepository) GetFoodByID(foodID uuid.UUID) (*models.Food, error) {
	query := `
	SELECT id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold, created_at, updated_at
	FROM foods
	WHERE id = $1
`
	row := r.DB.QueryRow(query, foodID)

	var food models.Food

	if err := row.Scan(&food.ID, &food.FarmID, &food.Name, pq.Array(&food.SuitableFor), &food.UnitOfMeasure, &food.Quantity, &food.MinThreshold, &food.CreatedAt, &food.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get food: %v", err)
	}

	return &food, nil
}

func (r *FoodRepository) UpdateFood(food *models.UpdateFoodReq) error {
	query := `
        UPDATE foods
        SET name = $1, suitable_for = $2, unit_of_measure = $3, quantity = $4, min_threshold = $5
        WHERE id = $6
    `
	_, err := r.DB.Exec(query, food.Name, pq.Array(food.SuitableFor), food.UnitOfMeasure, food.Quantity, food.MinThreshold, food.ID)
	if err != nil {
		return fmt.Errorf("failed to update food: %v", err)
	}
	return nil
}

func (r *FoodRepository) DeleteFood(foodID uuid.UUID) error {
	query := `DELETE FROM foods WHERE id = $1`
	_, err := r.DB.Exec(query, foodID)
	if err != nil {
		return fmt.Errorf("failed to delete food: %v", err)
	}
	return nil
}
