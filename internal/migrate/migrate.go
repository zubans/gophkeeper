package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Runner executes SQL migrations stored in an embedded filesystem directory.
type Runner struct {
	DB  *sql.DB
	FS  embed.FS
	Dir string
}

// EnsureSchemaTable ensures the schema_migrations bookkeeping table exists.
func (r *Runner) EnsureSchemaTable(driver string) error {
	var stmt string
	switch driver {
	case "postgres":
		stmt = `CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	default:
		// sqlite
		stmt = `CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	}
	_, err := r.DB.Exec(stmt)
	return err
}

// appliedSet returns a set of already applied migration names.
func (r *Runner) appliedSet() (map[string]struct{}, error) {
	rows, err := r.DB.Query(`SELECT name FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	set := make(map[string]struct{})
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		set[name] = struct{}{}
	}
	return set, nil
}

// Run executes pending migrations in lexical order.
func (r *Runner) Run(driver string) error {
	if err := r.EnsureSchemaTable(driver); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	entries, err := r.FS.ReadDir(r.Dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}
	var files []string
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".sql") {
			files = append(files, name)
		}
	}
	sort.Strings(files)

	applied, err := r.appliedSet()
	if err != nil {
		return fmt.Errorf("load applied set: %w", err)
	}

	for _, f := range files {
		if _, ok := applied[f]; ok {
			continue
		}
		path := filepath.Join(r.Dir, f)
		content, err := r.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}
		tx, err := r.DB.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", f, err)
		}
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec migration %s: %w", f, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (name) VALUES ($1)`, f); err != nil {
			// Try sqlite placeholder if postgres one failed
			if _, err2 := tx.Exec(`INSERT INTO schema_migrations (name) VALUES (?)`, f); err2 != nil {
				tx.Rollback()
				return fmt.Errorf("record migration %s: %v / %v", f, err, err2)
			}
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", f, err)
		}
	}
	return nil
}
