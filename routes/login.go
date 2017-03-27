
package routes

import (
    "net/http"
    "encoding/json"
    "log"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/el-komandante/gochat/models"
    "github.com/gorilla/mux"
)

func addLoginRoutes(r *mux.Router) *mux.Router {
    r.HandleFunc("/login", loginHandler)
    return r
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
    var (
        user models.User
        sess models.Session
        err error
    )

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    if req.Method == "OPTIONS" {
        return
    }

    decoder := json.NewDecoder(req.Body)
    defer req.Body.Close()
    err = decoder.Decode(&user)

    if err != nil {
        panic(err)
    }

    username := user.Username
    password := user.Password

    if models.DB.Where("Username = ?", username).First(&user).RecordNotFound() {
        log.Printf("User %v not found.", user.Username)
        http.Error(w, "User not found", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    sess, err = user.GetSession()
    if err != nil {
        sess, err = user.NewSession()
        if err != nil {
            panic(err)
        }
    } else if sess.Expires < int(time.Now().Unix()) {
        models.DB.Delete(&sess)
        sess, err = user.NewSession()
        if err != nil {
            panic(err)
        }
    }

    // if models.DB.Where("user_id = ?", user.ID).First(&sess).RecordNotFound() {
    //   sess = models.CreateSession(user.ID)
    // }

    expires := time.Now().Add(time.Duration(24) * time.Hour)
    cookie := http.Cookie{
      Name: "session",
      Value: sess.SessionID,
      Path: "/",
      Expires: expires,
      Secure: false,
      HttpOnly: true,
      Domain: "localhost",
    }
    // log.Printf("%v", cookie.String())

    http.SetCookie(w, &cookie)
    w.WriteHeader(200)

  }
