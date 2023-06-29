package safe_live

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

const pw_length = 32

func EncryptLines(lines []string, key string) ([]string, error) {
	preLine := key
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		nkey := newkey(preLine, key)
		salt := newkey(nkey, key)
		str, err := encrypt(salt+lines[i], nkey)
		if err != nil {
			return nil, err
		}
		preLine = lines[i]
		lines[i] = str
	}

	return lines, nil
}

func DecryptLines(lines []string, key string) ([]string, error) {
	preLine := key
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		nkey := newkey(preLine, key)
		str, err := decrypt(lines[i], nkey)
		if err != nil {
			return nil, err
		}
		lines[i] = str[pw_length:]
		preLine = lines[i]
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

var encoding = base64.NewEncoding("abcdefghijklmnOPQRSTUVWXYZ0123456789~!@#$%^&*()_+{}|:A>?<,./;'[]")

func newkey(str, key string) string {
	sha := sha256.Sum256([]byte(key + str))
	return encoding.EncodeToString(sha[:])[6 : 6+pw_length]
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

func EncryptPassword(key string) (string, string, error) {
	randomNum := fmt.Sprintf("%d", time.Now().UnixNano())
	nkey := newkey(randomNum, key)
	enKey, err := encrypt(nkey, newkey(key, key))
	return nkey, enKey, err
}

func DecryptPassword(key string, line string) (string, error) {
	return decrypt(line, newkey(key, key))
}
