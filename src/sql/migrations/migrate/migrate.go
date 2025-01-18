package main

import (
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/casts"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/migrations"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

var DB_HOST = "127.0.0.1"
var DB_PORT = 5432
var DB_USER = "admin"
var DB_PASS = "492631bcd77b2ca81f0fecb6d9bc803b"
var DB_DATABASE = "verifier"
var DB_SSL = false

func ReadEnv() {
	DB_HOST = utils.ParseEnvString("DB_HOST", DB_HOST)
	DB_PORT = utils.ParseEnvInt("DB_PORT", DB_PORT)
	DB_USER = utils.ParseEnvString("DB_USER", DB_USER)
	DB_PASS = utils.ParseEnvString("DB_PASS", DB_PASS)
	DB_DATABASE = utils.ParseEnvString("DB_DATABASE", DB_DATABASE)
	DB_SSL = casts.IntToBool(utils.ParseEnvInt("DB_SSL", casts.BoolToInt(DB_SSL)))
}

func main() {
	db := sql.ConnectToPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_DATABASE, DB_SSL)
	migrations.M01Initial(db)
	migrations.M02AddMetadata(db)
	migrations.M03AddCounts(db)
	migrations.M04CreateErrorTable(db)
}
