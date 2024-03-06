package model

import (
	"time"

	"github.com/google/uuid"
)

// User define
type User struct {
	Id uuid.UUID `gorm:"primaryKey;->;<-:create;type:char(36);"`

	Username string `gorm:"type:varchar(64);unique_index"`

	Password string `gorm:"type:varchar(64);unique_index"`

	Salt string `gorm:"type:varchar(36);unique_index"`

	EMail string `gorm:"column:email;type:varchar(256);unique_index"`

	Mobile string `gorm:"type:char(11);unique_index"`

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

func (s *User) GetPassword() string {
	return s.Password
}

func (s *User) GetSalt() string {
	return s.Salt
}

func (s *User) GetMobile() string {
	return s.Mobile
}

func (s *User) GetEmail() string {
	return s.EMail
}
