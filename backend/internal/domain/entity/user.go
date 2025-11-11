package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           string
	IconUrl        string
	CreatedAt      time.Time
	HashedPassword string
}

func NewUser(name string, iconUrl string, hashedPassword string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if iconUrl == "" {
		return nil, errors.New("iconUrl is required")
	}
	if hashedPassword == "" {
		return nil, errors.New("hashedPassword is required")
	}

	return &User{
		ID:             uuid.New(),
		Name:           name,
		IconUrl:        iconUrl,
		CreatedAt:      time.Now(),
		HashedPassword: hashedPassword,
	}, nil
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	u.Name = name
	return nil
}

func (u *User) ChangeIconUrl(iconUrl string) error {
	if iconUrl == "" {
		return errors.New("iconUrl is required")
	}
	u.IconUrl = iconUrl
	return nil
}

func (u *User) ChangeHashedPassword(hashedPassword string) error {
	if hashedPassword == "" {
		return errors.New("hashedPassword is required")
	}
	u.HashedPassword = hashedPassword
	return nil
}
