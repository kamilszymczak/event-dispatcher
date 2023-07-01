package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type config struct {
	Request struct {
		FetchRate int `yaml:"fetchRate"`
	} `yaml:"request"`
}

var cfg *config
var once sync.Once

func GetConfig() config {
	once.Do(func() {
		readFile(cfg)
	})
	return *cfg
}

func processError(err error) {
    fmt.Println(err)
    os.Exit(2)
}

func readFile(cfg *config) {
    f, err := os.Open("config.yml")
    if err != nil {
        processError(err)
    }
    defer f.Close()

    decoder := yaml.NewDecoder(f)
    err = decoder.Decode(cfg)
    if err != nil {
        processError(err)
    }
} 