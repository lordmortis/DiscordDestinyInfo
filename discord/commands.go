package discord

import (
  "fmt"
  "strings"
  "sort"

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
  commandList []string
)

func init() {
  commands = make(map[string]commandData)
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
}

func help(session *discordgo.Session, user *discordgo.User) {
  channel, err := session.UserChannelCreate(user.ID)
  if err != nil {
    fmt.Println("Could not create channel for DM!")
    return
  }

  session.ChannelMessageSend(channel.ID, "Hi!")
  session.ChannelMessageSend(channel.ID, "Here's some help:")
  for _, key := range ( commandList ) {
    var command = commands[key]
    var help = fmt.Sprintf("%s - %s", command.name, command.help)
    session.ChannelMessageSend(channel.ID, help)
  }
}
