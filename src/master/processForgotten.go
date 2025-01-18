package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/env"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/mq"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/rules"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/sql/dtos/emails"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

func ProcessForgotten(db *sqlx.DB, channel *amqp.Channel, uploadMutex *sync.Mutex) {
	for {
		uploadMutex.Lock()

		utils.Logger.Info("processing forgotten guys")

		forgottenEmails, err := emails.GetAllUnfinshed(db)
		if err != nil {
			utils.Logger.Error("error when getting unfinished emails", zap.Error(err))
			return
		}

		emailMap := make(map[int][]emails.EmailDTO)
		for _, email := range forgottenEmails {
			requestID := email.RequestID
			emailMap[requestID] = append(emailMap[requestID], email)
		}

		utils.Logger.Info("got unfinished emails", zap.Int("length", len(forgottenEmails)))

		for requestId, emailDTOS := range emailMap {
			utils.Logger.Info("sending unfinished emails from request to process", zap.Int("requestId", requestId), zap.Int("length", len(emailDTOS)))
			emailsAccounts := make([]string, 0, len(emailDTOS))
			for _, e := range emailDTOS {
				if !rules.IsEmailValid(e.Email) {
					continue
				}
				emailsAccounts = append(emailsAccounts, e.Email)
			}

			chunkedEmails := utils.ChunkArray(emailsAccounts, env.TASK_CHUNK_SIZE)
			chunkedDTOs := utils.ChunkArray(emailDTOS, env.TASK_CHUNK_SIZE)
			utils.Logger.Info("chunked unfinished emails from request", zap.Int("requestId", requestId), zap.Int("chunks", len(chunkedDTOs)))

			for i, chunk := range chunkedEmails {
				body := strings.Join(chunk, ",")
				mq.SendMessage(fmt.Sprintf("%d,%s", requestId, body), utils.ToProcessQueueName, channel)
				emails.BatchUpdateEmailsTimestamp(db, chunkedDTOs[i])
			}
		}

		uploadMutex.Unlock()
		time.Sleep(env.FORGOTTEN_TIMEOUT)
	}
}
