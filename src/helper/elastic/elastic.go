// Package elastic contains helper functions for operating with elasticsearch
package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

// List of common error
var (
	ErrInvalidIndex = fmt.Errorf("Invalid index")
	ErrNilData      = fmt.Errorf("Data value must not be nil")
	ErrInvalidID    = fmt.Errorf("Invalid ID format")
)

// Connect to elasticsearch
func Connect(host string, port int) (*elasticsearch.Client, error) {
	return newClient(host, port)
}

func newClient(host string, port int) (*elasticsearch.Client, error) {

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s:%d", host, port),
		},
	})
	if err != nil {
		return nil, err
	}

	// Test connection
	_, err = es.Info()
	if err != nil {
		return nil, err
	}

	return es, nil
}

// AddDocument for adding document to specified index
func AddDocument(es *elasticsearch.Client, index, id string, data interface{}) error {
	if data == nil {
		return ErrNilData
	}

	// Id must be number and cannot zero and below
	if idVal, err := strconv.Atoi(id); err != nil || idVal <= 0 {
		return ErrInvalidID
	}

	// Index cannot be empty
	if len(index) == 0 {
		return ErrInvalidIndex
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err)
	}

	// Set up the request object
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}

	// Perform the request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	// Check if
	if res.IsError() {
		return fmt.Errorf("[%s] error indexing document ID=%s", res.Status(), id)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("[%s] error parsing the response body: %s", res.Status(), err)
	}

	results, ok := r["_shards"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("insert failed")
	}

	if results["successful"].(float64) <= 0 {
		return fmt.Errorf("insert failed")
	}

	return nil
}

// Query is simple elasticsearch search query format
type Query struct {
	// If page is zero, get all data
	Page int `json:"from,omitempty"`
	// If limit is zero, get all data
	Limit int               `json:"size,omitempty"`
	Sort  map[string]string `json:"sort"`
}

// GetJSON : build Elasticsearch search query
func (q *Query) GetJSON() (string, error) {
	if q.Limit > 0 {
		if q.Page > 0 {
			q.Page--

			// Set FROM
			q.Page *= q.Limit
		}
	}

	b, err := json.Marshal(q)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// GetDocuments : get elasticsearch document from specified index
func GetDocuments(es *elasticsearch.Client, index string, q *Query) (ids []int, err error) {
	if len(index) == 0 {
		return ids, ErrInvalidIndex
	}

	s, err := q.GetJSON()
	if err != nil {
		return ids, err
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(strings.NewReader(s)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return ids, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return ids, fmt.Errorf("error parsing the response body: %s", err)
	}

	if res.IsError() {
		// Print the response status and error information.
		return ids, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			r["error"].(map[string]interface{})["type"],
			r["error"].(map[string]interface{})["reason"],
		)
	}

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, _ := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))

		ids = append(ids, id)
	}

	return ids, nil
}
