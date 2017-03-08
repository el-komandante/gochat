package models

import (
    "log"
    "crypto/rand"
    "encoding/base64"
    "time"

    "github.com/jinzhu/gorm"
)

type Session struct {
    gorm.Model
    Expires int
    SessionID string
    UserID uint
}
func CreateSession(userID uint) Session {
    var sess Session

    id, err := GenerateRandomString(27)

    if err != nil {
        log.Fatal(err)
    }

    sess.SessionID = id
    sess.UserID = userID
    sess.Expires = int(time.Now().Add(time.Duration(240) * time.Hour).Unix())

    DB.NewRecord(sess)

    return sess
}

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)

    _, err := rand.Read(b)

    if err != nil {
        return nil, err
    }
    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}
