package migrate

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type Migration struct {
	Version string
	Down    []byte
	Up      []byte
}

type Migrations []Migration

func LoadMigrations(e *Engine) (Migrations, error) {
	raw := map[string]Migration{}

	err := fs.WalkDir(e.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		parts := strings.SplitN(d.Name(), ".", 3)
		if len(parts) != 3 {
			return fmt.Errorf("invalid migration: %s", d.Name())
		}
		mm := raw[parts[0]]
		switch parts[1] {
		case "up":
			mm.Up, _ = fs.ReadFile(e.fs, path)
		case "down":
			mm.Down, _ = fs.ReadFile(e.fs, path)
		default:
			return fmt.Errorf("invalid migration: %s", d.Name())
		}
		raw[parts[0]] = mm
		return nil
	})
	if err != nil {
		return nil, err
	}

	for k, m := range raw {
		if m.Up == nil {
			return nil, fmt.Errorf("missing migration file: %s.up.sql", k)
		}
		if m.Down == nil {
			return nil, fmt.Errorf("missing migration file: %s.down.sql", k)
		}
	}

	ms := Migrations{}

	for k, m := range raw {
		m.Version = k
		ms = append(ms, m)
	}

	sort.Slice(ms, func(i, j int) bool { return ms[i].Version < ms[j].Version })

	return ms, nil
}

func (ms Migrations) Find(version string) (Migration, bool) {
	for _, m := range ms {
		if m.Version == version {
			return m, true
		}
	}

	return Migration{}, false
}
