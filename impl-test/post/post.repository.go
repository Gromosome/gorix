package post

import (
	"fmt"
	"time"

	"github.com/Gromosome/gorix/gorix"
	postdoc "github.com/Gromosome/gorix/impl-test/post/document"
)

type Post = postdoc.Post
type PostAudit = postdoc.PostAudit

type PostRepository struct {
	documents *gorix.DocumentManager

	posts  *gorix.DocumentRepository[postdoc.Post, string]
	audits *gorix.DocumentRepository[postdoc.PostAudit, string]
}

func NewPostRepository(
	documents *gorix.DocumentManager,
) *PostRepository {
	postRepo, err := gorix.NewDocumentRepository[postdoc.Post, string](
		documents,
		"mongo",
	)
	if err != nil {
		panic(err)
	}

	auditRepo, err := gorix.NewDocumentRepository[postdoc.PostAudit, string](
		documents,
		"mongo",
	)
	if err != nil {
		panic(err)
	}

	return &PostRepository{
		documents: documents,
		posts:     postRepo,
		audits:    auditRepo,
	}
}

func (r *PostRepository) EnsureSchema(
	ctx *gorix.Context,
) error {
	if err := r.posts.EnsureSchema(ctx); err != nil {
		return err
	}

	if err := r.audits.EnsureSchema(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) FindByID(
	ctx *gorix.Context,
	id string,
) (*Post, error) {
	post, err := r.posts.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("post %s not found: %w", id, err)
	}

	return post, nil
}

func (r *PostRepository) Find(
	ctx *gorix.Context,
	query PostQueryDto,
) ([]postdoc.Post, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}

	q := gorix.NewDocumentQuery().
		Limit(int64(limit)).
		Offset(int64(query.Offset)).
		SortDesc("createdAt")
	if query.Published != nil {
		q.Eq("published", *query.Published)
	}

	return r.posts.Find(ctx, q)
}

func (r *PostRepository) Create(
	ctx *gorix.Context,
	post *postdoc.Post,
) error {
	if post == nil {
		return fmt.Errorf("post repository: post cannot be nil")
	}

	if err := r.EnsureSchema(ctx); err != nil {
		return err
	}

	if err := r.posts.Insert(ctx, post); err != nil {
		return err
	}

	audit := &postdoc.PostAudit{
		PostID:    post.ID,
		Action:    "POST_CREATED",
		CreatedAt: time.Now().UTC(),
	}

	return r.audits.Insert(ctx, audit)
}

func (r *PostRepository) CreateWithTransaction(
	ctx *gorix.Context,
	post *postdoc.Post,
) error {
	if post == nil {
		return fmt.Errorf("post repository: post cannot be nil")
	}

	if err := r.EnsureSchema(ctx); err != nil {
		return err
	}

	return gorix.WithDocumentTransaction(
		ctx,
		r.documents,
		"default",
		nil,
		func(
			ctx *gorix.Context,
			tx gorix.DocumentTx,
		) error {
			txPostRepo := r.posts.WithExecutor(tx)
			txAuditRepo := r.audits.WithExecutor(tx)

			if err := txPostRepo.Insert(ctx, post); err != nil {
				return err
			}
			audit := &postdoc.PostAudit{
				PostID:    post.ID,
				Action:    "POST_CREATED_TX",
				CreatedAt: time.Now().UTC(),
			}

			return txAuditRepo.Insert(ctx, audit)
		},
	)
}

func (r *PostRepository) Update(
	ctx *gorix.Context,
	post *postdoc.Post,
) error {
	if post == nil {
		return fmt.Errorf("post repository: post cannot be nil")
	}

	if err := r.posts.Update(ctx, post); err != nil {
		return err
	}

	audit := &postdoc.PostAudit{
		PostID:    post.ID,
		Action:    "POST_UPDATED",
		CreatedAt: time.Now().UTC(),
	}

	return r.audits.Insert(ctx, audit)
}

func (r *PostRepository) DeleteByID(
	ctx *gorix.Context,
	id string,
) error {
	return r.posts.DeleteByID(ctx, id)
}

func (r *PostRepository) Count(
	ctx *gorix.Context,
	query PostQueryDto,
) (int64, error) {
	q := gorix.NewDocumentQuery()

	if query.Published != nil {
		q.Eq("published", *query.Published)
	}
	return r.posts.Count(ctx, q)
}
