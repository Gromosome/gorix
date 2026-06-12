package repository

import (
	"testing"

	repository2 "github.com/Gromosome/gorix/gorix/core/database/repository"
)

type metadataEntity struct {
	ID       int    `db:"id" repository:"primaryKey,autoIncrement"`
	FullName string `db:"full_name"`
	Secret   string `db:"-"`
}

func (metadataEntity) TableName() string {
	return "people"
}

func TestMetadataOfBuildsEntityMetadata(t *testing.T) {
	metadata, err := repository2.MetadataOf[metadataEntity]()
	if err != nil {
		t.Fatalf("MetadataOf returned error: %v", err)
	}
	if metadata.TableName != "people" {
		t.Fatalf("unexpected table name: %s", metadata.TableName)
	}
	if metadata.PrimaryKey == nil || metadata.PrimaryKey.ColumnName != "id" {
		t.Fatalf("unexpected primary key: %#v", metadata.PrimaryKey)
	}
	if len(metadata.Fields) != 2 {
		t.Fatalf("unexpected field count: %d", len(metadata.Fields))
	}
}

func TestMetadataOfRejectsInvalidEntities(t *testing.T) {
	if _, err := repository2.MetadataOf[int](); err == nil {
		t.Fatal("non-struct entity should be rejected")
	}

	type missingPrimaryKey struct {
		ID int
	}
	if _, err := repository2.MetadataOf[missingPrimaryKey](); err == nil {
		t.Fatal("missing primary key should be rejected")
	}
}
