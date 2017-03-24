package models

import (
    // "log"
    // "errors"

    "github.com/jinzhu/gorm"
)

type Message struct {
    gorm.Model
    From uint   `json: from`
    To uint
    Text string
}
