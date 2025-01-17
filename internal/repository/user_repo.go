package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"farmish/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

var ErrEmailAlreadyInUse = errors.New("email is already in use")

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (id, name, email, phone_number, password_hash)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.DB.Exec(query, user.ID, user.Name, user.Email, user.PhoneNumber, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrEmailAlreadyInUse
		}
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, name, email, phone_number, created_at FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %v", err)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	query := `SELECT id, name, email, phone_number, created_at FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, userID)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, phone_number, created_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.UpdateUser) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, phone_number = $3
    `

	updateValues := []interface{}{user.Name, user.Email, user.PhoneNumber}

	if user.Password != "" {
		query += ", password_hash = $4"
		updateValues = append(updateValues, user.Password)
	}
	query += " WHERE id = $" + strconv.Itoa(len(updateValues)+1)
	updateValues = append(updateValues, user.ID)

	_, err := r.DB.Exec(query, updateValues...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
