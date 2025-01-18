package api

import (
	"fmt"
	"log"
	"sync"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/master/env"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/mq"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func SetUpAndRunAPI(db *sqlx.DB, channel *amqp.Channel, uploadMutex *sync.Mutex) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = int64(env.API_MAX_FILE_SIZE_BYTES)

	state := &State{db: db, channel: channel, uploadMutex: uploadMutex, taskChannel: make(chan string, env.TASK_CHANNEL_QUEUE_SIZE)}

	go func() {
		log.Print("TASKS")
		for task := range state.taskChannel {
			log.Print(task)
			taskID := mq.SendMessage(task, utils.ToProcessQueueName, state.channel)
			utils.Logger.Info("sent task", zap.String("task", taskID))
		}
	}()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "f8eae4065791f7991e360d5c26d0af3d",
	}))

	authorized.POST("/upload", state.PostUpload)
	authorized.POST("/emails", state.PostGetEmails)
	authorized.GET("/filter", state.GetFilter)
	authorized.POST("/filter", state.PostFilter)
	authorized.GET("/requests", state.GetRequests)
	authorized.GET("/requests/:request_id", state.GetRequest)
	authorized.GET("/requests/:request_id/csv", state.GetRequestCSV)
	authorized.GET("/requests/:request_id/csv/filtered/default", state.GetDefaultFilteredCSV)
	authorized.POST("/requests/:request_id/csv/filtered", state.PostGetFilteredEmailsCSV)

	r.Run(fmt.Sprintf(":%d", env.API_PORT))
}
