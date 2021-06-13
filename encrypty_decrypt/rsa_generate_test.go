package encrypty_decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

const KeySize = 1024 * 4

func TestGenerateRSAKey(t *testing.T) {
	err := GenerateRSAKey(KeySize, "private.pem", "public.pem")
	if err != nil {
		t.Error("generate RSA key error:", err)
	}
	//加密
	buf, err := ReadFile("public.pem")
	if err != nil {
		t.Error("read public file error:", err)
	}
	b, _ := pem.Decode(buf)

	pub, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		t.Error("解析public key：", err)
	}
	pubKey := pub.(*rsa.PublicKey)
	plainText := "我是rsa加密的内容"
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(plainText))
	if err != nil {
		t.Error("加密错误：", err)
	}
	buffer, err := ReadFile("private.pem")
	if err != nil {
		t.Error("读取私钥错误")
	}
	bk, _ := pem.Decode(buffer)
	priKey, err := x509.ParsePKCS1PrivateKey(bk.Bytes)
	if err != nil {
		t.Error("解析私钥：", err)
	}
	result, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherText)
	if err != nil {
		t.Error("解密：", err)
	}
	t.Log("解码之后的内容:   ", string(result))
}
