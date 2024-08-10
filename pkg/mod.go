package pkg

import (
	"os"

	"golang.org/x/mod/modfile"
)

func GetModuleName(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}
