// Package general contains commonly use methods and variables
// such as config, server module, etc.
package general

import (
	"github.com/elastic/go-elasticsearch"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// KMPR is a global variable used for getting config values, establish db, redis and elasticsearch connections.
var KMPR Module

// Module stores DB connection, configuration, redis pool and elasticsearchclient.
type Module struct {
	DB     *sqlx.DB
	Config Config
	ES     *elasticsearch.Client
	Redis  *redis.Pool
}

// New is adding config and connections to global module
func New(db *sqlx.DB, config Config, es *elasticsearch.Client, red *redis.Pool) {
	KMPR = Module{
		DB:     db,
		Config: config,
		ES:     es,
		Redis:  red,
	}
}
