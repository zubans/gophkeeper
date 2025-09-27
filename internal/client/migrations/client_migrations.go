package migrations

import "embed"

//go:embed client/*.sql
var ClientMigrations embed.FS
