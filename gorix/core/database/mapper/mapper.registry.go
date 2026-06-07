package mapper

import (
	"fmt"
	"strings"
	"sync"
)

type StatementRegistry struct {
	mutex      sync.RWMutex
	statements map[string]string
}

func NewStatementRegistry() *StatementRegistry {
	return &StatementRegistry{
		statements: make(map[string]string),
	}
}

func (r *StatementRegistry) Register(
	name string,
	query string,
) error {
	name = strings.TrimSpace(name)
	query = strings.TrimSpace(query)

	if name == "" {
		return fmt.Errorf(
			"gorix mapper: statement name cannot be empty",
		)
	}

	if query == "" {
		return fmt.Errorf(
			"gorix mapper: statement %q cannot be empty",
			name,
		)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.statements[name]; exists {
		return fmt.Errorf(
			"gorix mapper: statement %q is already registered",
			name,
		)
	}

	r.statements[name] = query
	return nil
}

func (r *StatementRegistry) MustRegister(
	name string,
	query string,
) {
	if err := r.Register(name, query); err != nil {
		panic(err)
	}
}

func (r *StatementRegistry) Get(
	name string,
) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query, exists := r.statements[name]
	if !exists {
		return "", fmt.Errorf(
			"gorix mapper: statement %q is not registered",
			name,
		)
	}

	return query, nil
}

func (r *StatementRegistry) Has(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.statements[name]
	return exists
}
