package safe_live

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
)

func EncryptLines(lines []string, key string) ([]string, error) {
	f := fillFormat(key)
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		k := fmt.Sprintf(f, key, i)
		str, err := encrypt(lines[i], k)
		if err != nil {
			return nil, err
		}
		lines[i] = str
	}

	return lines, nil
}

func DecryptLines(lines []string, key string) ([]string, error) {
	f := fillFormat(key)
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		k := fmt.Sprintf(f, key, i)
		str, err := decrypt(lines[i], k)
		if err != nil {
			return nil, err
		}
		lines[i] = str
	}
	return lines, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func fillFormat(key string) string {
	kl := len(key)
	f := "%s%0" + strconv.Itoa(16-kl) + "d"
	return f
}

func encrypt(str, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	origData := []byte(str)
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(key)[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func decrypt(str, key string) (string, error) {

	crypted, err := base64.StdEncoding.DecodeString(str)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(key)[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return string(origData), nil
}
