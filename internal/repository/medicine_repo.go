package repository

import (
	"database/sql"
	"farmish/internal/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MedicineRepository struct {
	DB *sql.DB
}

func NewMedicineRepository(db *sql.DB) *MedicineRepository {
	return &MedicineRepository{DB: db}
}

func (r *MedicineRepository) CreateMedicine(medicine *models.MedicineWithoutTime) error {
	query := `
    INSERT INTO medicines (id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
  `
	_, err := r.DB.Exec(query, medicine.ID, medicine.FarmID, medicine.Name, pq.Array(medicine.SuitableFor), medicine.UnitOfMeasure, medicine.Quantity, medicine.MinThreshold)
	if err != nil {
		return fmt.Errorf("failed to create medicine: %v", err)
	}
	return nil
}

func (r *MedicineRepository) GetAllMedicines(farmID uuid.UUID) ([]models.Medicine, error) {
	query := `SELECT id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold, created_at, updated_at FROM medicines WHERE farm_id = $1`
	rows, err := r.DB.Query(query, farmID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch medicines: %v", err)
	}
	defer rows.Close()

	var medicines []models.Medicine
	for rows.Next() {
		var medicine models.Medicine
		err := rows.Scan(&medicine.ID, &medicine.FarmID, &medicine.Name, pq.Array(&medicine.SuitableFor), &medicine.UnitOfMeasure, &medicine.Quantity, &medicine.MinThreshold, &medicine.CreatedAt, &medicine.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan medicine row: %v", err)
		}
		medicines = append(medicines, medicine)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return medicines, nil
}

func (r *MedicineRepository) GetMedicineByID(id uuid.UUID) (*models.Medicine, error) {
	query := `SELECT id, farm_id, name, suitable_for, unit_of_measure, quantity, min_threshold, created_at, updated_at FROM medicines WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var medicine models.Medicine
	err := row.Scan(&medicine.ID, &medicine.FarmID, &medicine.Name, pq.Array(&medicine.SuitableFor), &medicine.UnitOfMeasure, &medicine.Quantity, &medicine.MinThreshold, &medicine.CreatedAt, &medicine.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch medicine by id: %v", err)
	}
	return &medicine, nil
}

func (r *MedicineRepository) UpdateMedicine(medicine *models.MedicineWithoutTime) error {
	query := `
    UPDATE medicines
    SET name = $1, suitable_for = $2, unit_of_measure = $3, quantity = $4, min_threshold = $5
    WHERE id = $6
  `
	_, err := r.DB.Exec(query, medicine.Name, pq.Array(medicine.SuitableFor), medicine.UnitOfMeasure, medicine.Quantity, medicine.MinThreshold, medicine.ID)
	if err != nil {
		return fmt.Errorf("failed to update medicine: %v", err)
	}
	return nil
}

func (r *MedicineRepository) DeleteMedicine(id uuid.UUID) error {
	query := `DELETE FROM medicines WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete medicine: %v", err)
	}
	return nil
}
