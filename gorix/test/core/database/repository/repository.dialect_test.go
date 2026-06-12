package repository

import (
	"testing"

	repository2 "github.com/Gromosome/gorix/gorix/core/database/repository"
)

func TestResolveDialectAliases(t *testing.T) {
	tests := map[string]string{
		"postgres":      "postgres",
		"pgx":           "postgres",
		"mysql":         "mysql",
		"sqlserver":     "mssql",
		"godror":        "oracle",
		"sqlite":        "sqlite-modern",
		"sqlite-modern": "sqlite-modern",
	}

	for input, want := range tests {
		dialect, err := repository2.ResolveDialect(input)
		if err != nil {
			t.Fatalf("ResolveDialect(%q) returned error: %v", input, err)
		}
		if dialect.Name() != want {
			t.Fatalf("ResolveDialect(%q) name = %q, want %q", input, dialect.Name(), want)
		}
	}
}

func TestDialectPlaceholdersAndQuoting(t *testing.T) {
	if (repository2.PostgresDialect{}).Placeholder(0) != "$1" {
		t.Fatal("postgres placeholder should normalize position")
	}
	if (repository2.MSSQLDialect{}).Placeholder(2) != "@p2" {
		t.Fatal("mssql placeholder mismatch")
	}
	if (repository2.OracleDialect{}).Placeholder(3) != ":3" {
		t.Fatal("oracle placeholder mismatch")
	}
	if (repository2.MySQLDialect{}).QuoteIdentifier("users.name") != "`users`.`name`" {
		t.Fatal("mysql qualified quote mismatch")
	}
	if (repository2.PostgresDialect{}).QuoteIdentifier(`user"name`) != `"user""name"` {
		t.Fatal("postgres quote escaping mismatch")
	}
	if (repository2.MSSQLDialect{}).QuoteIdentifier(`user]name`) != `[user]]name]` {
		t.Fatal("mssql quote escaping mismatch")
	}
}

func TestResolveDialectRejectsUnsupportedDriver(t *testing.T) {
	if _, err := repository2.ResolveDialect("unknown"); err == nil {
		t.Fatal("unsupported dialect should return error")
	}
}
