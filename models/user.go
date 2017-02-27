package models

import (
  "golang.org/x/crypto/bcrypt"
  "github.com/jinzhu/gorm"
)
type User struct {
  gorm.Model
  Email string    `json:"email"`
  Username string `json:"username"`
  Password string `json:"password"`
}

func CreateUser(user User) {
  // db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=gochat sslmode=disable password=postgres")
  // if err != nil {
  //   panic(err)
  // }
  user.Password = string(CreatePassword(user.Password))

  DB.NewRecord(user)
}

func CreatePassword(pw string) string {
  password := []byte(pw)
  hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }
  return string(hash)
}
