package main

import (
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/env"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/emails"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/requests"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/types"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

var DBResultsChannel chan types.ProcessingResult

func InitDBResultsChannel() {
	DBResultsChannel = make(chan types.ProcessingResult, env.DB_INCOMING_QUEUE_MAX_SIZE)
}

func ReceiveResults(r <-chan amqp.Delivery) {
	for result := range r {
		utils.Logger.Info("recieved task result", zap.String("task", result.MessageId))

		deserialized, err := types.DeserializeEmailResults(string(result.Body))
		if err != nil {
			utils.Logger.Error("failed to deserialize task result", zap.String("task", result.MessageId), zap.Int("requestId", deserialized.RequestId))
			continue
		}

		for _, dr := range deserialized.Results {
			dr.RequestId = deserialized.RequestId
			DBResultsChannel <- dr
		}
		utils.Logger.Info("pushed task result to queue", zap.String("task", result.MessageId))

		result.Ack(false)
	}
}

func ProcessResults(db *sqlx.DB) {
	for {
		utils.Logger.Info("processing results")

		toWrite := make([]types.ProcessingResult, 0, 100)
		timeout := time.After(env.DB_INCOMING_QUEUE_TIMEOUT)

	loop:
		for {
			select {
			case res := <-DBResultsChannel:
				toWrite = append(toWrite, res)
			case <-timeout:
				break loop
			}
		}

		if len(toWrite) > 0 {
			err := emails.BatchAddResults(db, toWrite)
			if err != nil {
				utils.Logger.Error("error when writing results to db", zap.Error(err))
			}

			updatedRequests := make(map[int]int)
			for _, n := range toWrite {
				updatedRequests[n.RequestId] = 1
			}

			for requestId := range updatedRequests {
				err := requests.UpdateCountsById(db, requestId)
				if err != nil {
					utils.Logger.Error("error when updating counts to db", zap.Error(err), zap.Int("requestId", requestId))
				}
			}
		}
	}
}
