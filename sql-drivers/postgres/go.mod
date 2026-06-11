module github.com/Gromosome/gorix/sql-drivers/postgres

go 1.26

require (
	github.com/Gromosome/gorix/sql-driver-manager v0.0.0
	github.com/jackc/pgx/v5 v5.7.6
)

replace github.com/Gromosome/gorix/sql-driver-manager => ../../sql-driver-manager
