package safe_live

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	str := "123"
	key := "1234567890123456"
	en, err := encrypt(str, key)
	fmt.Println(en, err)
	fmt.Println(decrypt(en, key))
}

func TestEncryptLines(t *testing.T) {
	block := []string{
		"黑客与画家 (保罗·格雷厄姆)",
		"- 您在位置 #685的标注 | 添加于 2015年10月28日星期三 上午9:03:17",
		"",
		"程序写出来是给人看的，附带能在机器上运行。",
	}
	key := "1234567890a"
	lines, _ := EncryptLines(block, key)

	for _, v := range lines {
		fmt.Println(v)
	}
	fmt.Println()

	lines2, _ := DecryptLines(block, key)
	for _, v := range lines2 {
		fmt.Println(v)
	}
}

func TestFillNum(t *testing.T) {
	fmt.Printf("%04d", 34124)
}

func TestFillFormat(t *testing.T) {
	str := "freedomdie"
	fmt.Printf(fillFormat(str), str, 1)
}

func TestPassword(t *testing.T) {
	str := "1234567890123456"
	npw, _ := EncryptPassword(str)
	fmt.Println(npw)
	fmt.Println(DecryptPassword(str, npw))
}
