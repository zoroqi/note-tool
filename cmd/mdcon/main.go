package main

import (
	"flag"
	"fmt"
	"github.com/zoroqi/note-tool/mdcon"
	"io/ioutil"
)

func main() {
	file := flag.String("f", "", "markdown file")
	flag.Parse()
	if *file == "" {
		fmt.Println("f is empty")
		return
	}

	bs, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("read file error%v", err)
		return
	}
	md := string(bs)
	if md == "" {
		return
	}

	mdContent := mdcon.SetContents(md)
	err = ioutil.WriteFile(*file, []byte(mdContent), 0666)
	if err != nil {
		fmt.Println("write err%v", err)
	}
}
