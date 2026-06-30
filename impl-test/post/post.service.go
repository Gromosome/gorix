package post

import (
	"time"

	"github.com/Gromosome/gorix/gorix"
	postdoc "github.com/Gromosome/gorix/impl-test/post/document"
)

type PostService struct {
	postRepository *PostRepository
}

func NewPostService(
	postRepository *PostRepository,
) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) EnsureSchema(
	ctx *gorix.Context,
) (any, error) {
	return map[string]any{
		"success": true,
		"message": "post document schema ensured",
	}, s.postRepository.EnsureSchema(ctx)
}

func (s *PostService) FindByID(
	ctx *gorix.Context,
	id string,
) (*Post, error) {
	return s.postRepository.FindByID(ctx, id)
}

func (s *PostService) Find(
	ctx *gorix.Context,
	query PostQueryDto,
) ([]postdoc.Post, error) {
	return s.postRepository.Find(ctx, query)
}

func (s *PostService) Count(
	ctx *gorix.Context,
	query PostQueryDto,
) (map[string]any, error) {
	count, err := s.postRepository.Count(ctx, query)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"count": count,
	}, nil
}

func (s *PostService) Create(
	ctx *gorix.Context,
	request CreatePostDto,
) (*Post, error) {
	now := time.Now().UTC()

	post := &postdoc.Post{
		Title:     request.Title,
		Slug:      request.Slug,
		Content:   request.Content,
		Tags:      request.Tags,
		Published: request.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.postRepository.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) CreateTx(
	ctx *gorix.Context,
	request CreatePostDto,
) (*Post, error) {
	now := time.Now().UTC()

	post := &postdoc.Post{
		Title:     request.Title,
		Slug:      request.Slug,
		Content:   request.Content,
		Tags:      request.Tags,
		Published: request.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.postRepository.CreateWithTransaction(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Update(
	ctx *gorix.Context,
	id string,
	request UpdatePostDto,
) (*Post, error) {
	post, err := s.postRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Title != "" {
		post.Title = request.Title
	}

	if request.Slug != "" {
		post.Slug = request.Slug
	}

	if request.Content != "" {
		post.Content = request.Content
	}

	if request.Tags != nil {
		post.Tags = request.Tags
	}

	if request.Published != nil {
		post.Published = *request.Published
	}

	post.UpdatedAt = time.Now().UTC()

	if err := s.postRepository.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeleteByID(
	ctx *gorix.Context,
	id string,
) (any, error) {
	return nil, s.postRepository.DeleteByID(ctx, id)
}
