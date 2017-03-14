package routes

import (
  "github.com/gorilla/mux"
  "net/http"
)


func GetRouter() *mux.Router {
  r := mux.NewRouter()
  r = addUserRoutes(r)
  r = addLoginRoutes(r)
  r = addLogoutRoutes(r)
  r = addWsRoutes(r)
  // r.Handle("/users", jwtMiddleware.Handler(CreateUserHandler)).Methods("POST")
  // r.HandleFunc("/login", loginHandler)
  // r.HandleFunc("/ws", connectionHandler)
  // r.HandleFunc("/token", GetTokenHandler)
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
  return r
}

func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
    for _, m := range middleware {
        h = m(h)
    }
    return h
}
