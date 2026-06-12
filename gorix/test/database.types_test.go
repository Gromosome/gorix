package test

import (
	"testing"

	"github.com/Gromosome/gorix/gorix"
	"github.com/Gromosome/gorix/gorix/core/database/repository"
)

type aliasEntity struct {
	ID int `repository:"primaryKey"`
}

func TestDatabaseTypeAliasesExposeConstructors(t *testing.T) {
	manager := &gorix.DBManager{}

	repo, err := gorix.NewSQLRepository[aliasEntity, int](manager)
	if err != nil {
		t.Fatalf("NewSQLRepository returned error: %v", err)
	}

	if repo.ConnectionName() != "default" {
		t.Fatalf("unexpected repository connection name: %s", repo.ConnectionName())
	}

	query, args, err := gorix.NewQueryBuilder(repository.NewSQLiteDialect("sqlite3"), "users").
		Select("id").
		Where("id = ?", 10).
		BuildSelect()
	if err != nil {
		t.Fatalf("BuildSelect returned error: %v", err)
	}
	if query != `SELECT "id" FROM "users" WHERE (id = ?)` {
		t.Fatalf("unexpected query: %s", query)
	}
	if len(args) != 1 || args[0] != 10 {
		t.Fatalf("unexpected args: %#v", args)
	}
}
