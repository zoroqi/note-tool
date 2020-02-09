package ember

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseClippings(t *testing.T) {
	filePath := "/Users/wuming/cache/k.txt"
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
		"- 您在第 135 页（位置 #1925）的笔记 | 添加于 2020年1月13日星期一 下午9:38:05^M",
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
