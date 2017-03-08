package routes

import (
  "github.com/gorilla/mux"
)


func GetRouter() *mux.Router {
  r := mux.NewRouter()
  r = addUserRoutes(r)
  r = addLoginRoutes(r)
  // r.Handle("/users", jwtMiddleware.Handler(CreateUserHandler)).Methods("POST")
  // r.HandleFunc("/login", loginHandler)
  // r.HandleFunc("/ws", connectionHandler)
  // r.HandleFunc("/token", GetTokenHandler)
  // r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
  return r
}
