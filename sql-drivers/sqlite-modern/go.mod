module github.com/Gromosome/gorix/sql-drivers/sqlite-modern

go 1.26

require (
	github.com/Gromosome/gorix/sql-driver-manager v0.0.0
	modernc.org/sqlite v1.39.1
)

replace github.com/Gromosome/gorix/sql-driver-manager => ../../sql-driver-manager
