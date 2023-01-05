package models

import (
	"crypto/md5"
	"fmt"

	"github.com/petaki/support-go/forms"
)

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	RememberToken string    `json:"-"`
	IsEnabled     bool      `json:"isEnabled"`
	CreatedAt     Timestamp `json:"createdAt"`
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
