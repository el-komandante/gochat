package routes

import (
  "net/http"
  "log"
  "encoding/json"
  "github.com/el-komandante/gochat/models"
  "github.com/gorilla/mux"
)

func addUserRoutes(r *mux.Router) *mux.Router {
  r.HandleFunc("/users", createUserHandler).Methods("POST")
  return r
}
func createUserHandler(w http.ResponseWriter, req *http.Request) {
  var user models.User
  var currentUser models.User

  decoder := json.NewDecoder(req.Body)
  defer req.Body.Close()

  err := decoder.Decode(&user)
  if err != nil {
    panic(err)
  }

  if models.DB.Where("Username = ?", user.Username).First(&currentUser).RecordNotFound() {
    user.Password = models.CreatePassword(user.Password)
    models.DB.Create(&user)
    log.Println(user)
  } else {
    log.Printf("User %v already exists", user.Username)
  }

}
