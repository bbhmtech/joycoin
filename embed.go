package joycoin

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist
var embedAssets embed.FS

var StaticContent, _ = fs.Sub(embedAssets, "frontend/dist")
