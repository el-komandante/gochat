package routes
//
// import (
//     "github.com/el-komandante/gochat/models"
//     "github.com/gorilla/mux"
//     "errors"
// )
//
// func authenticate(h http.HandlerFunc, w http.ResponseWriter, req *http.Request) {
//     var sess Session
//     if models.DB.Where("SessionID = ?", req.Cookie("session")).First(&sess).RecordNotFound() {
//         http.Error(w, errors.New("Session not found"), http.StatusUnauthorized)
//         return
//     }
//
// }
