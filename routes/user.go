package routes

import (
  "net/http"
  "log"
  "encoding/json"
  "strconv"
  // "net/url"
  // "bytes"

  "github.com/el-komandante/gochat/models"
  "github.com/gorilla/mux"
)

func addUserRoutes(r *mux.Router) *mux.Router {
  r.HandleFunc("/users/messages", getMessagesHandler)
  r.HandleFunc("/users", createUserHandler)/*.Methods("POST, GET, OPTIONS")*/
  return r
}
func createUserHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    if req.Method == "POST" {
        var user models.User
        var currentUser models.User

        err := json.NewDecoder(req.Body).Decode(&user)
        defer req.Body.Close()
        if err != nil {
            panic(err)
        }
        if models.DB.Where("Username = ?", user.Username).First(&currentUser).RecordNotFound() {
            user.Password = models.CreatePassword(user.Password)
            models.DB.Create(&user)
            err = json.NewEncoder(w).Encode(&user)
            if err != nil {
                panic(err)
            }
            log.Println(user)
        } else {
            log.Printf("User %v already exists", user.Username)
        }
    } else if req.Method == "GET" {
        var users []models.User

        models.DB.Find(&users)
        err := json.NewEncoder(w).Encode(&users)
        if err != nil {
            panic(err)
        }
    }

}

func getUsersHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    if req.Method == "GET" {
        var users []models.User

        models.DB.Find(&users)
        err := json.NewEncoder(w).Encode(&users)
        if err != nil {
            panic(err)
        }
    }
}

func getMessagesHandler(w http.ResponseWriter, req *http.Request) {
    var messages []models.Message
    var u models.User
    // var id []string
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    if req.Method == "GET" {
        // query, err := url.ParseQuery(req.URL.RawQuery)

        // if id, ok := query["id"]; !ok {
        //     http.Error(w, "User not found.", http.StatusNotFound)
        //     return
        // }
        query := req.URL.Query().Get("id")
        log.Printf("query value: %v, query type: %t", query, query)
        id, err := strconv.ParseUint(query, 10, 64)
        if err != nil {
            panic(err)
        }
        // id, err := strconv.Atoi(query)
        // if err != nil {
        //     panic(err)
        // }

        cookie, err := req.Cookie("session")
        if err != nil {
            http.Error(w, "You need to be logged in to do that.", http.StatusUnauthorized)
            return
        }

        u, err = u.FromSessionID(cookie.Value)
        if err != nil {
            panic(err)
        }
        log.Printf("Finding messages between users %v and %v.", id, u.ID)

        models.DB.Where("\"from\" = ? AND \"to\" = ?", 7, 6).Or("\"from\" = ? AND \"to\" = ?", 6, 7).Order("created_at desc").Limit(10)/*.Offset(0)*/.Find(&messages)
        err = json.NewEncoder(w).Encode(&messages)
        if err != nil {
            panic(err)
        }
    }
}
