**Example MustInsert**

```go
package main

import (
	"context"
	"fmt"
	"qry"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Model struct {
	ID int `column:"id,serial" table:"models"`
}

func main() {
	db, _ := pgxpool.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/postgres")
	m := Model{}
	qry.MustInsert(db, &m)
	fmt.Println(m.ID)
}
```

**Example Migrate**

Create 'embed.go' in migrations catalog.

```go
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
```

```go
package main

import (
	"context"
	"qry"
	"github.com/jackc/pgx/v4/pgxpool"
	"migrations"
)

func main() {
	db, _ := pgxpool.Connect(context.Background(), "postgresql://postgres:postgres@localhost:5432/postgres")
	qry.Migrate(db, migrations.FS, []string{"init.sql"})
}
```

**Context**

Custom context cannot be passed from caller.

Reason: to simplify signatures and usually there is no need to stop running query.