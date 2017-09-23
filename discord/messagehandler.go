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

  if handleCommand(suffix, s, m) { return }

  fmt.Println("Didn't recognize: ", suffix)
  s.ChannelMessageSend(m.ChannelID, "Sorry I didn't understand \"" + m.Content +"\" - ask me for \"help\" if you need")
}
