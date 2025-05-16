package utils

import (
	"os"
)

func CreateDir(path string) error {
	// 判断文件夹是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func RemoveDir(path string) error {
	// 判断文件夹是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件夹不存在，直接返回
		return nil
	}

	// 删除文件夹及其内容
	return os.RemoveAll(path)
}

func RemoveFile(path string) error {
	// 判断文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在，直接返回
		return nil
	}

	// 删除文件
	return os.Remove(path)
}
