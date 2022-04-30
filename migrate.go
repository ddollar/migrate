package migrate

import (
	"io/fs"

	"github.com/ddollar/stdcli"
	"github.com/go-pg/pg/v10"
)

func New(dburl string, migrations fs.FS) (*CLI, error) {
	opts, err := pg.ParseURL(dburl)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	e := &Engine{
		db: db,
		fs: migrations,
	}

	c := &CLI{
		cli:    stdcli.New("migrate", ""),
		engine: e,
	}

	c.Register()

	return c, nil
}
