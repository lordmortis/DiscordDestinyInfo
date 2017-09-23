package discord

import (
  "errors"
  "strings"

  "github.com/bwmarrin/discordgo"
)

var (
  session *discordgo.Session
)

type commandFunc func(*discordgo.Session, *discordgo.MessageCreate)

func Setup(token string) error {
  var err error
  session, err = discordgo.New("Bot " + token)
  if err != nil {
    return errors.New("error creating Discord session, " + err.Error())
  }

  session.AddHandler(handleMessage)

  // Open a websocket connection to Discord and begin listening.
  err = session.Open()
  if err != nil {
    return errors.New("error opening Discord connection," + err.Error())
  }

  return nil
}

func Close() {
  if (session != nil) { session.Close() }
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
