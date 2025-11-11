package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"pay-your-dues/internal/domain/entities"
	"pay-your-dues/internal/domain/interfaces"
	"pay-your-dues/internal/models"
)

// userRepositoryGORM implements the UserRepository interface using GORM
type userRepositoryGORM struct {
	db *gorm.DB
}

// NewUserRepositoryGORM creates a new user repository with GORM
func NewUserRepositoryGORM(db *gorm.DB) interfaces.UserRepository {
	return &userRepositoryGORM{
		db: db,
	}
}

func (r *userRepositoryGORM) Create(ctx context.Context, user *entities.User) error {
	gormUser := r.entityToGORM(user)
	if err := r.db.WithContext(ctx).Create(gormUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	// Update the entity with the created ID if it was auto-generated
	user.ID = gormUser.ID
	user.CreatedAt = gormUser.CreatedAt
	user.UpdatedAt = gormUser.UpdatedAt
	return nil
}

func (r *userRepositoryGORM) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var gormUser models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&gormUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return r.gormToEntity(&gormUser), nil
}

func (r *userRepositoryGORM) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var gormUser models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entities.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return r.gormToEntity(&gormUser), nil
}

func (r *userRepositoryGORM) Update(ctx context.Context, user *entities.User) error {
	gormUser := r.entityToGORM(user)
	if err := r.db.WithContext(ctx).Save(gormUser).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	// Update the entity with the updated timestamp
	user.UpdatedAt = gormUser.UpdatedAt
	return nil
}

func (r *userRepositoryGORM) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return entities.ErrUserNotFound
	}
	return nil
}

func (r *userRepositoryGORM) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	return count > 0, nil
}


func (r *userRepositoryGORM) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user exists by phone: %w", err)
	}
	return count > 0, nil
}

// entityToGORM converts a domain entity to GORM model
func (r *userRepositoryGORM) entityToGORM(user *entities.User) *models.User {
	return &models.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// gormToEntity converts a GORM model to domain entity
func (r *userRepositoryGORM) gormToEntity(gormUser *models.User) *entities.User {
	return &entities.User{
		ID:           gormUser.ID,
		Email:        gormUser.Email,
		PasswordHash: gormUser.PasswordHash,
		FirstName:    gormUser.FirstName,
		LastName:     gormUser.LastName,
		Phone:        gormUser.Phone,
		CreatedAt:    gormUser.CreatedAt,
		UpdatedAt:    gormUser.UpdatedAt,
	}
}
