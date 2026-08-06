package migration

import (
	"embed"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations/*.sql
var MigrationFS embed.FS
