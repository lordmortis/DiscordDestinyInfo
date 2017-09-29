package main

import (
  "database/sql"
  
  _ "github.com/mattn/go-sqlite3"
  "github.com/lordmortis/goBungieNet"

//  "github.com/davecgh/go-spew/spew"
)

var (
  db *sql.DB
)

func setDbPath(path string) error {
  adb, err := sql.Open("sqlite3", path)
  if err != nil { return err }

  db = adb
  err = migrateRegoTable()
  if err != nil { return err }

  return nil
}

