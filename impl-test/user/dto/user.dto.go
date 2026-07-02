package dto

type UserPathDto struct {
	ID int64 `param:"id" validate:"required,min=1"`
}

type CreateUserDto struct {
	Name string `json:"name" validate:"required,min=3,max=100"`

	Email string `json:"email" validate:"required,email"`
}

type UpdateUserDto struct {
	Name string `json:"name" validate:"required,min=3,max=100"`

	Active bool `json:"active"`
}

type UserQueryDto struct {
	Active *bool `query:"active"`

	Limit int `query:"limit" validate:"min=1,max=100"`

	Offset int `query:"offset" validate:"min=0"`
}
