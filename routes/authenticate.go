package routes

import (
    "net/http"
    "log"

    "github.com/el-komandante/gochat/models"
)

func authenticate(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        var sess models.Session
        cookie, err := r.Cookie("session")
        if err != nil {
            log.Printf("no cookie found: %v", err)
            http.Error(w, "You need to be logged in to do that", http.StatusUnauthorized)
            return
        }
        if models.DB.Where("session_id = ?", cookie.Value).First(&sess).RecordNotFound() {
            log.Printf("session not found in db")
            http.Error(w, "You need to be logged in to do that", http.StatusUnauthorized)
            return
        }
        h.ServeHTTP(w, r)
    })

}
