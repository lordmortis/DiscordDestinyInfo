package main

import (
	"fmt"

	"github.com/lordmortis/DiscordDestinyInfo/discord"
	"github.com/lordmortis/goBungieNet"

	"github.com/bwmarrin/discordgo"
)

func init() {
	discord.RegisterCommand("WhosOn", "Tell me who is on and what they are doing", handleWhosOn)
}

func handleWhosOn(session *discordgo.Session, message *discordgo.Message, parameters string) {
	channel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		discord.LogPMCreateError(message.Author)
		return
	}

	discord.LogChatCommand(message.Author, "WhosOn")
	session.ChannelMessageSend(message.ChannelID, "Checking Who is online...")

	var regos *[]Registration
	regos, err = loadRegos()
	if err != nil {
		discord.LogPMError(session, message.Author, channel, "Couldn't get registrations: %s", err.Error())
		return
	}

	msg := "I think the following players are online:\n"

	components := []goBungieNet.DestinyComponentType{
		goBungieNet.ComponentCharacters,
		goBungieNet.ComponentCharacterActivities,
	}

	onlineCount := 0

	for _, rego := range *regos {
		var response *goBungieNet.GetProfileResponse
		response, err = rego.GetProfile(components)
		if err != nil {
			discord.LogPMError(session, message.Author, channel, "Couldn't get profile info for %d because: %s", rego.bungieID, err.Error())
			continue
		}

		charID := response.CharacterActivities.MostRecentCharacterID()
		currentCharacter := response.Characters.Data[charID]
		currentActivity := response.CharacterActivities.Data[charID]

		// if the hash is 0 they aren't playing at the moment.
		if currentActivity.CurrentActivityHash == 0 {
			continue
		}

		var currentActivityData *goBungieNet.DestinyActivity
		currentActivityData, err = currentActivity.ActivityDefinition("en")
		if err != nil {
			discord.LogPMError(
				session,
				message.Author,
				channel,
				"Couldn't get activity details for %d because: %s",
				currentActivity.CurrentActivityModeHash,
				rego.bungieID,
				err.Error(),
			)
			continue
		}
		activityName := currentActivityData.DisplayProperties.Name

		activityModeName := ""
		var currentActivityModeData *goBungieNet.DestinyActivityModeDefinition

		if currentActivity.CurrentActivityModeType != goBungieNet.DestinyActivityModeNone {
			currentActivityModeData, err = currentActivity.ActivityModeDefinition("en")
			if err != nil {
				discord.LogPMError(
					session,
					message.Author,
					channel,
					"Couldn't get activity mode details %d for %d because: %s",
					currentActivity.CurrentActivityModeHash,
					rego.bungieID,
					err.Error(),
				)
				continue
			}
			activityModeName = currentActivityModeData.DisplayProperties.Name
		}

		var class *goBungieNet.DestinyClassDefinition
		class, err = currentCharacter.Class("en")
		if err != nil {
			discord.LogPMError(session, message.Author, channel, "Couldn't get class details for %d because: %s", rego.bungieID, err.Error())
			continue
		}

		levelString := ""
		if currentCharacter.LevelProgression.Level == currentCharacter.LevelProgression.LevelCap {
			levelString = fmt.Sprintf("%d Light", currentCharacter.Light)
		} else {
			levelString = fmt.Sprintf("Level %d", currentCharacter.LevelProgression.Level)
		}

		msgString := "<@%s> playing their %s %s on %s"
		msg += fmt.Sprintf(msgString,
			rego.discordID,
			levelString,
			class.DisplayProperties.Name,
			currentCharacter.MembershipType,
		)

		switch currentActivity.CurrentActivityModeType {
		case goBungieNet.DestinyActivityModeNone:
			msg += fmt.Sprintf(" and they're In Orbit")
		case goBungieNet.DestinyActivityModePatrol:
			msg += fmt.Sprintf(" doing %s %s", activityName, activityModeName)
		case goBungieNet.DestinyActivityModeStory:
			msg += fmt.Sprintf(" doing story mission: %s", activityModeName)
		case goBungieNet.DestinyActivityModeSocial:
			msg += fmt.Sprintf(" and they're at the %s", activityName)
		case goBungieNet.DestinyActivityModeClash:
		case goBungieNet.DestinyActivityModeControl:
		case goBungieNet.DestinyActivityModeSupremacy:
		case goBungieNet.DestinyActivityModeCountdown:
		case goBungieNet.DestinyActivityModeSurvival:
			msg += fmt.Sprintf(" is playing %s on the %s map", currentActivityModeData.ModeType, activityName)
		default:
			msg += fmt.Sprintf(
				" doing %s %s %s",
				currentActivityModeData.ModeType,
				activityName,
				activityModeName,
			)
		}

		msg += "\n"
		onlineCount++
	}

	if onlineCount == 0 {
		msg = "No one is online!"
	}

	session.ChannelMessageSend(message.ChannelID, msg)
}
