// Package redis contains helper functions for operating with redis
package redis

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// Connect to redis server
func Connect(host string, port int) *redis.Pool {
	return newPool(host, port)
}

func newPool(host string, port int) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				panic(err.Error())
			}
			return c, nil
		},
	}
}

// Get redis string value
func Get(pool *redis.Pool, key string) (string, error) {

	conn := pool.Get()
	defer conn.Close()

	// Perform request and convert result to string
	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error getting key %s: %v", key, err)
	}
	return data, nil
}

// GetStruct : get redis value and convert to struct
func GetStruct(pool *redis.Pool, key string, data interface{}) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	s, err := Get(pool, key)
	if err != nil {
		return fmt.Errorf("error getting key %s: %v", key, err)
	}

	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return fmt.Errorf("error unmarshal key %s data %s: %v", key, s, err)
	}

	return nil
}

// Set is adding new redis key
func Set(pool *redis.Pool, key string, value string) error {

	conn := pool.Get()
	defer conn.Close()

	// Perform request
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("error setting key %s to %s: %v", key, value, err)
	}
	return nil
}

// SetStruct is adding new redis key from struct
func SetStruct(pool *redis.Pool, key string, data interface{}) error {
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data %v: %v", data, err)
	}

	return Set(pool, key, string(b))
}

// Exists : Check if key exists
func Exists(pool *redis.Pool, key string) (bool, error) {

	conn := pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("error checking if key %s exists: %v", key, err)
	}
	return ok, nil
}

// Delete redis key
func Delete(pool *redis.Pool, key string) error {

	conn := pool.Get()
	defer conn.Close()

	// Perform request and convert result to string
	_, err := conn.Do("DEL", key)
	return err
}

// GetKeys : get redis keys using pattern
func GetKeys(pool *redis.Pool, pattern string) ([]string, error) {

	conn := pool.Get()
	defer conn.Close()

	iter := 0
	keys := []string{}

	// Loop until scan cursor return 0
	for {
		// Perform request
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", pattern))
		if err != nil {
			return keys, fmt.Errorf("error retrieving '%s' keys", pattern)
		}

		// Get cursor
		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
