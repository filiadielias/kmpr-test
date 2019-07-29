// Package general contains commonly use methods and variables
// such as config, server module, etc.
package general

import (
	"os"
)

// Config struct stores configuration values
type Config struct {
	App struct {
		Address  string `json:"address"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
	} `json:"app"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	ES struct { //elasticsearch
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"elasticsearch"`
	Redis struct { //redis
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	NSQ struct {
		Producer struct { //producer host
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"producer"`
		Consumer struct { //nsqlookupd host
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"consumer"`
	} `json:"nsq"`
}

// Parse config file into Config struct
func (c *Config) Parse(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err := JSONUnmarshal(f, &c); err != nil {
		return err
	}

	return nil
}
