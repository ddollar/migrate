package migrate

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/go-pg/pg/v10"
)

func Run(db *pg.DB, migrations fs.FS) error {
	e := &Engine{
		db: db,
		fs: migrations,
	}

	if err := e.Initialize(); err != nil {
		return err
	}

	ms, err := e.Pending()
	if err != nil {
		return err
	}

	for _, m := range ms {
		fmt.Printf("%s: ", m)

		if err := e.Migrate(m); err != nil {
			fmt.Printf("%s\n", err)
			return errors.New("migration failed")
		} else {
			fmt.Println("OK")
		}
	}

	return nil
}
