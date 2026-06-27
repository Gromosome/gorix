package entity

import (
	"time"
)

type User struct {
	ID int64 `json:"id" xml:"id" db:"id" repository:"primaryKey,autoIncrement"`

	Name string `json:"name" xml:"name" db:"name"`

	Email string `json:"email" xml:"email" db:"email"`

	Active bool `json:"active" xml:"active" db:"active"`

	CreatedAt time.Time `json:"createdAt" xml:"createdAt" db:"created_at" repository:"readOnly"`

	UpdatedAt time.Time `json:"updatedAt" xml:"updatedAt" db:"updated_at" repository:"readOnly"`
}

func (User) TableName() string {
	return "users"
}
