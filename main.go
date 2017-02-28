package main
import (
  "log"
  "net/http"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/el-komandante/gochat/models"
  "github.com/gorilla/websocket"
)
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
        return true
      },
}

type Message struct {
  Email string    `json:email`
  Username string `json:username`
  Message string  `json:message`
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer ws.Close()
  clients[ws] = true

  for {
    var msg Message
    err := ws.ReadJSON(&msg)
    if err != nil {
      log.Printf("error: %v", err)
      delete(clients, ws)
      break
    }
    broadcast <- msg
  }
}
//
// func handleMessages() {
//   for {
//     msg := <- broadcast
//
//     for client:= range clients {
//       err := client.WriteJSON(msg)
//       if err != nil {
//         log.Printf("error: %v", err)
//         client.Close()
//         delete(clients, client)
//       }
//     }
//   }
// }

func createUserHandler(w http.ResponseWriter, req *http.Request) {
  var user models.User

  decoder := json.NewDecoder(req.Body)
  defer req.Body.Close()

  err := decoder.Decode(&user)
  if err != nil {
    panic(err)
  }

  user.Password = models.CreatePassword(user.Password)
  models.DB.Create(&user)
  log.Println(user)
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/users", createUserHandler)
  // r.HandleFunc("/login", loginHandler)
  r.HandleFunc("/ws", connectionHandler)
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("../public")))
  http.Handle("/", r)
  // go handleMessages()

  log.Println("http server started on :8000")
  // err := http.ListenAndServe(":8000", nil)
  log.Fatal(http.ListenAndServe(":8000", r))
}
