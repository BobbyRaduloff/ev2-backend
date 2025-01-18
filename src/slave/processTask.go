package main

import (
	"context"
	"sync"

	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/types"
	"github.com/CorporateBusinessTechnologies/email-verifier-v2/src/utils"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func ProcessTask(emails []string, requestId int, taskId string, channel *amqp.Channel) []types.ProcessingResult {
	// set up so that it runs only 20 emails at a time
	var wg sync.WaitGroup
	resultsChannel := make(chan types.ProcessingResult, len(emails)+1)
	semaphore := make(chan struct{}, MAX_THREADS)
	utils.Logger.Info("beginning task processing", zap.String("task", taskId), zap.Int("maxThreads", MAX_THREADS))

	// process them all, up to 20 at a time
	for _, email := range emails {
		wg.Add(1)
		utils.Logger.Info("started processing email from task", zap.String("task", taskId), zap.String("email", email))

		// process one email
		go func(email string) {
			// mark as done at end
			defer wg.Done()

			// allocate a thread
			semaphore <- struct{}{}

			// set up timeout per email
			ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultTimeout)
			defer cancel()

			// process, return results, and free up max threads
			result, err := ProcessEmail(ctx, email, requestId)
			if err != nil {
				utils.Logger.Error("cant process email", zap.Error(err), zap.String("email", email), zap.Int("request", requestId))
			}

			resultsChannel <- result

			<-semaphore
		}(email)
	}

	// finish all emails
	wg.Wait()
	utils.Logger.Info("finished processing task", zap.String("task", taskId))

	// get emails in an array
	results := []types.ProcessingResult{}
	for i := 0; i < len(emails); i++ {
		results = append(results, <-resultsChannel)
	}
	utils.Logger.Info("finished collecting results", zap.String("task", taskId))

	close(resultsChannel)
	return results
}
