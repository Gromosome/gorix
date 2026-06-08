package entity

import "time"

type User struct {
	ID int64 `json:"id" db:"id" repository:"primaryKey,autoIncrement"`

	Name string `json:"name" db:"name"`

	Email string `json:"email" db:"email"`

	Active bool `json:"active" db:"active"`

	CreatedAt time.Time `json:"createdAt" db:"created_at" repository:"readOnly"`

	UpdatedAt time.Time `json:"updatedAt" db:"updated_at" repository:"readOnly"`
}

func (User) TableName() string {
	return "users"
}
