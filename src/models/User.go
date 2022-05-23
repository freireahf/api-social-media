package models

import (
	"api/src/secure"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

//User represent User in database
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

//Prepare execute methods validate and format in received user
func (user *User) Prepare(stage string) error {
	if err := user.validate(stage); err != nil {
		return err
	}

	if err := user.format(stage); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(stage string) error {
	if user.Name == "" {
		return errors.New("Name cannot be blank")
	}

	if user.Nick == "" {
		return errors.New("Nick cannot be blank")
	}

	if user.Email == "" {
		return errors.New("Email cannot be blank")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("Email format is invalid")
	}

	if stage == "create" && user.Password == "" {
		return errors.New("Password cannot be blank")
	}

	return nil
}

func (user *User) format(stage string) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)

	if stage == "create" {
		passwordWithHash, err := secure.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(passwordWithHash)
	}

	return nil
}
