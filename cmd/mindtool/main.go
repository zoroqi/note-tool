package main

import (
	"flag"
	"fmt"
	"github.com/zoroqi/note-tool/mind"
	"io/ioutil"
)

var handler = make(map[string]func(string, string, int) (string, error))

func init() {
	handler["2mm"] = md2mm
	handler["2md"] = mm2md
}

func main() {
	f := flag.String("f", "", "change file path")
	o := flag.String("o", "", "output file")
	m := flag.String("m", "", "2mm/2md")
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
	if err != nil || result == ""{
		fmt.Print("change err, ",err)
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
