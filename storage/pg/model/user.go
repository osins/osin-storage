package model

import (
	"time"

	"github.com/google/uuid"
)

// User define
type User struct {
	Id uuid.UUID `gorm:"primaryKey;->;<-:create;"`

	Username string

	Password []byte `json:"-"`

	Salt []byte `json:"-"`

	EMail string `gorm:"column:email";`

	Mobile string

	// Date created
	CreatedAt time.Time
}

func (s *User) GetId() string {
	if s == nil || s.Id == uuid.Nil {
		return ""
	}

	return s.Id.String()
}

func (s *User) GetUsername() string {
	return s.Username
}

func (s *User) GetPassword() []byte {
	return s.Password
}

func (s *User) GetSalt() []byte {
	return s.Salt
}

func (s *User) GetMobile() string {
	return s.Mobile
}

func (s *User) GetEmail() string {
	return s.EMail
}
