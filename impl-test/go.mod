module github.com/Gromosome/gorix/impl-test

go 1.25

require github.com/Gromosome/gorix/gorix v1.0.0
require github.com/Gromosome/gorix/sql-drivers/postgres v1.0.0
require github.com/Gromosome/gorix/document-drivers/mongo v0.0.0
require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/Gromosome/gorix/gorix => ../gorix
replace github.com/Gromosome/gorix/document-drivers/mongo => ../document-drivers/mongo