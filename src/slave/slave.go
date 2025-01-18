package main

import (
	"os"
	"os/signal"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/mq"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
)

func main() {
	// read env
	ReadEnv()

	// set up logging
	utils.CreateLogger()
	go utils.SyncLogger(LOGGER_SYNC_TIMEOUT)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			utils.Logger.Sync()
			os.Exit(1)
		}
	}()

	// set up mq
	connection, channel := mq.GetConnectionAndChannel(MQ_URL)
	defer connection.Close()
	defer channel.Close()
	mq.SetQos(channel, MQ_MAX_NEW_TASKS)
	mq.GetQueue(channel, utils.ToProcessQueueName)
	mq.GetQueue(channel, utils.ProcessedQueueName)
	mq.GetQueue(channel, utils.ErrorsQueueName)
	tasks := mq.GetMessages(utils.ToProcessQueueName, channel)

	// run forever async
	var forever chan struct{}
	go ProcessTasks(tasks, channel)

	<-forever
}
