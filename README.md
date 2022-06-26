# go-formatter

Конфигурируемый форматер кода для Go на Go

## Установка и использование

## Конфигурация

Файл config.json нужно положить в корень проекта

```json
{
  "formatters": [
    {
      "name": "many_arguments",
      "on": true
    },
    ...
  ]
}
```

Пример конфигурации

```json
{
  "formatters": [
    {
      "name": "context_first_parameter",
      "on": true
    },
    {
      "name": "empty_func_body",
      "on": true
    },
    {
      "name": "grouped_vars",
      "on": true
    },
    {
      "name": "many_arguments",
      "on": true
    },
    {
      "name": "methods_with_star_and_rename",
      "on": true
    },
    {
      "name": "new_line",
      "on": true
    },
    {
      "name": "rm_ignore_vars",
      "on": true
    },
    {
      "name": "start_enums_at_one",
      "on": true
    },
    {
      "name": "with",
      "on": true
    },
    {
      "name": "order",
      "on": true
    }
  ]
}
```

## Правила форматирования

### context_first_parameter

До

```go
package a

import "context"

func ctxFirstParameter(a int, ctx context.Context) {

}
```

После

```go
package a

import "context"

func ctxFirstParameter(ctx context.Context, a int) {

}

```

### empty_func_body

До

```go
package a

func WithCarColor(color string) {

}

```

После

```go
package a

func WithCarColor(color string) {}
```

### grouped_vars

До

```go
package main

var a = "a"
var b = "b"

var (
	с = "с"
)

var (
	d = "d"
)
var e = "e"
```

После

```go
package main

var (
	a = "a"
	b = "b"
)

var с = "с"

var (
	d = "d"
	e = "e"
)
```

### many_arguments

До

```go
package a

func ManyArguments(a int, b int, c string, d bool, l int64, m string) {

}


```

После

```go
package a

func ManyArguments(
	a int,
	b int,
	c string,
	d bool,
	l int64,
	m string,
) {

}


```

### methods_with_star_and_rename

До

```go
package main

type car struct{}

func (c *car) run() {}

func (e car) stop() {}


```

После

```go
package main

type car struct{}

func (c *car) run() {}

func (c *car) stop() {}


```

### new_line

До

```go
package main

type car struct{}

func (c *car) run()  {}
func (c *car) stop() {}
```

После

```go
package main

type car struct{}

func (c *car) run() {}

func (c *car) stop() {}
```

### rm_ignore_vars

До

```go
package main

var _, b = "b", "b"

var _, _ = "b", "b"
```

После

```go
package main

var _, b = "b", "b"
```

### start_enums_at_one

До

```go
package main

type Operation int

const (
	Add Operation = iota
	Subtract
	Multiply
)
```

После

```go
package main

type Operation int

const (
	Add Operation = iota + 1
	Subtract
	Multiply
)
```

### with

До

```go
package a

func CarWithColor(color string) {}

func carWithPrice(price float64) {}
```

После

```go
package a

func WithCarColor(color string) {}

func withCarPrice(price float64) {}
```

## Как пользоваться

### Запуск

Для запуска на весь проект

```shell
myfmt -fix .
```

Для запуска на отдельную директорию

```shell
myfmt -fix DIR_PATH
```

Также команду можно выполнять по сохранению или на git pre_commit hook
