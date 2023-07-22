package utils

import (
	"os"
	"regexp"
)

// 检查文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 判断字符串是否匹配正则
func CompareString(src, dst string) (bool, error) {
	re, err := regexp.Compile(dst)
	if err != nil {
		return false, err
	}
	return re.MatchString(src), nil
}
