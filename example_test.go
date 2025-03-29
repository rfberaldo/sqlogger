package sqlogger_test

import (
	"context"
	"log"
	"log/slog"

	"github.com/mattn/go-sqlite3"
	"github.com/rfberaldo/sqlogger"
)

func ExampleOpen() {
	// [sqlogger.Open] is similar to [sql.Open], with two additional parameters:
	// [slog.Logger] and [sqlogger.Options]
	db, err := sqlogger.Open("sqlite3", ":memory:", slog.Default(), nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE user (id INT PRIMARY KEY, name TEXT")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleOpen_options() {
	opts := &sqlogger.Options{
		CleanQuery: true,
	}

	// use sqlogger.Options as fourth parameter
	db, err := sqlogger.Open("sqlite3", ":memory:", slog.Default(), opts)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE user (id INT PRIMARY KEY, name TEXT")
	if err != nil {
		log.Fatal(err)
	}
}

type Logger struct {
	slog *slog.Logger
}

func (l *Logger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	attrs = append(attrs, slog.Any("userId", ctx.Value("userId")))
	l.slog.LogAttrs(ctx, level, msg, attrs...)
}

func ExampleOpen_contextValues() {
	// to log values from context, wrap the slog object
	lg := &Logger{slog.Default()}
	db, err := sqlogger.Open("sqlite3", ":memory:", lg, nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), "userId", 42)
	_, err = db.ExecContext(ctx, "CREATE TABLE user (id INT PRIMARY KEY, name TEXT")
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleNew() {
	// [sqlogger.New] receive the [driver.Driver] instead of the driver name
	db := sqlogger.New(&sqlite3.SQLiteDriver{}, ":memory:", slog.Default(), nil)

	_, err := db.Exec("CREATE TABLE user (id INT PRIMARY KEY, name TEXT")
	if err != nil {
		log.Fatal(err)
	}
}
