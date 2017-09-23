package main

import (
  "encoding/json"
//  "fmt"
  "io/ioutil"
//  "strconv"
)

type DiscordConfig struct {
  Token    string
  ClientID string
}

type Config struct {
  Discord      DiscordConfig           `json:"discord"`
}

func LoadConfig(filename string) (*Config, error) {
  filestring, err := ioutil.ReadFile(filename)
  if (err != nil) {
    return nil, err
  }

  var config = Config{}

  err = json.Unmarshal(filestring, &config)
  if (err != nil) {
    return nil, err
  }

  return &config, nil
}


