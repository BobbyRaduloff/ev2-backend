package mq

import (
	"context"
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func GetConnectionAndChannel(url string) (*amqp.Connection, *amqp.Channel) {
	connection, err := amqp.Dial(url)
	if err != nil {
		utils.Logger.Fatal("can't connect to rabbit", zap.Error(err))
	}

	channel, err := connection.Channel()
	if err != nil {
		utils.Logger.Fatal("can't create rabbit channel", zap.Error(err))
	}

	return connection, channel
}

func SetQos(channel *amqp.Channel, maxNewTasks int) {
	err := channel.Qos(maxNewTasks, 0, true)
	if err != nil {
		utils.Logger.Fatal("can't set rabbit qos", zap.Error(err))
	}
}

func GetQueue(channel *amqp.Channel, name string) amqp.Queue {
	queue, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		utils.Logger.Fatal("can't get outgoing task queue", zap.Error(err))
	}

	return queue
}

func GetMessages(queueName string, channel *amqp.Channel) <-chan amqp.Delivery {
	messages, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		utils.Logger.Fatal("can't get rabbit queue", zap.Error(err), zap.String("queue", queueName))
	}

	return messages
}

func SendMessage(data string, queueName string, channel *amqp.Channel) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	taskUUID, err := uuid.NewRandom()
	if err != nil {
		utils.Logger.Error("can't generate uuid for task")
	}

	err = channel.PublishWithContext(ctx,
		"", // exchange
		queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
			MessageId:   taskUUID.String(),
		})
	if err != nil {
		utils.Logger.Fatal("failed to send message on rabbit", zap.Error(err), zap.String("queue", queueName), zap.String("task", taskUUID.String()))
	}

	return taskUUID.String()
}
