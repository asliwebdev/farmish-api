package repository

import (
	"database/sql"
	"errors"
	"farmish/internal/models"

	"github.com/google/uuid"
)

type FeedingRecordRepository struct {
	db *sql.DB
}

func NewFeedingRecordRepository(db *sql.DB) *FeedingRecordRepository {
	return &FeedingRecordRepository{db: db}
}

var ErrRecordNotFound = errors.New("feeding record not found")

func (r *FeedingRecordRepository) CreateFeedingRecord(record *models.FeedingRecordWithoutTime, newFoodQuantity float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	updateQuery := `
		UPDATE foods
		SET quantity = $1
		WHERE id = $2
	`
	_, err = tx.Exec(updateQuery, newFoodQuantity, record.FoodID)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO feeding_records (id, animal_id, food_id, quantity, fed_at, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(insertQuery, record.ID, record.AnimalID, record.FoodID, record.Quantity, record.FedAt, record.Notes)
	if err != nil {
		return err
	}

	return nil
}

func (r *FeedingRecordRepository) GetFeedingRecordByID(id uuid.UUID) (*models.FeedingRecordDetailed, error) {
	query := `
	SELECT 
	  fr.id AS feeding_record_id, 
	  fr.quantity, 
	  fr.fed_at, 
	  fr.notes, 
	  fr.created_at, 
	  a.id AS animal_id, 
	  a.name AS animal_name, 
	  a.type AS animal_type, 
	  a.weight AS animal_weight, 
	  a.health_status AS animal_health_status, 
	  f.id AS food_id, 
	  f.name AS food_name, 
	  f.suitable_for AS food_suitable_for, 
	  f.unit_of_measure AS food_unit_of_measure
	FROM feeding_records fr
	INNER JOIN animals a ON fr.animal_id = a.id
	INNER JOIN foods f ON fr.food_id = f.id
	WHERE fr.id = $1;
	`

	var detailedRecord models.FeedingRecordDetailed
	err := r.db.QueryRow(query, id).Scan(
		&detailedRecord.FeedingRecordID,
		&detailedRecord.Quantity,
		&detailedRecord.FedAt,
		&detailedRecord.Notes,
		&detailedRecord.CreatedAt,
		&detailedRecord.Animal.ID,
		&detailedRecord.Animal.Name,
		&detailedRecord.Animal.Type,
		&detailedRecord.Animal.Weight,
		&detailedRecord.Animal.HealthStatus,
		&detailedRecord.Food.ID,
		&detailedRecord.Food.Name,
		&detailedRecord.Food.SuitableFor,
		&detailedRecord.Food.UnitOfMeasure,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &detailedRecord, nil
}

func (r *FeedingRecordRepository) GetFeedingRecordsByAnimalID(animalID uuid.UUID) ([]models.FeedingRecordDetailed, error) {
	query := `
	SELECT 
	  fr.id AS feeding_record_id, 
	  fr.quantity, 
	  fr.fed_at, 
	  fr.notes, 
	  fr.created_at, 
	  a.id AS animal_id, 
	  a.name AS animal_name, 
	  a.type AS animal_type, 
	  a.weight AS animal_weight, 
	  a.health_status AS animal_health_status, 
	  f.id AS food_id, 
	  f.name AS food_name, 
	  f.suitable_for AS food_suitable_for, 
	  f.unit_of_measure AS food_unit_of_measure
	FROM feeding_records fr
	INNER JOIN animals a ON fr.animal_id = a.id
	INNER JOIN foods f ON fr.food_id = f.id
	WHERE fr.animal_id = $1;
	`

	rows, err := r.db.Query(query, animalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.FeedingRecordDetailed
	for rows.Next() {
		var detailedRecord models.FeedingRecordDetailed
		err := rows.Scan(
			&detailedRecord.FeedingRecordID,
			&detailedRecord.Quantity,
			&detailedRecord.FedAt,
			&detailedRecord.Notes,
			&detailedRecord.CreatedAt,
			&detailedRecord.Animal.ID,
			&detailedRecord.Animal.Name,
			&detailedRecord.Animal.Type,
			&detailedRecord.Animal.Weight,
			&detailedRecord.Animal.HealthStatus,
			&detailedRecord.Food.ID,
			&detailedRecord.Food.Name,
			&detailedRecord.Food.SuitableFor,
			&detailedRecord.Food.UnitOfMeasure,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, detailedRecord)
	}

	return records, nil
}

func (r *FeedingRecordRepository) UpdateFeedingRecord(record *models.FeedingRecordWithoutTime) error {
	query := `
		UPDATE feeding_records
		SET quantity = $1, fed_at = $2, notes = $3
		WHERE id = $4
	`
	result, err := r.db.Exec(query, record.Quantity, record.FedAt, record.Notes, record.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (r *FeedingRecordRepository) DeleteFeedingRecord(id uuid.UUID) error {
	query := `DELETE FROM feeding_records WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
