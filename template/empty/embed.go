package emptytemplatex

import (
	"embed"
	_ "embed"
)

//go:embed *.tmpl
var EmptyTmplData embed.FS
