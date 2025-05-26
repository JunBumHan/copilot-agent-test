package domain

import "time"

// User represents the user entity in our domain
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id uint) (*User, error)
	FindAll() ([]*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}