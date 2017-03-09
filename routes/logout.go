package routes

import (
    "net/http"
    "log"

    "github.com/el-komandante/gochat/models"
    "github.com/gorilla/mux"
)

func addLogoutRoutes(r *mux.Router) *mux.Router {
    r.HandleFunc("/logout", logoutHandler)
    return r
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
    log.Printf("req: %v", req)
    cookie, err := req.Cookie("session")
    sess := models.Session{SessionID: cookie.Value}
    if err != nil {
        log.Printf("%v", err)
        return
    }
    if models.DB.Where("session_id = ?", sess.SessionID).First(&sess).RecordNotFound() {
        return
    }
    models.DB.Unscoped().Delete(&sess)
    log.Printf("session deleted")
}
