package qry

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io/fs"
	"sort"
	"strings"
)

// Migrations is a constant that represents the name of the migrations table.
var Migrations = "migrations"

func Migrate(db Querier, migrations embed.FS) {
	entries, err := migrations.ReadDir(".")
	if err != nil {
		panic(err)
	}
	files := filenames(entries)
	err = Transaction(db, func(tx pgx.Tx) error {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (filename TEXT UNIQUE NOT NULL)", Migrations)
		_, err = tx.Exec(context.Background(), sql)
		if err != nil {
			return err
		}
		for _, filename := range files {
			content, err := migrations.ReadFile(filename)
			if err != nil {
				return err
			}
			err = migrate(tx, filename, string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func filenames(entries []fs.DirEntry) (names []string) {
	names = make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

func migrate(tx Querier, filename, content string) error {
	ok, err := exists(tx, filename)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	_, err = tx.Exec(context.Background(), content)
	if err != nil {
		return err
	}
	return insert(tx, filename)
}

func insert(tx Querier, filename string) (err error) {
	sql := fmt.Sprintf("INSERT INTO %s (filename) VALUES ($1)", Migrations)
	_, err = tx.Exec(context.Background(), sql, filename)
	return err
}

func exists(tx Querier, version string) (ok bool, err error) {
	sql := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE filename=$1)", Migrations)
	err = tx.QueryRow(context.Background(), sql, version).Scan(&ok)
	return ok, err
}
