package static

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed dist/*
var StaticFile embed.FS

// writeFile 将嵌入式文件的内容写入本地文件系统。
func writeFile(fsys fs.FS, embeddedPath, localPath string) error {
	data, err := fs.ReadFile(fsys, embeddedPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(localPath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(localPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// writeDir 递归地将嵌入式文件系统中的所有文件和目录写入本地文件系统。
func writeDir(fsys fs.FS, embeddedPath, localPath string) error {
	entries, err := fs.ReadDir(fsys, embeddedPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		embeddedEntryPath := embeddedPath + "/" + entry.Name()
		localEntryPath := localPath + "/" + entry.Name()

		if entry.IsDir() {
			if err := writeDir(fsys, embeddedEntryPath, localEntryPath); err != nil {
				return err
			}
		} else {
			if err := writeFile(fsys, embeddedEntryPath, localEntryPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func WriteStaticFiles(localDir string) {
	if err := writeDir(StaticFile, "dist", localDir); err != nil {
		log.Fatal("写入静态文件到本地目录时出错:", err)
	}
}
