package templatex

import (
	"embed"
	_ "embed"
)

//go:embed *.tmpl
var TmplSysUserData embed.FS
