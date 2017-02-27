package models

import (
  "time"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
  time.Sleep(time.Duration(20) * time.Second)
  DB, err := gorm.Open("postgres", "postgres://postgres:postgres@postgres/postgres?sslmode=disable")
  if err != nil {
    panic(err)
  }
  DB.AutoMigrate(&User{})
}
