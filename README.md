# calc_go_yandex
## Финальное задание Yandex LMS 2 спринт.

Данная программа вычисляет значение арифметического выражения.

Программа поддерживает ввод рациональных чисел и арифметичексие операции (+ - * /).

Доступен графический интерфейс по адресу `http://localhost:8080/api/v1/`

## Структура проекта

Проект состоит из двух блоков:

**Сервер (orchestrator)** – управляет вычислениями, раздаёт задачи вычислителям, обрабатывает данные.

**Вычислитель (agent)** – получает задачу и возвращает полученный результат 

    calc_go_yandex_2/
    ├── cmd/
    │    └── main.go
    ├── internal/
    │    └── application/
    │        ├── agent.go
    │        ├── agent_test.go
    │        ├── application.go
    │        ├── application_test.go
    │        └── orchestrator.go
    ├── pkg/
    │    └── calcucaltion/
    │        ├── calculation.go
    │        ├── calculation_test.go
    │        └── errors.go
    ├── static/
    │        └── index.html
    ├── go.mod
    └── README.md

**main.go** - Точка входа в программу.

**application.go** - Web-сервер, методы для запуска программы.

**orchestrator.go** - Структуры и методы необходимые для работы оркестратора, экспортируются в application.go.

**agent.go** - Вычислитель, который будет вызван в отдельные горутины.

**calculation.go** - Функция calc для разбения выражния на действия (action). Action представляет из себя улучшенный task

**agent_test.go** - Тестирование вычислителя

**calclulation_test.go** - Тестирование разбияния на действия.

**application_test.go** - Тестирование оркестратора и вычислителя

**index.html** - Frontend

## Endpoint-ы web-сервера

**localhost/api/v1/calculate** - добавление арифметического выражения с помощью POST запроса `{"expression":"Выражение"}`

---
**localhost/api/v1/expressions** - получение текущего списка выражений
  
Тело ответа:

    {
        "expressions": [
            {
                "id": <идентификатор выражения>,
                "status": <статус вычисления выражения>,
                "result": <результат выражения>
            },
            {
                "id": <идентификатор выражения>,
                "status": <статус вычисления выражения>,
                "result": <результат выражения>
            }
        ]
    }

---
**localhost/api/v1/expressions/:id** - получение выражения по его id

Тело ответа:

    {
        "expression":
            {
                "id": <идентификатор выражения>,
                "status": <статус вычисления выражения>,
                "result": <результат выражения>
            }
    }

---
**localhost/internal/task** - получение задания вычислителем от сервера (не обращайтесь вручную, если не хотите ничего сломать)

Тело ответа:

    {
        "task":
            {
                "id": <идентификатор задачи>,
                "arg1": <имя первого аргумента>,
                "arg2": <имя второго аргумента>,
                "operation": <операция>,
                "operation_time": <время выполнения операции>
            }
    }

---
**localhost/api/v1/calculate** - прием результата обработки данных с помощью POST запроса `{"id":"идентификатор задания", "result": значения результата}` (не обращайтесь вручную, если не хотите ничего сломать)

## Чтобы запустить программу, необходимо:
### **Введите это в powershell:**
1) Скачать актуальную версию `git clone git@github.com:hidnt/calc_go_yandex_2.git`
2) Перейти в созданную папку `cd calc_go_yandex_2`
3) Запустить программу `$env:COMPUTING_POWER=3; $env:PORT=8080; $env:TIME_ADDITION_MS=500; $env:TIME_SUBTRACTION_MS=500; $env:TIME_MULTIPLICATIONS_MS=500; $env:TIME_DIVISIONS_MS=500; go run ./cmd/main.go`

COMPUTING_POWER - количество вычистителей

TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах

TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах

TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах

TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах


## Чтобы запустить тесты, необходимо:
1) Скачать актуальную версию `git clone git@github.com:hidnt/calc_go_yandex_2.git`
2) Перейти в созданную папку `cd calc_go_yandex_2`
3) Запустить тестирование `go test -v ./...`

## Примеры работы программы:
`curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"(2+2-(-2+7)*2)/2\" }" http:/localhost:8080/api/v1/calculate`

Возвращает 

    {
        "id": "0"
    }

Код 201.

---

`curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"123-(8*4\" }" http://localhost:8080/api/v1/calculate`

Возвращает

    {
        "id": "0"
    }

Код 422.

---

`curl -X POST -H "Content-Type: application/json" -d "" http://localhost:8080/api/v1/calculate`

Возвращает

    {
        "id": "0"
    }

Код 500.

---

`curl -X GET http://localhost:8080/api/v1/expressions`

Возвращает все выражения

    {
        "expressions": [
            {
                "id": "0",
                "status": "not enough nums",
                "result": 0
            },
            {
                "id": "1",
                "status": "completed",
                "result": 744
            }
        ]
    }

Код 200.

---

`curl -X GET http://localhost:8080/api/v1/expressions/:0`

Возвращает выражение 0

    {
        "expression": {
            "id": "0",
            "status": "not enough nums",
            "result": 0
        }
    }

Код 200.

---

`curl -X GET http://localhost:8080/api/v1/expressions/:123`

Возвращает все выражения, если что-то не найдено

    {
        "expressions": [
            {
                "id": "0",
                "status": "not enough nums",
                "result": 0
            },
            {
                "id": "1",
                "status": "completed",
                "result": 744
            }
        ]
    }

Код 404.