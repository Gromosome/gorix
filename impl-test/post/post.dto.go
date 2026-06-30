package post

type PostPathDto struct {
	ID string `param:"id" validate:"required"`
}

type CreatePostDto struct {
	Title     string   `json:"title" validate:"required,min=3,max=150"`
	Slug      string   `json:"slug" validate:"required,min=3,max=180"`
	Content   string   `json:"content" validate:"required,min=10"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

type UpdatePostDto struct {
	Title     string   `json:"title" validate:"omitempty,min=3,max=150"`
	Slug      string   `json:"slug" validate:"omitempty,min=3,max=180"`
	Content   string   `json:"content" validate:"omitempty,min=10"`
	Tags      []string `json:"tags"`
	Published *bool    `json:"published"`
}

type PostQueryDto struct {
	Published *bool `query:"published"`
	Limit     int   `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int   `query:"offset" validate:"omitempty,min=0"`
}
