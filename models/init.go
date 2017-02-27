package models

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
  db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=gochat sslmode=disable password=postgres")
  if err != nil {
    panic(err)
  }
  db.AutoMigrate(&User{})
}
