//db.go

package db

import (
    "database/sql"
    "log"
    "os"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX idx_date ON scheduler(date);
`

func Init(dbFile string) error {
    _, err := os.Stat(dbFile)
    install := false
    if err != nil {
        install = true
    }

    var errOpen error
    DB, errOpen = sql.Open("sqlite", dbFile)
    if errOpen != nil {
        return errOpen
    }

    if install {
        _, err = DB.Exec(schema)
        if err != nil {
            return err
        }
        log.Println("База данных создана и таблица scheduler добавлена")
    }

    return nil
}
