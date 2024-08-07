package pkg

import (
	"golang.org/x/mod/modfile"
	"os"
)

func GetModuleName(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}
