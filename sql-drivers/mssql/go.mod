module github.com/Gromosome/gorix/sql-drivers/mssql

go 1.26

require (
	github.com/Gromosome/gorix/sql-driver-manager v0.0.0
	github.com/microsoft/go-mssqldb v1.9.3
)

replace github.com/Gromosome/gorix/sql-driver-manager => ../../sql-driver-manager
