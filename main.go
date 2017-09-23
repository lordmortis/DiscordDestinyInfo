package main

import (
  "fmt"
  "runtime"
  "flag"
  "os"
  "os/signal"
  "syscall"

  "github.com/bwmarrin/discordgo"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  configFile := flag.String("config", "config.json", "JSON Config File")

  flag.Parse()

  var config, err1 = LoadConfig(*configFile)

  if (err1 != nil) {
    fmt.Println("Unable to parse/load config!")
    fmt.Println(err1)
    return
  }

  discord, err := discordgo.New("Bot " + config.Discord.Token)
  if err != nil {
    fmt.Println("error creating Discord session,", err)
    return
  }

  // Register the messageCreate func as a callback for MessageCreate events.
  discord.AddHandler(messageCreate)

  // Open a websocket connection to Discord and begin listening.
  err = discord.Open()
  if err != nil {
    fmt.Println("error opening connection,", err)
    return
  }

  // Wait here until CTRL-C or other term signal is received.
  fmt.Println("Bot is now running.  Press CTRL-C to exit.")
  fmt.Println("To add this bot to your server, visit https://discordapp.com/oauth2/authorize?scope=bot&permissions=o&client_id=" + config.Discord.ClientID)
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  // Cleanly close down the Discord session.
  discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

  // Ignore all messages created by the bot itself
  // This isn't required in this specific example but it's a good practice.
  if m.Author.ID == s.State.User.ID {
    return
  }
  // If the message is "ping" reply with "Pong!"
  if m.Content == "ping" {
    s.ChannelMessageSend(m.ChannelID, "Pong!")
  }

  // If the message is "pong" reply with "Ping!"
  if m.Content == "pong" {
    s.ChannelMessageSend(m.ChannelID, "Ping!")
  }
}
