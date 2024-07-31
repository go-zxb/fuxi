package templatex

import (
	"embed"
	_ "embed"
)

//go:embed *.tmpl
var TmplData embed.FS
