package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID `gorm:"primaryKey;->;<-:create;"`

	ClientId string

	Username string

	Password string

	EMail string

	Mobile string

	// Date created
	CreatedAt time.Time
}
