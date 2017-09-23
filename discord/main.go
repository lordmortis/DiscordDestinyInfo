package discord

import (
  "errors"
  "strings"

  "fmt"
//  "strconv"

  "github.com/bwmarrin/discordgo"
)

var (
  discord *discordgo.Session
  discordChannels map[string]channelInfo
)

type channelInfo struct {
  DM bool
}

type commandFunc func(*discordgo.Session, *discordgo.MessageCreate)

func Setup(token string) error {
  var err error
  discord, err = discordgo.New("Bot " + token)
  if err != nil {
    return errors.New("error creating Discord session, " + err.Error())
  }

  discord.AddHandler(discordNewMessage)

  // Open a websocket connection to Discord and begin listening.
  err = discord.Open()
  if err != nil {
    return errors.New("error opening Discord connection," + err.Error())
  }

  discordChannels = make(map[string]channelInfo)

  return nil
}

func Close() {
  if (discord != nil) { discord.Close() }
}

func discordGetCommand(user *discordgo.User, message *discordgo.MessageCreate) (string, bool) {
  var searchprefixes []string = make([]string, 3)
  var command string
  commandFound := false

  searchprefixes[0] = "<@" + user.ID + ">"
  searchprefixes[1] = "@" + user.Username
  searchprefixes[2] = user.Username

  for _, prefix := range searchprefixes {
   if strings.HasPrefix(message.Content, prefix) {
      command = strings.TrimPrefix(message.Content, prefix)
      commandFound = true
      break
    }
  }

  if commandFound {
    command = strings.Trim(command, " ")
  }

  return command, commandFound
}

func discordSendHelp(session *discordgo.Session, user *discordgo.User) {
  channel, err := session.UserChannelCreate(user.ID)
  if err != nil {
    fmt.Println("Could not create channel!")
    return
  }

  session.ChannelMessageSend(channel.ID, "Hi!")
  session.ChannelMessageSend(channel.ID, "send me 'ping' or 'pong' to have me send you a message!")
}

func discordNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID { return }

  var channelInfo = discordGetChannelInfo(s, m.ChannelID)

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
