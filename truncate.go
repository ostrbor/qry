package qry

import (
	"context"
	"fmt"
	"strings"
)

func Truncate(db Querier, tables []string) {
	sql := fmt.Sprintf(
		`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`,
		strings.Join(tables, ","))
	_, err := db.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}
}
