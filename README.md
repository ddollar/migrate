# migrate

Postgres migration tool

## Installation

Create `cmd/migrate/main.go` in your project:

```golang
package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/ddollar/migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
	}
}

func run() error {
	m, err := migrate.New(os.Getenv("POSTGRES_URL"), migrations)
	if err != nil {
		return err
	}

	if err := m.Run(os.Args[1:]); err != nil {
		return err
	}

	return nil
}
```

Add this target to your `Makefile`:

```makefile
migration:
        $(if $(name),,$(error name is not set))
        touch cmd/migrate/migrations/$(date)_$(name).up.sql
        touch cmd/migrate/migrations/$(date)_$(name).down.sql
```

## Usage

### Create a migration

```
$ make migration name=create_users
```

### Run migrations locally

```
$ go run cmd/migrate/main.go
```

### Use in production

Make sure the `migrate` binary makes it into your final deployment. The migration files will
be compiled into it.
