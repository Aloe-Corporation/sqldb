package test

import "os"

var (
	WorkDir            = os.Getenv("GOPATH") + "/src/github.com/Aloe-Corporation/sqldb"
	PathScriptPostgres = WorkDir + "/test/schema_postgres.sql"
	PathScriptMySQL    = WorkDir + "/test/schema_mysql.sql"
)
