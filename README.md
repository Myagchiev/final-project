# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Планировщик задач — Диплом

**Джамбулат Мягчиев** | Спринт 13–14

---

## Выполнено: **100% + все звёздочки**

- Сервер + `web/`
- `TODO_PORT`, `TODO_DBFILE`
- SQLite + индекс
- `NextDate()`: `d`, `y`, `w`, `m` (включая `-1`, `-2`, месяцы)
- API: `/api/nextdate`, `/api/task`, `/api/tasks`, `/api/task/done`
- Поиск: `?search=текст` или `?search=08.02.2024`
- **Аутентификация**: `/api/signin` → JWT в куке `token`
- **Middleware**: защита всех `/api/*`
- **Docker**: `distroless`, ~30 МБ, volume для БД
- Все тесты: `PASS`

---

## Запуск (Docker)

```bash
# Сборка
docker build -t planner .

# Запуск
docker run -d \
  -p 7540:7540 \
  -v $(pwd)/scheduler.db:/app/scheduler.db \
  -e TODO_PASSWORD=12345 \
  -e TODO_DBFILE=/app/scheduler.db \
  --name planner-app \
  planner