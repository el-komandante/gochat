package models

import (
    "errors"

    "golang.org/x/crypto/bcrypt"
    "github.com/jinzhu/gorm"
)
type User struct {
    gorm.Model
    Email    string   `json:"email"`
    Username string   `json:"username"`
    Password string   `json:"password"`
    Session  Session
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

func (u *User) GetSession() (Session, error) {
    var (
        sess Session
    )
    if DB.Model(u).Related(&sess).RecordNotFound() {
        return sess, errors.New("No sessions found for this user.")
    }
    return sess, nil
}

func (u *User) NewSession() (Session, error) {
    var sess Session

    if u.ID == 0 {
        return sess, errors.New("Missing user ID")
    }
    sess = CreateSession(u.ID)
    return sess, nil
}
