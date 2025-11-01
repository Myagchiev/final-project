//server.go

package server

import (
    "log"
    "net/http"
    "os"

    "github.com/Myagchiev/final-project/pkg/api"
)

const defaultPort = "7540"

func Run() {
    webDir := "./web"

    port := os.Getenv("TODO_PORT")
    if port == "" {
        port = defaultPort
    }

    api.Init()

    http.Handle("/", http.FileServer(http.Dir(webDir)))

    log.Printf("Сервер запущен на порту %s...\n", port)

    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}
