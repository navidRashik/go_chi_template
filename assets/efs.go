package assets

import (
	"embed"
)

//go:embed "migrations" "VERSION"
var EmbeddedFiles embed.FS
