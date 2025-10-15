# SQLogger

[![Test Status](https://github.com/rfberaldo/sqlogger/actions/workflows/test.yaml/badge.svg)](https://github.com/rfberaldo/sqlogger/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rfberaldo/sqlogger)](https://goreportcard.com/report/github.com/rfberaldo/sqlogger)
[![Go Reference](https://pkg.go.dev/badge/github.com/rfberaldo/sqlogger.svg)](https://pkg.go.dev/github.com/rfberaldo/sqlogger)

SQLogger is a lightweight zero-deps logger library for `database/sql`, it returns a standard `sql.DB` object, making it very transparent.
Easy to debug with cascading connections ids.

It's like [sqldb-logger](https://github.com/simukti/sqldb-logger) with `slog` support and more opinionated.

## Example

```bash
2025/03/27 17:52:27 INFO ExecContext query="CREATE TABLE user ( id INT PRIMARY KEY, name VARCHAR(255), age INT )" conn_id=Ty7tgM duration=3.913661ms
2025/03/27 17:52:27 INFO ExecContext query="INSERT INTO user (id, name, age) VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9)" args="[1 Alice 18 2 Rob 38 3 John 4]" conn_id=Ty7tgM duration=614.307µs
2025/03/27 17:52:27 INFO QueryContext query="SELECT count(1) FROM user" conn_id=Ty7tgM duration=311.123µs
2025/03/27 17:52:27 INFO ExecContext query="DELETE FROM user" conn_id=Ty7tgM duration=486.935µs
2025/03/27 17:52:27 INFO ExecContext query="INSERT INTO user (id, name, age) VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9)" args="[1 Alice 18 2 Rob 38 3 John 4]" conn_id=Ty7tgM duration=128.631µs
2025/03/27 17:52:27 INFO QueryContext query="SELECT count(1) FROM user" conn_id=Ty7tgM duration=116.081µs
2025/03/27 17:52:27 INFO ExecContext query="INSERT INTO user (id, name, age) VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9)" args="[1 Alice 18 2 Rob 38 3 John 4]" conn_id=Ty7tgM duration=160.512µs
2025/03/27 17:52:27 ERROR Rollback error="timeout: context already done: context canceled" conn_id=Ty7tgM tx_id=FoPyiM duration=53.201µs
2025/03/27 17:52:27 ERROR ResetSession error="driver: bad connection" conn_id=Ty7tgM duration=490ns
2025/03/27 17:52:27 INFO QueryContext query="SELECT count(1) FROM user" conn_id=MIcA61 duration=645.607µs
2025/03/27 17:52:27 INFO ExecContext query="DROP TABLE user" conn_id=MIcA61 duration=1.360044ms
```

## Getting started

### Install

```bash
go get github.com/rfberaldo/sqlogger
```

### Open a database

```go
// there's 2 ways of opening a database:

// 1. using [sqlogger.Open]
db, err := sqlogger.Open("sqlite3", ":memory:", slog.Default(), nil)

// 2. using [sqlogger.New]
// instead of driver name, it expects the [driver.Driver]
db := sqlogger.New(&sqlite3.SQLiteDriver{}, ":memory:", slog.Default(), nil)
```

## Options

```go
type Options struct {
  // IdGenerator is used to generate the id of connections.
  // Default: 6-length random string.
  IdGenerator func() string

  // CleanQuery removes any redundant whitespace before logging.
  // Default: false.
  CleanQuery bool
}
```

Use `sqlogger.Options` as a fourth parameter of `New` or `Open`:

```go
opts := &sqlogger.Options{
  // options...
}
db, err := sqlogger.Open("sqlite3", ":memory:", slog.Default(), opts)
```

## Logging values from context

For that you can wrap the `slog.Logger` object:

```go
type Logger struct {
  slog *slog.Logger
}

func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
  attrs = append(attrs, slog.Any("userId", ctx.Value("userId")))
  l.slog.LogAttrs(ctx, level, msg, attrs...)
}

lg := &Logger{slog.Default()}
db, err := sqlogger.Open("sqlite3", ":memory:", lg, nil)
```
