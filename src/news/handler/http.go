// Package handler contains http handlers and NSQ consumers
package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/filiadielias/kmpr-test/src/general"
	"github.com/filiadielias/kmpr-test/src/helper/redis"
	writer_lib "github.com/filiadielias/kmpr-test/src/helper/writer"
	"github.com/filiadielias/kmpr-test/src/news"

	"github.com/julienschmidt/httprouter"
)

var kmpr *general.Module

func init() {
	kmpr = &general.KMPR
}

// AddNewsHandler : to handle add news endpoint
func AddNewsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	writer := writer_lib.New(w)

	var n news.News

	if err := general.JSONUnmarshal(r.Body, &n); err != nil {
		log.Println(err)
		writer.Error(err)
		return
	}

	// Insert into database
	if err := news.AddNews(n.Author, n.Body); err != nil {
		log.Println(err)
		writer.Error(err)
		return
	}
	writer.Success(nil)

	// Delete stored cache (data is not up to date)
	keys, err := redis.GetKeys(kmpr.Redis, "news:search:page:*")
	if err != nil {
		log.Println(err)
		writer.Error(err)
		return
	}

	ch := make(chan error)
	for i := 0; i < len(keys); i++ {

		go func(value string, ch chan<- error) {

			err := redis.Delete(kmpr.Redis, value)
			ch <- err
		}(keys[i], ch)
	}

	for i := 0; i < len(keys); i++ {
		err = <-ch
		log.Println(err)
	}
}

// GetNewsHandler : get news by page number
func GetNewsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	writer := writer_lib.New(w)

	page := 0
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	var ns news.Newses

	// Check cache
	key := fmt.Sprintf("news:search:page:%d", page)
	exists, err := redis.Exists(kmpr.Redis, key)
	if err != nil {
		log.Println(err)
		writer.Error(err)
		return
	}
	// Get cache
	if exists {
		err = redis.GetStruct(kmpr.Redis, key, &ns)

		// If not error, return the data
		if err == nil {
			writer.Success(ns)
			return
		}

		// If error, display error then continue fetching from database
		log.Println(err)
	}

	// Fetching from database
	ns, err = news.GetNews(page, 10) // Hardcoded temporarily
	if err != nil {
		log.Println(err)
		writer.Error(err)
		return
	}

	// Store to cache server
	if err := redis.SetStruct(kmpr.Redis, key, ns); err != nil {
		// Just display the error
		log.Println(err)
	}

	writer.Success(ns)
}
