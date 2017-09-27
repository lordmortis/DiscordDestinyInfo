package main

import (
  "fmt"
  "strings"
  "strconv"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
  "github.com/lordmortis/goBungieNet"

  "github.com/bwmarrin/discordgo"
)

var (
)

func init() {
  discord.RegisterCommand("RegisterSearch", "Search for your PSN/Xbox account (if only one result matches, will register) - `RegisterSearch <GamerTag/Nickname> <Xbox/Psn/BattleNet>`", handleRegisterSearch)
  discord.RegisterCommand("Register", "Register your bungo account - `Register <Bungie.Net Membership ID> <Xbox/Psn/BattleNet>`", handleRegister)
  discord.RegisterCommand("ShowRegistration", "Show your registration if it exists", handleRegisterShow)
}

func handleRegister(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil {
    fmt.Println("Could not create channel for DM!")
    return
  }

  paramParts := strings.Split(parameters, " ")

  if len(paramParts) < 2 {
    msg := "Please supply your Bungie.Net ID - `Register <Bungie.Net Membership ID>`"
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  var id int64
  id, err = strconv.ParseInt(paramParts[0], 10, 64)
  if err != nil {
    msg := "didn't recognize member ID - `Register <Bungie.Net Membership ID> <Xbox/Psn/BattleNet>`"
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  accountType := goBungieNet.StringToBungieMembershipType(paramParts[1])

  if accountType == goBungieNet.NetworkNone {
    msg := "didn't recognize account type - `Register <Bungie.Net Membership ID> <Xbox/Psn/BattleNet>`"
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  components := []goBungieNet.DestinyComponentType{goBungieNet.ComponentCharacters}
  _, err1 := goBungieNet.GetProfile(id, accountType, components)
  if err1 != nil {
    msg := fmt.Sprintf("Error registering: %s", err1.Error())
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  err = createRego(message.Author.ID, id, accountType)
  if err != nil {
    msg := fmt.Sprintf("Error registering: %s", err.Error())
    session.ChannelMessageSend(channel.ID, msg)
    return
  }
}

func handleRegisterSearch(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil {
    fmt.Println("Could not create channel for DM!")
    return
  }

  paramParts := strings.Split(parameters, " ")

  if len(paramParts) < 2 {
    msg := "Please supply your nickname and account type - `RegisterSearch <GamerTag/Nickname> <Xbox/Psn/BattleNet>`"
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  accountType := goBungieNet.StringToBungieMembershipType(paramParts[1])

  if accountType == goBungieNet.NetworkNone {
    msg := "didn't recognize account type - `RegisterSearch <GamerTag/Nickname> <Xbox/Psn/BattleNet>`"
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  msg := fmt.Sprintf("Searching %s for %s", accountType, paramParts[0])
  session.ChannelMessageSend(message.ChannelID, msg)

  users, err := goBungieNet.FindDestinyUser(paramParts[0], accountType)
  if err != nil {
    msg := fmt.Sprintf("Could not search for player on Bungie.net: %s", err)
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  if len(*users) == 0 {
    msg := fmt.Sprintf("Could not find any matching users :(")
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  if len(*users) == 1 {
    user := (*users)[0]
    err := createRego(message.Author.ID, user.MembershipID, user.MembershipType)
    if err != nil {
      session.ChannelMessageSend(channel.ID, "Could not record registration :(")
      return
    }
    session.ChannelMessageSend(channel.ID, "Registered!")
    return
  }

  for _, info := range( *users ) {
    msg := fmt.Sprintf("Found Bungie.Net Membership ID: %d", info.MembershipID)
    session.ChannelMessageSend(channel.ID, msg)
  }
}

func handleRegisterShow(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil {
    fmt.Println("Could not create channel for DM!")
    return
  }

  var rego *Registration
  rego, err = loadRego(message.Author.ID)
  if err != nil {
    msg := fmt.Sprintf("Could not query database! %s", err)
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  if rego == nil {
    msg := fmt.Sprintf("No Registration Found")
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  msg := fmt.Sprintf("Your Bungie.Net ID is %d on %s", rego.bungieID, rego.network)
  session.ChannelMessageSend(channel.ID, msg)
}
