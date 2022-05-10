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

func TestParseKfxBlock(t *testing.T) {
	block := []string{
		"弱者的武器：农民反抗的日常形式 (人文与社会译丛) ([美国]詹姆斯•C•斯科特)",
		"- 您在第 91 页（位置 #1347）的笔记 | 添加于 2022年5月10日星期二 上午10:15:36",
		"",
		"我们以1966年作为比较的起点，穆达地区绝大多数的农民家庭确实比以前富裕了很多，但与此同时，收入分配差距也越来越大，而且，相当多的农民——可能要占总人口的35%—40%—已经被远远地抛在了后面，他们的收入和10年前差不多。",
	}

	e := parseClippingBlock(block)
	fmt.Println(e.id)
	fmt.Println(e.text)
	fmt.Println(e.offsetEnd)
	fmt.Println(e.offsetStart)
	fmt.Println(e.date)
	fmt.Println(e.clippingType)

}
