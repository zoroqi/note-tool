package ember

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseClippings(t *testing.T) {
	filePath := "/Users/wuming/note/My Clippings.txt"
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	book := ParseClippings(string(bs))
	for _, e := range book {
		fmt.Printf("%+v\n", e)
	}
}

func TestParseClippingsBlock(t *testing.T) {
	block := []string{
		"黑客与画家 (保罗·格雷厄姆)",
		"- 您在位置 #685的标注 | 添加于 2015年10月28日星期三 上午9:03:17",
		"",
		"程序写出来是给人看的，附带能在机器上运行。",
	}

	e := parseClippingBlock(block)
	fmt.Printf("%+v\n", e)

}

func TestBookCache_Format(t *testing.T) {
	filePath := "/Users/wuming/note/My Clippings.txt"
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	book := ParseClippings(string(bs))
	bc := BuildBookCache(book)
	bc.Format(25)
}
