package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New initializes the database connection pool.
func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

// Models holds all models.
type Models struct {
	User User
}

// User represents a user in the system.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"user_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll retrieves all users from the database.
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at FROM users ORDER BY last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			log.Println("Error scanning:", err)
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail retrieves a user by their email.
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at FROM users WHERE email = $1`
	var user User
	row := db.QueryRowContext(ctx, query, email)

	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no user found with that email")
		}
		log.Println("Error scanning:", err)
		return nil, err
	}

	log.Println("Inside models.go file line 103")
	return &user, nil
}

// GetOne retrieves a user by their ID.
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at FROM users WHERE id = ?`
	var user User
	row := db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no user found with that ID")
		}
		log.Println("Error scanning:", err)
		return nil, err
	}

	return &user, nil
}

// Update modifies an existing user's details.
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE users SET email = ?, first_name = ?, last_name = ?, user_active = ?, updated_at = ? WHERE id = ?`
	if _, err := db.ExecContext(ctx, stmt, u.Email, u.FirstName, u.LastName, u.Active, time.Now(), u.ID); err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}

// Delete removes a user from the database.
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM users WHERE id = ?`
	if _, err := db.ExecContext(ctx, stmt, u.ID); err != nil {
		log.Println("Error deleting user:", err)
		return err
	}

	return nil
}

// Insert adds a new user to the database.
func (u *User) Insert(user User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var newID int64
	stmt := `INSERT INTO users (email, first_name, last_name, password, user_active, created_at) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now()).Scan(&newID)

	if err != nil {
		log.Println("Error inserting new user:", err)
		return 0, err
	}

	return newID, nil
}

// ResetPassword updates a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET password = ? WHERE id = ?`
	if _, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID); err != nil {
		log.Println("Error resetting password:", err)
		return err
	}

	return nil
}

// PasswordMatches checks if the provided plain text password matches the hashed password.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil // Passwords do not match
	}

	return true, nil // Passwords match or other error occurred
}
