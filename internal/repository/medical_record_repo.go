package repository

import (
	"database/sql"
	"errors"
	"farmish/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type MedicalRecordRepository struct {
	db *sql.DB
}

func NewMedicalRecordRepository(db *sql.DB) *MedicalRecordRepository {
	return &MedicalRecordRepository{db: db}
}

var ErrMedicalRecordNotFound = errors.New("feeding record not found")

func (r *MedicalRecordRepository) CreateMedicalRecord(record *models.MedicalRecordWithoutTime, newMedicineQuantity float64) error {
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
	UPDATE medicines
	SET quantity = $1
	WHERE id = $2
`
	_, err = tx.Exec(updateQuery, newMedicineQuantity, record.MedicineID)
	if err != nil {
		return err
	}

	insertQuery := `
    INSERT INTO medical_records (id, animal_id, medicine_id, quantity, treatment_date, notes)
    VALUES ($1, $2, $3, $4, $5, $6)
  `
	_, err = tx.Exec(insertQuery, record.ID, record.AnimalID, record.MedicineID, record.Quantity, record.TreatmentDate, record.Notes)
	if err != nil {
		return err
	}
	return nil
}

func (r *MedicalRecordRepository) GetMedicalRecordByID(recordID uuid.UUID) (*models.MedicalRecordDetailed, error) {
	query := `
    SELECT
      mr.id AS medical_record_id, 
      mr.quantity, 
	  mr.treatment_date, 
	  mr.notes,
	  mr.created_at,
      a.id AS animal_id, 
	  a.name AS animal_name, 
	  a.type AS animal_type, 
	  a.weight AS animal_weight, 
	  a.health_status AS animal_health_status, 
	  m.id AS medicine_id,
	  m.name AS medicine_name,
	  m.suitable_for AS medicine_suitable_for
	  m.unit_of_measure AS medicine_unit_of_measure
    FROM medical_records mr
    INNER JOIN animals a ON mr.animal_id = a.id
    INNER JOIN medicines m ON mr.medicine_id = m.id
    WHERE mr.id = $1
  `
	row := r.db.QueryRow(query, recordID)

	var record models.MedicalRecordDetailed
	err := row.Scan(
		&record.ID,
		&record.Quantity,
		&record.TreatmentDate,
		&record.Notes,
		&record.CreatedAt,
		&record.Animal.ID,
		&record.Animal.Name,
		&record.Animal.Type,
		&record.Animal.Weight,
		&record.Animal.HealthStatus,
		&record.Medicine.ID,
		&record.Medicine.Name,
		&record.Medicine.SuitableFor,
		&record.Medicine.UnitOfMeasure,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get medical record by ID: %v", err)
	}
	return &record, nil
}

func (r *MedicalRecordRepository) GetMedicalRecordsByAnimalID(animalID uuid.UUID) ([]*models.MedicalRecordDetailed, error) {
	query := `
  SELECT
	mr.id AS medical_record_id, 
	mr.quantity, 
	mr.treatment_date, 
	mr.notes,
	mr.created_at,
	a.id AS animal_id, 
	a.name AS animal_name, 
	a.type AS animal_type, 
	a.weight AS animal_weight, 
	a.health_status AS animal_health_status, 
	m.id AS medicine_id,
	m.name AS medicine_name,
	m.suitable_for AS medicine_suitable_for
	m.unit_of_measure AS medicine_unit_of_measure
  FROM medical_records mr
  INNER JOIN animals a ON mr.animal_id = a.id
  INNER JOIN medicines m ON mr.medicine_id = m.id
  WHERE mr.animal_id = $1
  ORDER BY mr.treatment_date DESC
`
	rows, err := r.db.Query(query, animalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical records by animal ID: %v", err)
	}
	defer rows.Close()

	var records []*models.MedicalRecordDetailed
	for rows.Next() {
		var record models.MedicalRecordDetailed
		err := rows.Scan(
			&record.ID,
			&record.Quantity,
			&record.TreatmentDate,
			&record.Notes,
			&record.CreatedAt,
			&record.Animal.ID,
			&record.Animal.Name,
			&record.Animal.Type,
			&record.Animal.Weight,
			&record.Animal.HealthStatus,
			&record.Medicine.ID,
			&record.Medicine.Name,
			&record.Medicine.SuitableFor,
			&record.Medicine.UnitOfMeasure,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan medical record: %v", err)
		}
		records = append(records, &record)
	}
	return records, nil
}

func (r *MedicalRecordRepository) UpdateMedicalRecord(record *models.MedicalRecordWithoutTime) error {
	query := `
    UPDATE medical_records
    SET quantity = $1,
      treatment_date = $2, notes = $3
    WHERE id = $4
  `
	result, err := r.db.Exec(query, record.Quantity, record.TreatmentDate, record.Notes, record.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMedicalRecordNotFound
	}

	return nil
}

func (r *MedicalRecordRepository) DeleteMedicalRecord(recordID uuid.UUID) error {
	query := `DELETE FROM medical_records WHERE id = $1`
	result, err := r.db.Exec(query, recordID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMedicalRecordNotFound
	}

	return nil
}
