package models

import (
  "time"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
  time.Sleep(time.Duration(20) * time.Second)
  db, err := gorm.Open("postgres", "user=postgres password=postgres host=localhost port=5432 dbname=gochat sslmode=disable")
  if err != nil {
    panic(err)
  }
  db.AutoMigrate(&User{})
}
