// pkg/api/api.go
package api

import (
    "fmt"
    "net/http"
    "time"

    "github.com/Myagchiev/final-project/pkg/utils"
)

func Init() {
    http.HandleFunc("/api/nextdate", NextDateHandler)
    http.HandleFunc("/api/signin", SignInHandler)
    http.HandleFunc("/api/task", Auth(taskCRUDHandler))
    http.HandleFunc("/api/tasks", Auth(tasksListHandler))
    http.HandleFunc("/api/task/done", Auth(taskCRUDHandler))
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
    nowStr := r.FormValue("now")
    dateStr := r.FormValue("date")
    repeat := r.FormValue("repeat")

    now := time.Now()
    if nowStr != "" {
        if parsed, err := time.Parse(utils.DateLayout, nowStr); err == nil {
            now = parsed
        }
    }

    next, err := utils.NextDate(now, dateStr, repeat)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Fprint(w, next)
}