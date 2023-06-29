package safe_live

import (
	"testing"
)

func Test_newkey(t *testing.T) {
	str := "123"
	key := "1234567890123456"
	nkey := newkey(str, key)
	if nkey != "OQ[9#(WU?:U6(+g4<e~dXXhec<;0*hZ8" {
		t.Errorf("nkey is %s", nkey)
	}
}

func TestEncrypt(t *testing.T) {
	str := "123"
	key := "1234567890123456"
	en, err := encrypt(str, key)
	if err != nil {
		t.Error(err)
	}
	if en != "yGjLrkpCmzyuVSHt7ZUflw==" {
		t.Errorf("en is %s", en)
	}
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
	lines2, _ := DecryptLines(lines, key)
	for i, v := range lines2 {
		if block[i] != v {
			t.Errorf("%s%s", block[i], v)
		}
	}
}

func TestPassword(t *testing.T) {
	pw := "1234567890123456"
	nkey, enkey, err := EncryptPassword(pw)
	if err != nil {
		t.Error(err)
	}
	decPw, err := DecryptPassword(pw, enkey)
	if err != nil {
		t.Error(err)
	}
	if decPw != nkey {
		t.Errorf("decPw is %s", decPw)
	}
}
