package models

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "time"
)

var DB *gorm.DB

func init() {
  time.Sleep(time.Duration(10) * time.Second)
  var err error
  DB, err = gorm.Open("postgres", "postgres://postgres:postgres@localhost:5432/gochat?sslmode=disable")
  if err != nil {
    panic(err)
  }
  DB.AutoMigrate(&User{})
}
