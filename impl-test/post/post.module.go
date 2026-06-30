package post

import "github.com/Gromosome/gorix/gorix"

type PostModule struct{}

func NewPostModule() *PostModule {
	return &PostModule{}
}

func (m *PostModule) BasePath() gorix.BasePath {
	return "/post"
}

func (m *PostModule) Providers() []any {
	return []any{
		NewPostRepository,
		NewPostService,
	}
}

func (m *PostModule) Controllers() []any {
	return []any{
		NewPostController,
	}
}
