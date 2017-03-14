package main
import (
  "log"
  "net/http"
  "os"

  "github.com/el-komandante/gochat/routes"
  "github.com/gorilla/handlers"
)
func main() {
  r := routes.GetRouter()
  // r := mux.NewRouter()
  // r.Handle("/users", jwtMiddleware.Handler(CreateUserHandler)).Methods("POST")
  // // r.HandleFunc("/login", loginHandler)
  // r.HandleFunc("/ws", connectionHandler)
  // r.HandleFunc("/token", GetTokenHandler)
  // r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
  // http.Handle("/", r)
  // go handleMessages()

  log.Println("https server started on :8000")
  // err := http.ListenAndServe(":8000", nil)
  log.Fatal(http.ListenAndServeTLS(":8000", "server.crt", "server.key", handlers.LoggingHandler(os.Stdout, r)))
}
