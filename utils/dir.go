package utils

import (
	"os"
)

// 判斷資料夾是否存在，不存在就新增
func PathExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if !os.IsExist(err) {
			// 沒有資料夾，那就建立
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
