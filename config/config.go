package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type config struct {
	Request struct {
		FetchRate	int `yaml:"fetchRate"`
		DelayGap	int `yaml:"delayGap"`	
	} `yaml:"request"`
}

var cfg config
var once sync.Once

func GetConfig() *config {
	once.Do(func() {
		readFile(&cfg)
	})
	return &cfg
}

func processError(err error) {
    fmt.Println(err)
    os.Exit(2)
}

func readFile(cfg *config) {
    f, err := os.ReadFile("S:/DevProjects/event-dispatcher/config/config.yaml")
    if err != nil {
        processError(err)
    }

    unmarshalErr := yaml.Unmarshal(f, cfg)
    if unmarshalErr != nil {
        processError(err)
    }
} 

func (c config) FetchRate() int {
    return c.Request.FetchRate
}

func (c config) GetDatabase() Database {
	return c.Database
}