package discord

import (
  "strings"
  "sort"

  "github.com/bwmarrin/discordgo"
)

type Command func(session *discordgo.Session, message *discordgo.Message, parameters string)

type commandData struct {
  name string
  help string
  fn   Command
}

var (
  commands map[string]commandData
  commandList []string
  maxCommandLength int
)

func init() {
  commands = make(map[string]commandData)
  RegisterCommand("Help", "show help", helpCommand)
}

func RegisterCommand(name string, help string, fn Command) {
  var lowerName = strings.ToLower(name)
  commands[lowerName] = commandData{
    name: name,
    help: help,
    fn: fn,
  }

  commandList = append(commandList, lowerName)
  sort.Strings( commandList )

  if len(lowerName) > maxCommandLength { maxCommandLength = len(lowerName) }
}

func handleCommand(message string, s *discordgo.Session, m *discordgo.MessageCreate) bool {
  var substrings = strings.SplitN(message, " ", 2)
  data, exists := commands[strings.ToLower(substrings[0])]

  if !exists { return false }
  if len(substrings) == 1 {
    data.fn(s, m.Message, "")
  } else {
    data.fn(s, m.Message, substrings[1])
  }

  return true
}
