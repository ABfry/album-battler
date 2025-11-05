package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
	Name string
	IconUrl string
	CreatedAt time.Time
	HashedPassword string
}

func NewUser(name string, iconUrl string, hashedPassword string) *User {
	return &User{
		ID: uuid.New(),
		Name: name,
		IconUrl: iconUrl,
		CreatedAt: time.Now(),
		HashedPassword: hashedPassword,
	}
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	u.Name = name
	return nil
}

// todo