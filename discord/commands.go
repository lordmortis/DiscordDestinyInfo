package discord

import (
  //"fmt"
  "strings"

  "github.com/bwmarrin/discordgo"
)

type Command func(*discordgo.Session, *discordgo.Message)

type commandData struct {
  name string
  help string
  fn   Command
}

var (
  commands map[string]commandData
)

func init() {
  commands = make(map[string]commandData)
}

func RegisterCommand(name string, help string, fn Command) {
  commands[strings.ToLower(name)] = commandData{
    name: name,
    help: help,
    fn: fn,
  }
}
