package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Val(s string) (string, error) {
	// 创建一个新的MD5散列对象
	hash := md5.New()

	// 将字符串写入MD5散列对象
	_, err := hash.Write([]byte(s))
	if err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)
	md5Hex := hex.EncodeToString(hashBytes)
	return md5Hex, nil
}
