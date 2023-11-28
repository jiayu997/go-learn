package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"bytes"
	"crypto"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLicenseHandler(t *testing.T) {
	licInstance := &License{
		"0",
		[]string{"0A-00-27-00-00-12"},
		"BFEBFBFF000306C3",
		"0",
		strconv.FormatInt(time.Now().UnixNano(),10),
	}

	licBytes,_ := json.Marshal(licInstance)
	fmt.Println("origin:",licBytes)
	b, err := rsaEncrypt(licBytes)
	if err!=nil {
		fmt.Println("加密失败:",err)
	}
	fmt.Print("encrypted:",b)
	body := bytes.NewBuffer(b)

	res,err := http.Post("http://localhost:8082","application/octet-stream",body)

	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	signErr := resUnsign([]byte(string(licInstance.RequestTime)+":ok"),result)
	if signErr != nil {
		fmt.Println("签名验证失败")
	}else{
		fmt.Println("签名验证成功")
	}

}

func readBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func rsaEncrypt(origData []byte) ([]byte, error) {
	//pk,_ := readBytes(`public.pem`) ;
	pk := []byte(`MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDNv/C8QParDSKb64npiLEszp0U9nL8/cSmazM28+bWWPzOsLcs0xhushf//zi/7890p0zidbij88yKgV9QyIEmx3r0uXlOc8cc2loB6ioU8bCLNZA8/AFeImqKqq/9y+ywknj1ZhOE/KVgFI+dFk58eAuIRfrekHFGx2SYtmcv+wIDAQAB`)

	block, _ := pem.Decode(pk)
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

func resUnsign(message []byte, sig []byte) error {
	//pk,_ := readBytes(`public.pem`) ;
	pk := []byte(`MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDNv/C8QParDSKb64npiLEszp0U9nL8/cSmazM28+bWWPzOsLcs0xhushf//zi/7890p0zidbij88yKgV9QyIEmx3r0uXlOc8cc2loB6ioU8bCLNZA8/AFeImqKqq/9y+ywknj1ZhOE/KVgFI+dFk58eAuIRfrekHFGx2SYtmcv+wIDAQAB`)

	block, _ := pem.Decode(pk)
	if block == nil {
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, d, sig)
}




