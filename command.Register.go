package main

import (
  "fmt"
  "strings"
  "strconv"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
  "github.com/lordmortis/goBungieNet"

  "github.com/bwmarrin/discordgo"
)

func init() {
  discord.RegisterCommand("RegisterSearch", "Search for your PSN/Xbox account (if only one result matches, will register) - `RegisterSearch <GamerTag/Nickname> <Xbox/Psn/BattleNet>`", handleRegisterSearch)
  discord.RegisterCommand("Register", "Register your bungo account - `Register <Bungie.Net Membership ID> <Xbox/Psn/BattleNet>`", handleRegister)
  discord.RegisterCommand("ShowRegistration", "Show your registration if it exists", handleRegisterShow)
}

func handleRegister(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil { discord.LogPMCreateError(message.Author); return; }

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

  discord.LogChatCommand(message.Author, "Register %d %s", id, accountType)

  components := []goBungieNet.DestinyComponentType{goBungieNet.ComponentCharacters}
  _, err = goBungieNet.GetProfile(id, accountType, components)
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Couldn't get profile: %s", err.Error())
    return
  }

  err = createRego(message.Author.ID, id, accountType)
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Couldn't save registration: %s", err.Error())
    return
  }

  discord.LogChatMsg(message.Author, "Registered %d %s", id, accountType)
}

func handleRegisterSearch(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil { discord.LogPMCreateError(message.Author); return; }

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

  discord.LogChatCommand(message.Author, "RegisterSearch %s %s", paramParts[0], accountType)

  users, err := goBungieNet.FindDestinyUser(paramParts[0], accountType)
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Error from Bungie: %s", err.Error())
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
      discord.LogPMError(session, message.Author, channel, "Couldn't save registration: %s", err.Error())
      return
    }
    session.ChannelMessageSend(channel.ID, "Registered!")
    discord.LogChatMsg(message.Author, "Registered %d %s", user.MembershipID, user.MembershipType)
    return
  }

  session.ChannelMessageSend(channel.ID, "We found multiple matching bungie memberships - use Register with the correct id")
  for _, info := range( *users ) {
    msg := fmt.Sprintf("Found Bungie.Net Membership ID: %d on network %s - use `Register %d %s` to register this.", info.MembershipID, info.MembershipType, info.MembershipID, info.MembershipType)
    session.ChannelMessageSend(channel.ID, msg)
  }
}

func handleRegisterShow(session *discordgo.Session, message *discordgo.Message, parameters string) {
  channel, err := session.UserChannelCreate(message.Author.ID)
  if err != nil { discord.LogPMCreateError(message.Author); return; }

  discord.LogChatCommand(message.Author, "ShowRegistration")

  var rego *Registration
  rego, err = loadRego(message.Author.ID)
  if err != nil {
    discord.LogPMError(session, message.Author, channel, "Couldn't query database: %s", err.Error())
    return
  }

  if rego == nil {
    discord.LogChatMsg(message.Author, "No registration found")
    msg := fmt.Sprintf("No Registration Found")
    session.ChannelMessageSend(channel.ID, msg)
    return
  }

  discord.LogChatMsg(message.Author, "found registration: %d on %s", rego.bungieID, rego.network)
  msg := fmt.Sprintf("Your Bungie.Net ID is %d on %s", rego.bungieID, rego.network)
  session.ChannelMessageSend(channel.ID, msg)
}
