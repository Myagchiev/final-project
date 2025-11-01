package main

import (
    "log"
    "os"

    "github.com/Myagchiev/final-project/pkg/db"
    "github.com/Myagchiev/final-project/pkg/server"
)

func main() {
    dbFile := os.Getenv("TODO_DBFILE")
    if dbFile == "" {
        dbFile = "scheduler.db"
    }

    err := db.Init(dbFile)
    if err != nil {
        log.Fatalf("Ошибка инициализации БД: %v", err)
    }

    server.Run()
}
