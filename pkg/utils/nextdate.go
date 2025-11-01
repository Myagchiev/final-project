// pkg/utils/nextdate.go
package utils

import (
    "errors"
    "strconv"
    "strings"
    "time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
    date, err := time.Parse(DateLayout, dstart)
    if err != nil {
        return "", errors.New("неверная дата dstart")
    }

    if repeat == "" {
        return "", errors.New("пустое правило повторения")
    }

    parts := strings.Split(strings.TrimSpace(repeat), " ")
    if len(parts) == 0 {
        return "", errors.New("некорректное правило повторения")
    }

    rule := parts[0]

    switch rule {
    case "d":
        if len(parts) < 2 {
            return "", errors.New("не указан интервал в днях")
        }
        days, err := strconv.Atoi(parts[1])
        if err != nil || days <= 0 || days > 400 {
            return "", errors.New("недопустимый интервал дней")
        }
        for {
            date = date.AddDate(0, 0, days)
            if date.After(now) {
                return date.Format(DateLayout), nil
            }
        }

    case "y":
        for {
            date = date.AddDate(1, 0, 0)
            if date.After(now) {
                return date.Format(DateLayout), nil
            }
        }

    case "w":
        if len(parts) < 2 {
            return "", errors.New("не указаны дни недели")
        }
        targetDays, err := parseDays(parts[1])
        if err != nil {
            return "", err
        }
        for {
            date = date.AddDate(0, 0, 1)
            wd := int(date.Weekday())
            if wd == 0 {
                wd = 7
            }
            if targetDays[wd] && date.After(now) {
                return date.Format(DateLayout), nil
            }
        }

    case "m":
        if len(parts) < 2 {
            return "", errors.New("не указаны дни месяца")
        }
        dayStr := parts[1]
        monthStr := ""
        if len(parts) >= 3 {
            monthStr = parts[2]
        }
        targetDays, err := parseMonthDays(dayStr)
        if err != nil {
            return "", err
        }
        targetMonths, err := parseMonths(monthStr)
        if err != nil {
            return "", err
        }
        for {
            date = date.AddDate(0, 0, 1)
            _, m, d := date.Year(), date.Month(), date.Day()
            if !targetMonths[int(m)] {
                continue
            }
            matched := false
            if targetDays[d] {
                matched = true
            }
            if targetDays[-1] && d == getLastDay(date) {
                matched = true
            }
            if targetDays[-2] && d == getSecondLastDay(date) {
                matched = true
            }
            if matched && date.After(now) {
                return date.Format(DateLayout), nil
            }
        }

    default:
        return "", errors.New("неподдерживаемый формат правила")
    }
}

func getLastDay(t time.Time) int {
    return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func getSecondLastDay(t time.Time) int {
    last := getLastDay(t)
    if last >= 2 {
        return last - 1
    }
    return last
}

func parseDays(s string) (map[int]bool, error) {
    days := make(map[int]bool)
    for _, p := range strings.Split(s, ",") {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        d, err := strconv.Atoi(p)
        if err != nil || d < 1 || d > 7 {
            return nil, errors.New("недопустимый день недели")
        }
        days[d] = true
    }
    if len(days) == 0 {
        return nil, errors.New("не указаны дни недели")
    }
    return days, nil
}

func parseMonthDays(s string) (map[int]bool, error) {
    days := make(map[int]bool)
    for _, p := range strings.Split(s, ",") {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        if p == "-1" {
            days[-1] = true
            continue
        }
        if p == "-2" {
            days[-2] = true
            continue
        }
        d, err := strconv.Atoi(p)
        if err != nil || d < 1 || d > 31 {
            return nil, errors.New("недопустимый день месяца")
        }
        days[d] = true
    }
    if len(days) == 0 {
        return nil, errors.New("не указаны дни месяца")
    }
    return days, nil
}

func parseMonths(s string) (map[int]bool, error) {
    months := make(map[int]bool)
    if s == "" {
        for i := 1; i <= 12; i++ {
            months[i] = true
        }
        return months, nil
    }
    for _, p := range strings.Split(s, ",") {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        m, err := strconv.Atoi(p)
        if err != nil || m < 1 || m > 12 {
            return nil, errors.New("недопустимый месяц")
        }
        months[m] = true
    }
    if len(months) == 0 {
        return nil, errors.New("не указаны месяцы")
    }
    return months, nil
}