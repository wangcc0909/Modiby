package encrypty_decrypt

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

/*
加密模式：
hmac
*/

//key随意设置 data 要加密的数据
func Hmac(key, data string) string {
	hash := hmac.New(md5.New, []byte(key)) //创建对应的md5哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

func HmacSha256(key, data string) string {
	hash := hmac.New(sha256.New, []byte(key)) //创建对应的sha256哈希加密算法
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

func Sha1(data string) string {
	sha12 := sha1.New()
	sha12.Write([]byte(data))
	return hex.EncodeToString(sha12.Sum([]byte("")))
}
