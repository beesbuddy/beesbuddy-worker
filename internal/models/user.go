package models

import (
	"crypto/md5"
	"fmt"
	"time"

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

type UserRepository interface {
	Create(*User) error
	Find(int) (*User, error)
	FindAll() ([]*User, error)
	Update(*User, *User) error
	UpdateRememberToken(*User, string) error
	Authenticate(string, string) (*User, error)
	AuthenticateByRememberToken(int, string) (*User, error)
	Delete(*User) error
}

func NewUser() *User {
	return &User{
		IsEnabled: true,
	}
}

func UserCreateRules(form *forms.Form) {
	form.Required("username", "email", "password", "isEnabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
}

func UserUpdateRules(form *forms.Form) {
	form.Required("username", "email", "isEnabled")
	form.MatchesPattern("username", forms.UsernameRegexp)
	form.Min("username", 3)
	form.Max("username", 20)
	form.MatchesPattern("email", forms.EmailRegexp)
	form.Min("password", 8)
}

func (u *User) Fill(form *forms.Form) *User {
	u.Username = form.Data["username"].(string)
	u.Email = form.Data["email"].(string)
	u.IsEnabled = form.Data["is_enabled"].(bool)
	u.Password = form.Data["password"].(string)

	return u
}

func (u *User) Gravatar(size int) string {
	return fmt.Sprintf("https://gravatar.com/avatar/%x?s=%d", md5.Sum([]byte(u.Email)), size)
}

func (u *User) RememberCookie() []byte {
	return []byte(fmt.Sprintf("%v|%s", u.ID, u.RememberToken))
}
