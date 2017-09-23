package discord

import (
  "strings"

  "fmt"

  "github.com/bwmarrin/discordgo"
)

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID { return }

  var channelInfo = getChannelInfo(s, m.ChannelID)

  var command string
  var commandFound bool

  if channelInfo.DM {
    command = strings.Trim(m.Content, " ")
    commandFound = true
  } else {
    command, commandFound = discordGetCommand(s.State.User, m);
  }

  if !commandFound {
    fmt.Printf("Saw message: '%s'\n", m.Content)
    return
  }

  fmt.Println("Saw command: ", command)

  if strings.HasPrefix(command, "help") {
    discordSendHelp(s, m.Author)
    return
  }

  // If the message is "ping" reply with "Pong!"
  if strings.HasPrefix(command, "ping") {
    s.ChannelMessageSend(m.ChannelID, "Pong " + m.Author.Username)
    return
  }

  // If the message is "pong" reply with "Ping!"
  if strings.HasPrefix(command, "pong") {
    s.ChannelMessageSend(m.ChannelID, "Pong " + m.Author.Username)
    return
  }

  s.ChannelMessageSend(m.ChannelID, "Sorry I didn't understand \"" + command +"\" - ask me for \"help\" if you need")
}
