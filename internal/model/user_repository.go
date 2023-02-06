package model

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func (r UserRepository) Create(user *User) error {
	result := r.orm.Create(&user)

	if result.Error != nil {
		return fmt.Errorf("unable to create a user")
	}

	return nil
}

func (r UserRepository) Authenticate(usernameOrEmail, password string) (*User, error) {
	usernameOrEmail = strings.ToLower(usernameOrEmail)

	user := &User{}
	result := r.orm.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsEnabled {
		return nil, ErrInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	return user, nil
}

func (r UserRepository) UpdateRememberToken(user *User, token string) error {
	tx := r.orm.Model(user).Update("remember_token", token)

	if !tx.Statement.Changed() {
		return ErrUnableUpdateRecord
	}

	return nil
}

func (r UserRepository) Find(ID int) (*User, error) {
	user := &User{}

	result := r.orm.Where("id = ?", ID).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNoRecord
	}

	return user, nil
}
