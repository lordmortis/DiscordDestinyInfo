package main

import (
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/lordmortis/goBungieNet"

  //"github.com/davecgh/go-spew/spew"
)

var (
  db *sql.DB
)

type Registration struct {
  newRecord bool
  discordID string
  bungieID int64
  network goBungieNet.BungieMembershipType
}

func setDbPath(path string) error {
  adb, err := sql.Open("sqlite3", path)
  if err != nil { return err }

  db = adb
  err = migrateRegoTable()
  if err != nil { return err }

  return nil
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
    sql = `UPDATE registrations
      SET bungie_id = $1, network_type = $2
      WHERE discord_id = $3`
  } else {
    sql = `INSERT INTO registrations
      (bungie_id, network_type, discord_id)
      VALUES ($1, $2, $3)`
  }

  _, err := db.Exec(sql, rego.bungieID, rego.network, rego.discordID)
  return err
}

func loadRego(discordID string) (*Registration, error){
  sql := `SELECT bungie_id, network_type FROM registrations WHERE discord_id = $1`
  rows, err := db.Query(sql, discordID)
  if err != nil { return nil, err }
  if !rows.Next() { return nil, nil }

  rego := Registration{ newRecord: false, }
  err = rows.Scan(&rego.bungieID, &rego.network)

  return &rego, err
}

func createRego(discordID string, bungieID int64, network goBungieNet.BungieMembershipType) error {
  rego := Registration{
    discordID: discordID,
    bungieID: bungieID,
    network: network,
  }

  err := rego.Save()
  if err == nil { rego.newRecord = false }
  return err
}


