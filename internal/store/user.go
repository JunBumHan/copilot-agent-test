package store

import (
	"github.com/JunBumHan/copilot-agent-test/internal/domain"
	"gorm.io/gorm"
)

// UserModel is the GORM model for user persistence
type UserModel struct {
	gorm.Model
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"size:100;unique;not null"`
}

// ToEntity converts the model to a domain entity
func (u *UserModel) ToEntity() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a model
func (u *UserModel) FromEntity(user *domain.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
}

// UserRepository implements domain.UserRepository using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository with the given database connection
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByID fetches a user by ID
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user UserModel
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user.ToEntity(), nil
}

// FindAll fetches all users
func (r *UserRepository) FindAll() ([]*domain.User, error) {
	var userModels []UserModel
	if err := r.db.Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(userModels))
	for i, model := range userModels {
		users[i] = model.ToEntity()
	}
	return users, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	model := UserModel{}
	model.FromEntity(user)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	*user = *model.ToEntity()
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	model := UserModel{}
	model.FromEntity(user)
	return r.db.Save(&model).Error
}

// Delete removes a user by ID
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&UserModel{}, id).Error
}