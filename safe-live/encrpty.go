package safe_live

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const md5_start_index = 5
const pw_length = 32

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
	f := "%s%0" + strconv.Itoa(pw_length-kl) + "d"
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

func EncryptPassword(pw string) (string, string, error) {
	now := time.Now().UnixNano()
	h := md5.New()
	h.Write([]byte(strconv.FormatInt(now, 10)))
	mm := hex.EncodeToString(h.Sum(nil))
	begin := mm[md5_start_index : pw_length-6-len(pw)+md5_start_index]
	npw := begin + pw
	p, err := encrypt(npw, fmt.Sprintf(fillFormat(pw), pw, 1))
	return npw, p, err
}

func DecryptPassword(pw string, line string) (string, error) {
	return decrypt(line, fmt.Sprintf(fillFormat(pw), pw, 1))
}
