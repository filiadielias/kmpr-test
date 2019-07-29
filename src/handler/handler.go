// Package handler : Initialize http handlers and NSQ consumers
package handler

import (
	news_handler "github.com/filiadielias/kmpr-test/src/news/handler"

	"github.com/julienschmidt/httprouter"
)

// Init endpoints, do not add global middleware here
func initHandlers() *httprouter.Router {

	router := httprouter.New()
	router.POST("/news", news_handler.AddNewsHandler)
	router.GET("/news", news_handler.GetNewsHandler)

	return router
}

// GetHandlers : Get endpoint handlers
// global middleware handlers should be added here
func GetHandlers() *httprouter.Router {
	return initHandlers()
}

// StartNSQConsumer : Similar to http.ListenAndServe, but for NSQ consumer
func StartNSQConsumer(host string, port int) error {
	err := news_handler.InitNSQHandlers(host, port)
	if err != nil {
		return err
	}

	return nil
}
