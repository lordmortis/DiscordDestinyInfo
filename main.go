package main

import (
  "fmt"
  "runtime"
  "flag"
  "os"
  "os/signal"
  "syscall"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
)

var (
  config *Config
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())
  configFile := flag.String("config", "config.json", "JSON Config File")

  flag.Parse()

  var err error

  config, err = LoadConfig(*configFile)

  if (err != nil) {
    fmt.Println("Unable to parse/load config!")
    fmt.Println(err)
    return
  }

  err = discord.Setup(config.Discord.Token)
  if (err != nil) {
    fmt.Println("Unable to connect to discord, ", err)
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
