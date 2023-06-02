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

	if p == nil || *p == "" || len(*p) > 16 {
		fmt.Print("mode is empty or len(p) > 16")
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
		password, headLine, err := safe_live.EncryptPassword(*p)
		if err != nil {
			if err != nil {
				fmt.Println("encrypt pw error,", err)
				return
			}
		}
		lines, err = safe_live.EncryptLines(lines, password)
		if err != nil {
			fmt.Println("encrypt error,", err)
			return
		}
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, headLine)
		newLines = append(newLines, lines...)
		lines = newLines
	} else if *m == "des" {
		password, err := safe_live.DecryptPassword(*p, lines[0])
		if err != nil {
			fmt.Println("decrypt pw error,", err)
			return
		}
		lines, err = safe_live.DecryptLines(lines[1:], password)
		if err != nil {
			fmt.Println("decrypt error,", err)
			return
		}
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
	all := sb.String()
	if err := ioutil.WriteFile(*o, []byte(all[0:len(all)-1]), 0644); err != nil {
		fmt.Print(err)
	}
}
