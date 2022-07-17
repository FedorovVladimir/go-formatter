# go-formatter

Конфигурируемый форматер кода для Go на Go

## Установка и использование

// TODO

## Конфигурация

Файл config.json нужно положить в корень проекта

Пример конфигурации

```json
{
  "formatters": [
    {
      "name": "formatter_order",
      "enabled": true
    }
  ]
}
```

## Правила форматирования

### formatter_order

До

```go
package main

var v1 = "v1"

type e struct{}

var v2 = "v2"
```

После

```go
package main

var v1 = "v1"

var v2 = "v2"

type e struct{}
```

## Как пользоваться

### Запуск

Для запуска на весь проект

```shell
go-formatter -fix .
```

Для запуска на отдельную директорию

```shell
go-formatter -fix DIR_PATH
```

Также команду можно выполнять по сохранению или на git pre_commit hook
