package main

import (
	"flag"
	"fmt"
	safe_live "github.com/zoroqi/note-tool/safe-live"
	"io/ioutil"
	"strings"
)

func main() {
	f := flag.String("f", "", "change file path")
	o := flag.String("o", "", "output file")
	m := flag.String("m", "", "mode en/des")
	p := flag.String("p", "", "password")
	flag.Parse()
	if f == nil || *f == "" {
		fmt.Print("f is empty")
		return
	}
	if o == nil || *o == "" {
		fmt.Print("o is empty")
		return
	}

	if m == nil || *m == "" {
		fmt.Print("mode is empty")
		return
	}

	if p == nil || *p == "" || len(*p) > 12 {
		fmt.Print("mode is empty or len(p) > 12")
		return
	}
	bs, err := ioutil.ReadFile(*f)
	if err != nil {
		fmt.Print(err)
		return
	}

	str := string(bs)

	md := strings.ReplaceAll(str, "\r", "")
	lines := strings.Split(md, "\n")
	if *m == "en" {
		lines, err = safe_live.EncryptLines(lines, *p)
	} else if *m == "des" {
		lines, err = safe_live.DecryptLines(lines, *p)
	} else {
		fmt.Println("mode error")
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	sb := strings.Builder{}
	for _, v := range lines {
		sb.WriteString(v)
		sb.WriteString("\n")
	}

	if err := ioutil.WriteFile(*o, []byte(sb.String()), 0666); err != nil {
		fmt.Print(err)
	}
}
