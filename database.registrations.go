package main

import (
  "github.com/lordmortis/goBungieNet"
)

type Registration struct {
  newRecord bool
  discordID string
  bungieID int64
  network goBungieNet.BungieMembershipType
}

func (rego *Registration)GetProfile(components []goBungieNet.DestinyComponentType) (*goBungieNet.GetProfileResponse, error) {
	return goBungieNet.GetProfile(rego.bungieID, rego.network, components)
}

func migrateRegoTable() error {
  createSQL := `CREATE TABLE IF NOT EXISTS 'registrations' (
    'discord_id' TEXT,
    'bungie_id' INTEGER,
    'network_type' INTEGER
  )`

  // todo - handle if table exists
  _, err := db.Exec(createSQL)

  return err
}

func (rego *Registration)Save() error {
  sql := ""
  if rego.newRecord {
    sql = `INSERT INTO registrations
      (bungie_id, network_type, discord_id)
      VALUES (?, ?, ?)`
  } else {
    sql = `UPDATE registrations
      SET bungie_id = ?, network_type = ?
      WHERE discord_id = ?`
  }

  _, err := db.Exec(sql, rego.bungieID, rego.network, rego.discordID)

  return err
}

func loadRego(discordID string) (*Registration, error){
  sql := `SELECT bungie_id, network_type FROM registrations WHERE discord_id = $1`
  rows, err := db.Query(sql, discordID)
  if err != nil { return nil, err }
  if !rows.Next() { return nil, nil }
  defer rows.Close()

  rego := Registration{ newRecord: false, discordID: discordID, }
  err = rows.Scan(&rego.bungieID, &rego.network)
  return &rego, err
}

func loadRegos() (*[]Registration, error) {
  sql := `SELECT discord_id, bungie_id, network_type FROM registrations`
  rows, err := db.Query(sql)
  if err != nil { return nil, err }
  defer rows.Close()
  var regos []Registration
  for rows.Next() {
    rego := Registration{newRecord: false}
    err = rows.Scan(&rego.discordID, &rego.bungieID, &rego.network)
    if err != nil { return &regos, err }
    regos = append(regos, rego)
  }

  return &regos, nil
}

func createRego(discordID string, bungieID int64, network goBungieNet.BungieMembershipType) error {
  rego, err := loadRego(discordID)
  if err != nil { return err }

  if rego == nil {
    rego = &Registration{
      newRecord: true,
      discordID: discordID,
      bungieID: bungieID,
      network: network,
    }
  } else {
    rego.bungieID = bungieID
    rego.network = network
  }

  err = rego.Save()
  if err == nil { rego.newRecord = false }
  return err
}
