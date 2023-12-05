package qry

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Migrate(db Querier, fs embed.FS, files []string) {
	err := Transaction(db, func(tx pgx.Tx) error {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (version TEXT UNIQUE NOT NULL)", Migrations)
		_, err := tx.Exec(context.Background(), sql)
		if err != nil {
			return err
		}
		for _, filename := range files {
			if exists(tx, filename) {
				continue
			}
			migrate(tx, fs, filename)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func migrate(tx pgx.Tx, fs embed.FS, filename string) {
	content, err := fs.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if _, err := tx.Exec(ctx, string(content)); err != nil {
		panic(err)
	}
	sql := fmt.Sprintf("INSERT INTO %s (version) VALUES ($1)", Migrations)
	if _, err := tx.Exec(ctx, sql, filename); err != nil {
		panic(err)
	}
}

func exists(tx pgx.Tx, version string) bool {
	sql := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE version=$1)", Migrations)
	var e bool
	err := tx.QueryRow(context.Background(), sql, version).Scan(&e)
	if err != nil {
		panic(err)
	}
	return e
}
