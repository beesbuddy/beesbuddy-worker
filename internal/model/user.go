package model

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/beesbuddy/beesbuddy-worker/internal"
	"github.com/petaki/support-go/forms"
)

type User struct {
	ID            int       `gorm:"primaryKey;autoIncrement:true;column:user_id"`
	Username      string    `gorm:"unique;column:username;type:varchar(32);not null"`
	Email         string    `gorm:"unique;" validate:"required,email,min=6,max=32"`
	Password      string    `gorm:"column:password;type:varchar(128);not null" validate:"required,min=6"`
	RememberToken string    `gorm:"column:remember_token;type:varchar(128);default:null"`
	IsEnabled     bool      `gorm:"column:is_enabled;type:boolean;default:null"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;default:null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp;default:null"`
}

func NewUser() *User {
	return &User{
		IsEnabled: true,
	}
}

func UserCreateRules(form *forms.Form) {
	form.Required("username", "email", "password", "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 5)
}

func UserUpdateRules(form *forms.Form) {
	form.Required("username", "email", "is_enabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 5)
}

func (u *User) Fill(form *forms.Form) (*User, error) {
	u.Username = form.Data["username"].(string)
	u.Email = form.Data["email"].(string)
	u.IsEnabled = form.Data["is_enabled"].(bool)

	hash, err := internal.HashPassword(form.Data["password"].(string))

	if err != nil {
		return nil, err
	}

	u.Password = hash

	return u, nil
}

func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(u.Email)), size)
}

func (u *User) RememberCookie() []byte {
	return []byte(fmt.Sprintf("%v|%s", u.ID, u.RememberToken))
}
