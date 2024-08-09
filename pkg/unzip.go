package pkg

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unzip 将ZIP文件解压到指定的目标目录中。
// zipFileName 是要解压的ZIP文件的路径。
// destDir 是解压目标目录的路径。
func Unzip(zipFileName, destDir string) error {
	// 打开ZIP文件
	r, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建目标目录（如果不存在）
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 遍历ZIP文件中的每个文件/目录
	for _, f := range r.File {
		err := unzipFile(f, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// unzipFile 将ZIP文件中的单个文件或目录解压到目标目录中。
func unzipFile(f *zip.File, destDir string) error {
	// 构建目标路径
	path := filepath.Join(destDir, f.Name)

	// 如果是目录，创建目录并返回
	if f.FileInfo().IsDir() {
		return os.MkdirAll(path, f.Mode())
	}

	// 创建目标文件的父目录
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// 打开ZIP文件中的文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// 创建目标文件
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 将ZIP文件中的文件内容复制到目标文件中
	_, err = io.Copy(outFile, rc)
	return err
}
