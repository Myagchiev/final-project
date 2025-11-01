// pkg/db/task.go
package db

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    "github.com/Myagchiev/final-project/pkg/utils"
)

type Task struct {
    ID      int    `json:"id,string"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func AddTask(task Task) (int, error) {
    if DB == nil {
        return 0, sql.ErrConnDone
    }

    res, err := DB.Exec(
        "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
        task.Date, task.Title, task.Comment, task.Repeat,
    )
    if err != nil {
        return 0, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func Tasks(limit int) ([]Task, error) {
    return TasksWithFilter(limit, "", "")
}

func TasksWithFilter(limit int, searchText, searchDate string) ([]Task, error) {
    var args []interface{}
    var where []string
    query := "SELECT id, date, title, comment, repeat FROM scheduler"

    if searchDate != "" {
        where = append(where, "date = ?")
        args = append(args, searchDate)
    } else if searchText != "" {
        like := "%" + searchText + "%"
        where = append(where, "(title LIKE ? OR comment LIKE ?)")
        args = append(args, like, like)
    }

    if len(where) > 0 {
        query += " WHERE " + strings.Join(where, " AND ")
    }

    query += " ORDER BY date ASC"
    if limit > 0 {
        query += " LIMIT ?"
        args = append(args, limit)
    }

    rows, err := DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var t Task
        var id int64
        if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
            return nil, err
        }
        t.ID = int(id)
        tasks = append(tasks, t)
    }

    return tasks, rows.Err()
}

func GetTask(id int) (Task, error) {
    var t Task
    err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).
        Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
    if err != nil {
        if err == sql.ErrNoRows {
            return t, fmt.Errorf("task not found")
        }
        return t, err
    }
    return t, nil
}

func UpdateTask(task Task) error {
    res, err := DB.Exec(`
        UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?`,
        task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    return nil
}

func DeleteTask(id int) error {
    res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    return nil
}

func MarkDone(id int) error {
    task, err := GetTask(id)
    if err != nil {
        return err
    }

    if task.Repeat == "" {
        return DeleteTask(id)
    }

    nextDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
    if err != nil {
        return err
    }

    res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, id)
    if err != nil {
        return err
    }

    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("task not found")
    }
    return nil
}