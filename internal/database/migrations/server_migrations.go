package migrations

import "embed"

//go:embed server/*.sql
var ServerMigrations embed.FS
