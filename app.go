package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filiadielias/kmpr-test/src/general"
	"github.com/filiadielias/kmpr-test/src/handler"
	"github.com/filiadielias/kmpr-test/src/helper/db"
	"github.com/filiadielias/kmpr-test/src/helper/elastic"
	"github.com/filiadielias/kmpr-test/src/helper/redis"
)

var kmpr *general.Module

func init() {
	var c general.Config
	if err := initConfig(&c); err != nil {
		log.Fatal("Fail to load config: ", err)
		return
	}

	dbconn, err := db.InitDB(c.Database.Host, c.Database.Port).
		Credential(c.Database.User, c.Database.Password).
		Database(c.Database.DBName).
		Connect()
	if err != nil {
		log.Fatal("Fail to connect to database: ", err)
		return
	}

	es, err := elastic.Connect(c.ES.Host, c.ES.Port)
	if err != nil {
		log.Fatal("Fail to connect to elasticsearch service: ", err)
		return
	}

	red := redis.Connect(c.Redis.Host, c.Redis.Port)

	general.New(dbconn, c, es, red)

	kmpr = &general.KMPR
}

func initConfig(c *general.Config) error {
	confFile := "./files/config/"

	switch env := os.Getenv("APPSENV"); env {
	case "development":
		fallthrough
	case "staging":
		fallthrough
	case "production":
		confFile += env
	default:
		confFile += "development"
	}

	confFile += ".json"

	return c.Parse(confFile)
}

func main() {

	address := fmt.Sprintf("%s:%d", kmpr.Config.App.Host, kmpr.Config.App.Port)
	router := handler.GetHandlers()

	//start NSQ Consumer
	go func() {
		log.Fatal(handler.StartNSQConsumer(kmpr.Config.NSQ.Consumer.Host, kmpr.Config.NSQ.Consumer.Port))
	}()

	log.Fatal(http.ListenAndServe(address, router))

}
