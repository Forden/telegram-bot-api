package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type MyConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
	Auth struct {
		AppId   int    `json:"app_id"`
		AppHash string `json:"app_hash"`
	} `json:"auth"`
}

func loadConfig(filename string) (MyConfig, error) {
	var cfg MyConfig
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&cfg)
	if err != nil {
		return MyConfig{}, err
	} else {
		return cfg, nil
	}
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	} else {
		router := NewRouter(cfg)
		address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		log.Printf("Starting server on: %s", address)
		log.Fatal(http.ListenAndServe(address, router))
	}

}
