package templatex

import (
	"embed"
)

//go:embed *.tmpl
var TmplCasbinData embed.FS
