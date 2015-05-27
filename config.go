package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	URL       string
	SessionID string
}

var filename = ".mutterrc"

func SaveConfig(conf *Config) (err error) {
	b, err := json.Marshal(conf)
	if err == nil {
		err = ioutil.WriteFile(filename, b, 0644)
	}

	return
}

func LoadConfig() (conf Config, err error) {
	file, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(file, &conf)
	}

	return
}

func ClearConfig() {
	ioutil.WriteFile(filename, []byte(nil), 0644)
}
