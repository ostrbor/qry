package qry

import (
	"context"
)

// TruncateAll is a function that truncates all tables in the PostgreSQL database,
// excluding the migrations table. It uses an anonymous function (DO block) to execute
// the SQL command. This function is intended to be used in a testing environment where
// you need to clean up the database before each test.
func TruncateAll(db Querier) {
	sql := `
 DO $$
 BEGIN
  EXECUTE
  (SELECT 'TRUNCATE TABLE '
  -- Escape table names using %I to handle potential keywords or spaces
  || string_agg(format('%I', tablename), ', ')
  -- CASCADE ensures truncating tables with references to the specified tables
  -- RESTART IDENTITY restarts sequences owned by truncated tables
  || ' CASCADE RESTART IDENTITY'
  FROM pg_tables
  WHERE schemaname = 'public' AND tablename != '` + Migrations + `'
  );
 END
 $$;`
	_, err := db.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}
}
