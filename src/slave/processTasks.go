package main

import (
	"strconv"
	"strings"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/mq"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/types"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func ProcessTasks(tasks <-chan amqp.Delivery, channel *amqp.Channel) {
	// for all tasks, although it should be one
	for task := range tasks {
		if len(task.Body) <= 0 {
			utils.Logger.Info("got an empty task, acknowledging and moving on", zap.String("task", task.MessageId))
			task.Ack(false)
			continue
		}

		// get emails from tasks
		deserialized := strings.Split(string(task.Body), ",")
		deserialized = utils.FilterEmptyStrings(deserialized)
		if len(deserialized) <= 1 {
			utils.Logger.Error("empty task", zap.String("task", task.MessageId))
			task.Ack(false)
			continue
		}

		requestId := deserialized[0]
		requestIdInt, err := strconv.Atoi(requestId)
		if err != nil {
			utils.Logger.Error("malformed id", zap.String("task", task.MessageId))
			task.Ack(false)
			continue
		}
		emails := deserialized[1:]

		results := ProcessTask(emails, requestIdInt, task.MessageId, channel)

		// serialize emails to json so we can send back
		serializedResults, err := (&types.ProcessingResults{RequestId: requestIdInt, Results: results}).SerializeEmailResults()
		if err != nil {
			utils.Logger.Error("couldn't serialize task results to json", zap.String("task", task.MessageId), zap.Error(err))
		}

		// send it back to mq
		mq.SendMessage(serializedResults, utils.ProcessedQueueName, channel)
		utils.Logger.Info("finished sending results", zap.String("task", task.MessageId))

		// acknowledge task
		task.Ack(false)
		utils.Logger.Info("acknowledged task as done to mq", zap.String("task", task.MessageId))
	}

}
