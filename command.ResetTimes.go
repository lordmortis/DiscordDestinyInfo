package main

import (
  "time"
  "fmt"
//  "strings"
//  "strconv"

  "github.com/lordmortis/DiscordDestinyInfo/discord"

  "github.com/bwmarrin/discordgo"
)

var (
  dayDuration time.Duration
  temp time.Duration
)

func init() {
  dayDuration = 24 * time.Hour

  discord.RegisterCommand("NextDaily", "how long until next daily reset?", handleDailyReset)
  discord.RegisterCommand("NextWeekly", "how long until next weekly reset?", handleWeeklyReset)
}

func friendlyDuration(duration time.Duration) string {
  var timeStrings []string
  durationStrings := []string{ "day", "hour", "minute", "second" }
  durations := []time.Duration{ dayDuration, time.Hour, time.Minute, time.Second }

  lastSegment := 0

  for index, interval := range(durations) {
    if (index > 0) {
      duration = duration - time.Duration(lastSegment) * durations[index - 1]
    }

    segment := int(duration.Truncate(interval) / interval)
    lastSegment = segment
    format := "%d %ss"
    if segment == 0 {
      continue
    } else if segment == 1 {
      format = "%d %s"
    }

    timeStrings = append(timeStrings, fmt.Sprintf(format, segment, durationStrings[index]))
  }

  timeString := ""

  for index, section := range(timeStrings) {
    if (index == 0) {
      timeString += section
    } else if (index + 1 == len(timeStrings)) {
      timeString += fmt.Sprintf(" and %s", section)
    } else {
      timeString += fmt.Sprintf(", %s", section)
    }
  }

  return timeString
}

func handleDailyReset(session *discordgo.Session, message *discordgo.Message, parameters string) {
  discord.LogChatCommand(message.Author, "NextDaily")
  // Resets are at 0900 UTC - when's the next 0900?
  now := time.Now().UTC()
  nextReset := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
  if now.Hour() > 9 { nextReset = nextReset.AddDate(0, 0, 1) }

  response := fmt.Sprintf("the next daily reset is in %s", friendlyDuration(nextReset.Sub(now)))

  session.ChannelMessageSend(message.ChannelID, response)
}

func handleWeeklyReset(session *discordgo.Session, message *discordgo.Message, parameters string) {
  discord.LogChatCommand(message.Author, "NextWeekly")
  // Resets are at 0900 UTC on Tuesdays - when's the next 0900 tuesday?
  now := time.Now().UTC()
  daysToTuesday := int(time.Tuesday - now.Weekday())
  nextReset := time.Date(now.Year(), now.Month(), now.Day() + daysToTuesday, 9, 0, 0, 0, time.UTC)
  if (now.After(nextReset)) { nextReset = nextReset.AddDate(0, 0, 7) }

  response := fmt.Sprintf("the next weekly reset is in %s", friendlyDuration(nextReset.Sub(now)))

  session.ChannelMessageSend(message.ChannelID, response)
}
