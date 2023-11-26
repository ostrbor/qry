package qry

import (
	"context"
)

func TruncateAll(db Querier) {
	sql := `
	-- Truncate all tables in the public schema
	CREATE OR REPLACE FUNCTION truncate_all() RETURNS void LANGUAGE plpgsql AS
	$$
	BEGIN
		EXECUTE
		(SELECT 'TRUNCATE TABLE '
		-- Escape table names using %I to handle potential keywords or spaces
		|| string_agg(format('%I', tablename), ', ')
		-- CASCADE ensures truncating tables with references to the specified tables
		-- RESTART IDENTITY restarts sequences owned by truncated tables
		|| ' CASCADE RESTART IDENTITY'
		FROM pg_tables
		WHERE schemaname = 'public'
		);
	END
	$$;
	-- Call the truncate_all function
	SELECT truncate_all();`
	_, err := db.Exec(context.Background(), sql)
	if err != nil {
		panic(err)
	}
}
