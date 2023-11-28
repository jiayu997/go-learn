package utils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func bytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func RsaEncrypt(publickey []byte, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}


// 解密
func RsaDecrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)

	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var decryptResult []byte
	var i,offset int
	var inputLen = len(ciphertext)
	for inputLen - offset > 0 {
		//fmt.Println(len(ciphertext)-offset)
		if len(ciphertext)-offset>128 {
			result,errr := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext[i*128:(i*128)+128])
			if errr != nil {
				return nil,errr
			}
			//fmt.Println(ciphertext[i*128:(i*128)+128])
			//fmt.Println(string(result))
			decryptResult = bytesCombine(decryptResult,result)
		}else{
			result,errr := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext[i*128:])
			if errr != nil {
				return nil,errr
		}
			decryptResult = bytesCombine(decryptResult,result)
		}
		i++
		offset = i * 128
	}
	//fmt.Println("de:",decryptResult)
	return decryptResult,nil
}
//签名
func RsaSign(privateKey []byte, ciphertext []byte) ([]byte, error){
	block, _ := pem.Decode(privateKey)

	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(ciphertext)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, d)
}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")

	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil

}
