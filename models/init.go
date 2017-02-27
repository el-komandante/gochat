package models

import (
  "time"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
  time.Sleep(time.Duration(20) * time.Second)
  db, err := gorm.Open("postgres", "postgres://postgres:postgres@gochat?sslmode=disable")
  if err != nil {
    panic(err)
  }
  db.AutoMigrate(&User{})
}
