package encrypty_decrypt

import (
	"crypto/md5"
	"encoding/hex"
)

/*
加密模式：
	md5
*/
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
