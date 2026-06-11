module github.com/Gromosome/gorix/sql-drivers/sqlite3

go 1.26

require (
	github.com/Gromosome/gorix/sql-driver-manager v0.0.0
	github.com/mattn/go-sqlite3 v1.14.32
)

replace github.com/Gromosome/gorix/sql-driver-manager => ../../sql-driver-manager
