package document

import (
	docdriver "github.com/Gromosome/gorix/document-driver-manager"
)

type ScopedRepository[T any, ID comparable] struct {
	*Repository[T, ID]
}

func (r *Repository[T, ID]) WithExecutor(
	executor docdriver.Executor,
) *ScopedRepository[T, ID] {
	clone := *r
	clone.executor = executor

	return &ScopedRepository[T, ID]{
		Repository: &clone,
	}
}
