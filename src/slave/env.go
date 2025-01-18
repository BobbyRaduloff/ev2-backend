package main

import (
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

var MQ_MAX_NEW_TASKS = 1
var MQ_URL = "amqp://admin:admin@localhost:5672/"
var MAX_THREADS = 25
var LOGGER_SYNC_TIMEOUT = 60 * time.Second

func ReadEnv() {
	MQ_MAX_NEW_TASKS = utils.ParseEnvInt("MQ_MAX_NEW_TASKS", MQ_MAX_NEW_TASKS)
	MQ_URL = utils.ParseEnvString("MQ_URL", MQ_URL)
	MAX_THREADS = utils.ParseEnvInt("MAX_THREADS", MAX_THREADS)
	LOGGER_SYNC_TIMEOUT = utils.ParseEnvDuration("LOGGER_SYNC_TIMEOUT_SECONDS", LOGGER_SYNC_TIMEOUT)
}
