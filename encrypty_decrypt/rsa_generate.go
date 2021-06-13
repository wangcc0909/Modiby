package encrypty_decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateRSAKey(keySize int, priPath string, pubPath string) error {
	//生成密钥对
	priKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}
	//保存密钥对  1。x509序列化  2。pem序列化
	priByte := x509.MarshalPKCS1PrivateKey(priKey)
	//创建私钥保存的文件
	file, err := os.Create(priPath)
	if err != nil {
		return err
	}
	defer file.Close()
	//创建block
	block := pem.Block{
		Type:  "RSA Private Key",
		Bytes: priByte,
	}
	err = pem.Encode(file, &block)
	if err != nil {
		return err
	}

	//保存公钥
	pubKey := priKey.PublicKey
	pubByte, err := x509.MarshalPKIXPublicKey(&pubKey) //这里一定要取地址
	if err != nil {
		return err
	}
	pubFile, err := os.Create(pubPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()
	pubBlock := pem.Block{
		Type:  "RSA Public Key",
		Bytes: pubByte,
	}
	err = pem.Encode(pubFile, &pubBlock)
	return err
}

//读取文件内容
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, fileInfo.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
