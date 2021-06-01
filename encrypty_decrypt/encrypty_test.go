package encrypty_decrypt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

//移动端解密 https://segmentfault.com/a/1190000006863198?utm_source=sf-related
func TestAES(t *testing.T) {
	source := "hello world"
	fmt.Println("原字符：", source)
	//16byte密钥
	key := "1443flfsaWfdas"
	encryptCode := AESEncrypt([]byte(source), []byte(key))
	fmt.Println("密文：", string(encryptCode))

	decryptCode := AESDecrypt(encryptCode, []byte(key))

	fmt.Println("解密：", string(decryptCode))
}

func TestAes(t *testing.T) {
	orig := "hello world"
	key := "0123456789012345"
	fmt.Println("原文：", orig)
	encryptCode := AesEncrypt(orig, key)
	fmt.Println("密文：", encryptCode)
	decryptCode := AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)
}

func TestAesctr(t *testing.T) {
	source := "hello world"
	fmt.Println("原字符：", source)

	key := "1443flfsaWfdasds"
	encryptCode, _ := aesCtrCrypt([]byte(source), []byte(key))
	fmt.Println("密文：", string(encryptCode))

	decryptCode, _ := aesCtrCrypt(encryptCode, []byte(key))

	fmt.Println("解密：", string(decryptCode))
}

func TestAesDecryptCFB(t *testing.T) {
	source := "hello world"
	fmt.Println("原字符：", source)
	key := "ABCDEFGHIJKLMNO1" //16位
	encryptCode := AesEncryptCFB([]byte(source), []byte(key))
	fmt.Println("密文：", hex.EncodeToString(encryptCode))
	decryptCode := AesDecryptCFB(encryptCode, []byte(key))

	fmt.Println("解密：", string(decryptCode))
}

func TestOFB(t *testing.T) {
	source := "hello world"
	fmt.Println("原字符：", source)
	key := "1111111111111111" //16位  32位均可
	encryptCode, _ := aesEncryptOFB([]byte(source), []byte(key))
	fmt.Println("密文：", hex.EncodeToString(encryptCode))
	decryptCode, _ := aesDecryptOFB(encryptCode, []byte(key))

	fmt.Println("解密：", string(decryptCode))
}

func TestRSA(t *testing.T) {
	data, _ := RsaEncrypt([]byte("hello world"))
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	origData, _ := RsaDecrypt(data)
	fmt.Println(string(origData))
}
