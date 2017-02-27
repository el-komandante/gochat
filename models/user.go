package models

import (
  "golang.org/x/crypto/bcrypt"
)
type User struct {
  gorm.model
  Email string    `json:"email"`
  Username string `json:"username"`
  Password string `json:"password"`
}

func createUser(user User) {
  user := User{}
  user.Email = email
  user.Username = username
  user.Password = createPassword(password)

  db.NewRecord(user)
}

func createPassword(pw string) {
  hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }
  return hash
}
