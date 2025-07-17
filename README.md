# 15_07_2025

## Основные компоненты
1. Worker Pool
Реализован паттерн Worker Pool для ограничения количества одновременно выполняемых задач архивации (максимум 3). Это позволяет эффективно управлять ресурсами и избегать перегрузки системы.

2. Graceful Shutdown
Система корректно обрабатывает сигналы завершения, дожидается завершения текущих задач перед выходом.

4. Конфигурация
Настройки порта и разрешенных расширений файлов можно задавать через флаги командной строки.

---

## Для начала

**Сборка:**

1. Клонируйте репозиторий 15_07_2025:
```
git clone https://github.com/HACK3R911/15_07_2025
```
2. Перейдите в каталог проекта:
```
cd 15_07_2025
```
3. Установите зависимости проекта:
```
go mod init 15_07_2025
```
4. Запуск на порту 8080
```
go run ./cmd/server/main.go
```
или
```
go run ./cmd/server/main.go -port {free_port}
```

для фильтрации типов объектов можно использовать флаг `-ext`, например:
```
go run ./cmd/server/main.go -ext pdf,png
```
будет работать только с pdf и png


**Доступные типы объектов**: pdf, jpg, jpeg, png

---

## API Endpoints

`POST /tasks` - Создать новую задачу архивации
```JSON
{
    "id": "dbeny4xi6dug",
    "urls": null,
    "status": "pending",
    "created_at": "2025-07-18T00:54:52.8260182+03:00"
}
```

`GET /tasks/{id}` - Просмотр задачи по id

`GET /tasks` - Получение списка всех задач
```JSON
[
    {
        "id": "dbenax63v05g",
        "urls": null,
        "status": "pending",
        "created_at": "2025-07-18T00:24:33.5561353+03:00"
    },
    {
        "id": "dbenaybgb4p8",
        "urls": null,
        "status": "pending",
        "created_at": "2025-07-18T00:24:36.0561563+03:00"
    }
]
```
`DELETE /tasks/{id}` - Удаление задачи по id 


`POST /tasks/{id}/urls` - Добавить URL в задачу
request:
```JSON
{"url": "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f2/Felis_silvestris_silvestris_small_gradual_decrease_of_quality_-_JPEG_compression.jpg/250px-Felis_silvestris_silvestris_small_gradual_decrease_of_quality_-_JPEG_compression.jpg"}
```
response:
```
{
    "id": "dbeny4xi6dug",
    "urls": [
        "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f2/Felis_silvestris_silvestris_small_gradual_decrease_of_quality_-_JPEG_compression.jpg/250px-Felis_silvestris_silvestris_small_gradual_decrease_of_quality_-_JPEG_compression.jpg"
    ],
    "status": "pending",
    "created_at": "2025-07-18T00:54:52.8260182+03:00"
}
```

`GET /tasks/{id}/download` - Скачивание архива

## Проблемы

При добавлении в задачу файлов с расширением .pdf и скачивании zip-архива нельзя просмотреть файлы, для этого можно поменять расширение архива на .7z