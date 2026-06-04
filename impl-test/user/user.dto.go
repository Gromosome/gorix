package user

type CreateUserBodyDto struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"regex=^\\+?[0-9]{7,15}$"`
	Age   int    `json:"age" validate:"min=18,max=100"`
}
type UserSearchQueryDto struct {
	Page  int    `query:"page" validate:"min=1"`
	Limit int    `query:"limit" validate:"min=1,max=100"`
	Sort  string `query:"sort" validate:"oneof=asc desc"`
}
type UserPathDto struct {
	ID string `param:"id" validate:"required"`
}
