// Package migrations holds the SQL migrations, built into the program so they
// can be applied without depending on the current working directory.
package migrations

import "embed"

// FS contains every migration file in this directory.
//
//go:embed *.sql
var FS embed.FS
