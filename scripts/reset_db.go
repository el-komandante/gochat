package main

import (
    "fmt"
    "github.com/el-komandante/gochat/models"
)

func main() {
    fmt.Println("\nDeleting users...\n")
    models.DB.Unscoped().Model(&models.User{}).Delete(&models.User{})
    fmt.Println("Done\n")
    fmt.Println("Deleting sessions...\n")
    models.DB.Unscoped().Model(&models.Session{}).Delete(&models.Session{})
    fmt.Println("Done\n")
}
