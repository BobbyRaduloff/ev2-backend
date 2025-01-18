package env

import (
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/casts"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

var DB_INCOMING_QUEUE_TIMEOUT = 30 * time.Second
var DB_INCOMING_QUEUE_MAX_SIZE = 5000
var DB_HOST = "127.0.0.1"
var DB_PORT = 5432
var DB_USER = "admin"
var DB_PASS = "admin"
var DB_DATABASE = "verifier"
var DB_SSL = false
var LOGGER_SYNC_TIMEOUT = 60 * time.Second
var TASK_CHUNK_SIZE = 100
var API_PORT = 6969
var API_MAX_FILE_SIZE_BYTES = 1073741824 // 1 GiB
var API_UPLOADS_DIR = "./uploads"
var MQ_MAX_NEW_TASKS = 1
var MQ_URL = "amqp://admin:admin@localhost:5672/"
var FORGOTTEN_TIMEOUT = 5 * time.Minute
var TASK_CHANNEL_QUEUE_SIZE = 10000

func ReadEnv() {
	DB_INCOMING_QUEUE_TIMEOUT = utils.ParseEnvDuration("DB_QUEUE_TIMEOUT_SECONDS", DB_INCOMING_QUEUE_TIMEOUT)
	DB_INCOMING_QUEUE_MAX_SIZE = utils.ParseEnvInt("DB_QUEUE_MAX_SIZE", DB_INCOMING_QUEUE_MAX_SIZE)
	DB_HOST = utils.ParseEnvString("DB_HOST", DB_HOST)
	DB_PORT = utils.ParseEnvInt("DB_PORT", DB_PORT)
	DB_USER = utils.ParseEnvString("DB_USER", DB_USER)
	DB_PASS = utils.ParseEnvString("DB_PASS", DB_PASS)
	DB_DATABASE = utils.ParseEnvString("DB_DATABASE", DB_DATABASE)
	DB_SSL = casts.IntToBool(utils.ParseEnvInt("DB_SSL", casts.BoolToInt(DB_SSL)))
	LOGGER_SYNC_TIMEOUT = utils.ParseEnvDuration("LOGGER_SYNC_TIMEOUT_SECONDS", LOGGER_SYNC_TIMEOUT)
	TASK_CHUNK_SIZE = utils.ParseEnvInt("TASK_CHUNK_SIZE", TASK_CHUNK_SIZE)
	API_PORT = utils.ParseEnvInt("API_PORT", API_PORT)
	API_MAX_FILE_SIZE_BYTES = utils.ParseEnvInt("API_MAX_FILE_SIZE_BYTES", API_MAX_FILE_SIZE_BYTES)
	API_UPLOADS_DIR = utils.ParseEnvString("API_UPLOADS_DIR", API_UPLOADS_DIR)
	MQ_MAX_NEW_TASKS = utils.ParseEnvInt("MQ_MAX_NEW_TASKS", MQ_MAX_NEW_TASKS)
	MQ_URL = utils.ParseEnvString("MQ_URL", MQ_URL)
	FORGOTTEN_TIMEOUT = utils.ParseEnvDuration("FORGOTTEN_TIMEOUT_MINUTES", FORGOTTEN_TIMEOUT)
	TASK_CHANNEL_QUEUE_SIZE = utils.ParseEnvInt("TASK_CHANNEL_QUEUE_SIZE", TASK_CHANNEL_QUEUE_SIZE)
}
