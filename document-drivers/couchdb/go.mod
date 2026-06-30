module github.com/Gromosome/gorix/document-drivers/couchdb

go 1.25

require (
	github.com/Gromosome/gorix/document-driver-manager v0.0.0
	github.com/go-kivik/couchdb/v4 v4.0.0-20230828195858-5c44e9a72d49
	github.com/go-kivik/kivik/v4 v4.5.2
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
)

replace github.com/Gromosome/gorix/document-driver-manager => ../../document-driver-manager
