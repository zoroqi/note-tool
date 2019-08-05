package main

import (
	"flag"
	"fmt"
	"github.com/zoroqi/note-tool/mind"
	"io/ioutil"
	"strings"
)

var handler = make(map[string]func(string, string, int) (string, error))

func init() {
	handler["2mm"] = md2mm
	handler["2md"] = mm2md
	handler["sp"] = special
}

func main() {
	f := flag.String("f", "", "change file path")
	o := flag.String("o", "", "output file")
	m := flag.String("m", "", "2mm/2md/sp(special input file)")
	b := flag.Int("b", 1, "2mm balance")
	indent := flag.String("i", "    ", "md indent")
	flag.Parse()
	if f == nil || *f == "" {
		fmt.Print("f is empty")
		return
	}
	if o == nil || *o == "" {
		fmt.Print("o is empty")
		return
	}
	if m == nil || handler[*m] == nil {
		fmt.Print("m is err, ", *m)
		return
	}
	bs, err := ioutil.ReadFile(*f)
	if err != nil {
		fmt.Print(err)
		return
	}
	str := string(bs)

	h := handler[*m]
	result, err := h(str, *indent, *b)
	if err != nil || result == "" {
		fmt.Print("change err, ", err)
		return
	}
	if err := ioutil.WriteFile(*o, []byte(result), 0666); err != nil {
		fmt.Print(err)
	}
}

func md2mm(str, indent string, balance int) (string, error) {
	md, err := mind.ParseMd(str, indent)
	if err != nil {
		fmt.Print("parse err", err)
		return "", err
	}
	doc := mind.Md2Mm(md, balance)
	doc.Indent(4)
	return doc.WriteToString()
}

func mm2md(str, indent string, balance int) (string, error) {
	return "", nil
}

// ...一些内容
//### 思维导图
//...
//
//> 脑图概述
//
//* 主题
//    * 主题2
//    * 主题3
//
//file end/### 其他内容
func special(str, indent string, balance int) (string, error) {

	startOffset := strings.Index(str, "### 思维导图")
	s := str[startOffset:]
	startOffset = strings.Index(s, "> 脑图概述")
	s = s[startOffset:]
	s = strings.Replace(s, "> 脑图概述", "", 1)
	s = strings.Trim(s, " ")

	linefeed := false
	end := 0
	for i, v := range s {
		if v == '\n' {
			linefeed = true
		} else {
			if linefeed && v == '#' {
				end = i
				break
			} else {
				linefeed = false
			}
		}
	}
	if end > 0 {
		s = s[0:end]
	}
	return md2mm(s, indent, balance)
}
