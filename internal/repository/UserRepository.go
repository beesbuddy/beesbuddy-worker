package repository

import (
	"fmt"

	"github.com/beesbuddy/beesbuddy-worker/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	orm *gorm.DB
}

func NewUserRepository(orm *gorm.DB) *UserRepository {
	return &UserRepository{
		orm,
	}
}

func (r UserRepository) Create(user *model.User) error {
	result := r.orm.Create(&user)

	if result.Error != nil {
		return fmt.Errorf("unable to create a user")
	}

	return nil
}
