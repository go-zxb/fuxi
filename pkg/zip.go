package pkg

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

// FilesToZip 将文件添加到已有的ZIP文件中，如果ZIP文件不存在则创建新的ZIP文件。
// zipFileName 是现有ZIP文件的路径。
// files 是要添加到ZIP文件中的文件列表。
func FilesToZip(zipFileName string, files []string) error {
	// 检查ZIP文件是否存在
	_, err := os.Stat(zipFileName)
	zipFileExists := !os.IsNotExist(err)

	// 创建一个新的ZIP文件（如果ZIP文件不存在）或打开现有的ZIP文件
	var newZipFile *os.File
	if zipFileExists {
		newZipFile, err = os.OpenFile(zipFileName, os.O_RDWR|os.O_CREATE, 0666)
	} else {
		newZipFile, err = os.Create(zipFileName)
	}
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	// 创建一个新的ZIP写入器
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 如果ZIP文件存在，先复制现有的文件到新的ZIP文件中
	if zipFileExists {
		// 打开现有的ZIP文件进行读取
		existingZipFile, err := zip.OpenReader(zipFileName)
		if err != nil {
			return err
		}
		defer existingZipFile.Close()

		for _, file := range existingZipFile.File {
			err := copyFileToZip(zipWriter, file)
			if err != nil {
				return err
			}
		}
	}

	// 添加新的文件到ZIP文件中
	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			if strings.Contains(err.Error(), "The system cannot find the path specified") {
				continue
			}
			return err
		}
	}

	// 关闭ZIP写入器以确保所有数据被刷新
	if err := zipWriter.Close(); err != nil {
		return err
	}

	return nil
}

// 将ZIP读取器中的文件复制到ZIP写入器中。
func copyFileToZip(zipWriter *zip.Writer, file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	header, err := zip.FileInfoHeader(file.FileInfo())
	if err != nil {
		return err
	}

	header.Name = file.Name
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, rc)
	return err
}

// addFileToZip 将文件添加到ZIP存档中。
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}
