package repository

import (
	"reflect"
	"testing"

	repository2 "github.com/Gromosome/gorix/gorix/core/database/repository"
)

func TestQueryBuilderBuildSelect(t *testing.T) {
	query, args, err := repository2.NewQueryBuilder(repository2.PostgresDialect{}, "public.users").
		Select("id", "email").
		Where("email = ?", "a@example.com").
		Where("age > ?", 18).
		GroupBy(`"id"`).
		OrderBy(`"email" ASC`).
		Limit(10).
		Offset(20).
		BuildSelect()
	if err != nil {
		t.Fatalf("BuildSelect returned error: %v", err)
	}

	wantQuery := `SELECT "id", "email" FROM "public"."users" WHERE (email = $1) AND (age > $2) GROUP BY "id" ORDER BY "email" ASC LIMIT 10 OFFSET 20`
	if query != wantQuery {
		t.Fatalf("unexpected query:\n%s", query)
	}
	if !reflect.DeepEqual(args, []any{"a@example.com", 18}) {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestQueryBuilderRejectsMissingDialectOrTable(t *testing.T) {
	if _, _, err := repository2.NewQueryBuilder(nil, "users").BuildSelect(); err == nil {
		t.Fatal("nil dialect should return error")
	}
	if _, _, err := repository2.NewQueryBuilder(repository2.PostgresDialect{}, " ").BuildSelect(); err == nil {
		t.Fatal("empty table should return error")
	}
}
