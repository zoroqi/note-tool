package mdcon

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestBuildMdContents(t *testing.T)  {
	path := "test.md"

	bs,err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	md := BuildMdContents(string(bs))
	fmt.Println(md)
}



func TestFindTitles(t *testing.T) {

	path := "test.md"

	bs,err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(bs),"\n")
	titles :=  findTitle(lines)
	for _,t := range titles {
		fmt.Println(t)
	}
}

func TestParseTitle(t *testing.T) {
	title := parseTitle("### test")
	if title.level != 3 {
		t.Error("level is not 3")
	}
	if title.text != "test" {
		t.Error("text is not text")
	}
	fmt.Printf("%+v",title)
}

func TestBuildContents(t *testing.T) {
	lines := []string{"# test","## test2", "### test3", "##### test4","#### test5","## test6"}

	r := "-----\n* 目录\n- [test](#test)\n    - [test2](#test2)\n        - [test3](#test3)\n" +
		"                - [test4](#test4)\n            - [test5](#test5)\n    - [test6](#test6)\n-----\n"

	contents := buildContents(lines)
	if r != contents {
		t.Error("build error")
	}
	fmt.Println(contents)

}
