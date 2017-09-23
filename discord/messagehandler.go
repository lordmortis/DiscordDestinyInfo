package discord

import (
  "strings"

  "fmt"

  "github.com/bwmarrin/discordgo"
)

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID { return }

  var channelInfo = getChannelInfo(s, m.ChannelID)

  var suffix string
  var mymessage bool

  if channelInfo.DM {
    suffix = strings.Trim(m.Content, " ")
    mymessage = true
  } else {
    suffix, mymessage = discordGetCommand(s.State.User, m);
  }

  if !mymessage { return }

  var command = strings.ToLower(suffix)

  if strings.HasPrefix(command, "help") {
    help(s, m.Author)
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

  fmt.Println("Didn't recognize: ", suffix)
  s.ChannelMessageSend(m.ChannelID, "Sorry I didn't understand \"" + command +"\" - ask me for \"help\" if you need")
}
