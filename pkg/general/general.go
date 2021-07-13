package general

import (
	"os"
	"path/filepath"
)

// 基础函数

func GetFilePath(file, ip, virpath string) (path string) {
	if file == "" || ip == "" || virpath == "" {
		path = ""
	} else {
		path += "\\\\"
		path += ip
		path += "\\"
		path += virpath
		path += "\\"
		path += file
	}
	return
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
