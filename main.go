package main

import (
  "fmt"
  "runtime"
  "flag"
  "os"
  "os/signal"
  "syscall"

  "github.com/lordmortis/DiscordDestinyInfo/discord"
  "github.com/lordmortis/goBungieNet"
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

  goBungieNet.ApiKey = config.BungieNet.ApiKey
  err = goBungieNet.SetManifestPath(config.BungieNet.ManifestPath)
  if (err != nil) {
    fmt.Println("Destiny Manifest Directory error:")
    fmt.Println(err)
    return
  }

  err = goBungieNet.ManifestUpdate()
  if (err != nil) {
    fmt.Println("Manifest did not update")
    fmt.Println(err)
    return
  }

  //components := []goBungieNet.DestinyComponentType{goBungieNet.ComponentCharacters,goBungieNet.ComponentCharacterActivities}
  //response, err1 := goBungieNet.GetProfile(id, goBungieNet.NetworkPsn, components)
  /*profiles, err1 := goBungieNet.FindDestinyUser("maeglinhiei", goBungieNet.NetworkPsn)*/
  //if (err1 != nil) {
//    fmt.Println(err1.Error())
//    return
  //}

  /*for _, profile := range( *profiles ) {
    fmt.Printf("ID: %d\n", profile.MembershipID)
  }*/


//  currentCharID := response.CharacterActivities.MostRecentCharacter()
//  fmt.Printf("Character ID: %d\n", id)
//  fmt.Printf("\tPrivacy:%d\n", response.CharacterActivities.Privacy)
//  fmt.Printf("\tStarted: %s\n", response.CharacterActivities.Data[currentCharID].DateActivityStarted)
//  fmt.Printf("\tActivityHash: %d\n", response.CharacterActivities.Data[currentCharID].CurrentActivityHash)
//  fmt.Printf("\tActivity: %d\n", response.CharacterActivities.Data[currentCharID].CurrentActivityModeType)

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
