package api

import (
	"sync"

	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type State struct {
	db          *sqlx.DB
	channel     *amqp.Channel
	uploadMutex *sync.Mutex
	taskChannel chan string
}
