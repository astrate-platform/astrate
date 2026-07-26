// Package migrations embeds Astrate's SQL schema migrations (transcribed from
// docs/DESIGN.md §2.2–2.5) so the binary self-migrates at startup through
// golang-migrate's iofs source (docs/ROADMAP.md §3.1 file 2.6). A Go file is
// required here because go:embed cannot reach outside the declaring package's
// directory.
package migrations

import "embed"

// FS holds every numbered .up.sql/.down.sql migration file in this directory.
//
//go:embed *.sql
var FS embed.FS
