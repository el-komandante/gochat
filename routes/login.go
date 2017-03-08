package routes

import (
    "net/http"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    "github.com/el-komandante/gochat/models"
    "github.com/gorilla/mux"
    "log"
    "time"
)

// type Credentials struct{
//   Username string `json:username`
//   Password string `json:password`
// }

func addLoginRoutes(r *mux.Router) *mux.Router {
    r.HandleFunc("/login", loginHandler)/*.Methods("OPTIONS, POST")*/
    return r
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
    var (
        user models.User
        // storedUser models.User
        sess models.Session
        err error
    )



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
        http.Error(w, err.Error(), http.StatusNotFound)
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

    if models.DB.Where("UserID = ?", user.ID).First(&sess).RecordNotFound() {
      sess = models.CreateSession(user.ID)
    }
    sess = models.CreateSession(user.ID)

    maxAge := int(time.Now().Add(time.Duration(240) * time.Hour).Unix())

    cookie := http.Cookie{
      Name: "sessionid",
      Value: sess.SessionID,
      Path: "/",
      MaxAge: maxAge,
      Secure: true,
      HttpOnly: true,
    }
    log.Printf("%v", cookie.String())
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    http.SetCookie(w, &cookie)

  }



  // w.Header().Set("Content-Type", "application/json")
  // w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

//   if req.Method == "OPTIONS" {
//       return
//   }
// }
