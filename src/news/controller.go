// Package news contains business logic from news, store to database, etc
package news

import (
	"errors"
	"fmt"
	//"sync"

	"github.com/filiadielias/kmpr-test/src/general"
	"github.com/filiadielias/kmpr-test/src/helper/db"
	"github.com/filiadielias/kmpr-test/src/helper/elastic"

	"github.com/bitly/go-nsq"
)

// func AddNews(author, body string) error {
// 	n := News{Author: author, Body: body}

// 	// Insert to database
// 	if err := n.insert(); err != nil {
// 		return err
// 	}

// 	data := struct {
// 		ID      int    `json:"id"`
// 		Created string `json:"created"`
// 	}{
// 		n.ID,
// 		// Format timestamp for compatibility with elasticsearch Date format
// 		// make sure the time has 6-digit fractional second
// 		n.Created.Format("2006-01-02 15:04:05.000000"),
// 	}

// 	err := elastic.AddDocument(kmpr.ES, "news", fmt.Sprintf("%d", n.ID), data)
// 	if err != nil {
// 		return err
// 	}

// 	return err
// }

// InsertNews to add news to database
func InsertNews(n *News) error {
	//n := News{Author: "", Body: ""}

	// Insert to database
	if err := n.insert(); err != nil {
		return err
	}

	data := struct {
		ID      int    `json:"id"`
		Created string `json:"created"`
	}{
		n.ID,
		// Format timestamp for compatibility with elasticsearch Date format
		// make sure the time has 6-digit fractional second
		n.Created.Format("2006-01-02 15:04:05.000000"),
	}

	err := elastic.AddDocument(kmpr.ES, "news", fmt.Sprintf("%d", n.ID), data)
	if err != nil {
		return err
	}

	return nil
}

// AddNews to publish news to consumers
func AddNews(author, body string) error {
	n := News{Author: author, Body: body}

	// Publish to NSQ
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(fmt.Sprintf("%s:%d", kmpr.Config.NSQ.Producer.Host, kmpr.Config.NSQ.Producer.Port), config)
	if err != nil {
		return err
	}

	// Encode data to bytes
	b, err := general.GobEncode(n)
	if err != nil {
		return err
	}

	// Publish to nsq
	if err := w.Publish("NEWS_ADD", b); err != nil {
		return err
	}
	w.Stop()

	return nil
}

// GetNews getting list of news by page and speficy size per page
func GetNews(page, size int) (ns Newses, err error) {
	if size <= 0 {
		return ns, errors.New("Invalid size number")
	}

	if page <= 0 {
		page = 1
	}

	// Set filter, sort information
	eq := elastic.Query{}
	eq.Page = page
	eq.Limit = size
	eq.Sort = map[string]string{
		"created": "desc",
	}

	// Get from elastic
	ids, err := elastic.GetDocuments(kmpr.ES, "news", &eq)
	if err != nil {
		return ns, err
	}

	// Append all news first
	for i := 0; i < len(ids); i++ {
		ns = append(ns, News{ID: ids[i]})
	}

	ch := make(chan error)
	for i := 0; i < len(ns); i++ {
		go ns[i].get(ch)
	}

	for i := 0; i < len(ns); i++ {
		err := <-ch
		if err != nil {
			return ns, err
		}
	}

	return ns, nil
}

// Get news detail
func (n *News) get(ch chan<- error) {
	qb := db.QueryBuilder{}
	qb.AddFilter("id", n.ID, "")

	err := n.getFromDB(&qb)
	if err != nil {
		ch <- err
		return
	}

	ch <- nil
	return
}
