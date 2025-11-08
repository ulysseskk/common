package compress

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DecompressTarGz(filePath, destFile string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	// 创建gzip解压器
	unGzipStream, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer unGzipStream.Close()

	// 创建tar阅读器
	tarReader := tar.NewReader(unGzipStream)

	// 遍历tar档案中的每一项
	for {
		header, err := tarReader.Next()

		// 如果没有更多的文件，则跳出循环
		if err == io.EOF {
			break
		}

		// 遇到错误时返回错误
		if err != nil {
			return err
		}
		// 忽略以"._"开头的文件
		splits := strings.Split(header.Name, "/")
		if strings.HasPrefix(splits[len(splits)-1], "._") {
			continue

		}
		// 目标路径
		targetPath := destFile

		// 根据tar中对象的类型来处理（文件或目录）
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(filepath.Join(targetPath, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			buffer := &bytes.Buffer{}
			// 将文件内容写入到buffer中
			if _, err := io.Copy(buffer, tarReader); err != nil {
				return err
			}
			// 创建文件
			err = os.WriteFile(filepath.Join(targetPath, header.Name), buffer.Bytes(), os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
