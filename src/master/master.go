package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/api"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/env"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/mq"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

func main() {
	// read env
	env.ReadEnv()

	// set up logging
	utils.CreateLogger()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			utils.Logger.Sync()
			os.Exit(1)
		}
	}()
	go utils.SyncLogger(env.LOGGER_SYNC_TIMEOUT)

	// set up mq
	connection, channel := mq.GetConnectionAndChannel(env.MQ_URL)
	defer connection.Close()
	defer channel.Close()
	mq.SetQos(channel, env.MQ_MAX_NEW_TASKS)
	mq.GetQueue(channel, utils.ToProcessQueueName)
	mq.GetQueue(channel, utils.ProcessedQueueName)

	// set up db
	db := sql.ConnectToPostgres(env.DB_HOST, env.DB_PORT, env.DB_USER, env.DB_PASS, env.DB_DATABASE, env.DB_SSL)

	// get results incoming channel and pipe
	results := mq.GetMessages(utils.ProcessedQueueName, channel)
	InitDBResultsChannel()

	// run forever
	var forever chan struct{}

	// only one writer with new shit to db and taskqueue
	var uploadMutex sync.Mutex
	go api.SetUpAndRunAPI(db, channel, &uploadMutex)

	// recieve results
	go ReceiveResults(results)

	// write results to db
	go ProcessResults(db)

	// check for forgotten cunts every 5 minutes
	time.Sleep(env.DB_INCOMING_QUEUE_TIMEOUT / 3 * 2)
	go ProcessForgotten(db, channel, &uploadMutex)

	<-forever
}
