// Package handler contains http handlers and NSQ consumers
package handler

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/filiadielias/kmpr-test/src/general"
	"github.com/filiadielias/kmpr-test/src/news"

	"github.com/bitly/go-nsq"
)

type messageHandler struct{}

// HandleMessage to handle AddNews nsq consumer
func (h *messageHandler) HandleMessage(m *nsq.Message) error {

	if len(m.Body) == 0 {
		return fmt.Errorf("body is blank")
	}

	var n news.News

	// Decode message
	if err := general.GobDecode(m.Body, &n); err != nil {
		return err
	}

	// Insert data
	if err := news.InsertNews(&n); err != nil {
		return err
	}

	return nil
}

// InitNSQHandlers initialized nsq consumers then
// make connection to NSQ server
func InitNSQHandlers(host string, port int) error {
	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer("NEWS_ADD", "database", config)
	if err != nil {
		log.Panicf("fail to init NSQ consumer: %v", err)
		return err
	}

	consumer.ChangeMaxInFlight(100)

	consumer.AddConcurrentHandlers(
		&messageHandler{},
		20,
	)

	if err := consumer.ConnectToNSQLookupds([]string{fmt.Sprintf("%s:%d", host, port)}); err != nil {
		log.Panicf("fail to connect to nsqlookupd: %v", err)
		return err
	}

	// Listen to SIGINT (ctrl+c) to make sure the queues finish properly on shutdown
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	for {
		select {
		case <-consumer.StopChan:
			return nil // Consumer disconnected
		case <-shutdown:
			// Synchronously drain the queue before falling out of main
			consumer.Stop()
		}
	}

	return nil
}
