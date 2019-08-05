package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/zoroqi/note-tool/ember"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	filePath := flag.String("f", "", "kindle 'My Clippings.txt' path")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("file path is empty")
		return
	}
	bs, err := ioutil.ReadFile(*filePath)
	if err != nil {
		fmt.Println("-f %s, %v", *filePath, err)
		return
	}
	books := ember.ParseClippings(string(bs))
	bookCache := ember.BuildBookCache(books)

	help := `
Can only appear once
b list all books
f find book name
s list specify book by id
`
	for true {
		fmt.Println(help)
		input := strings.TrimSpace(readConsoleLine())
		if "" == input {
			continue
		}

		param := strings.SplitN(input, " ", 2)

		switch param[0] {
		case "b":
			bookCache.ListAllBooks()
		case "f":
			if len(param) <= 1 {
				fmt.Println("input keyword")
				continue
			}
			bookCache.FindBookByWord(param[1])
		case "s":
			{
				if len(param) <= 1 {
					fmt.Println("input book id")
					continue
				}
				i, err := strconv.ParseInt(param[1], 10, 32)
				if err != nil {
					fmt.Println("id is not number")
				}
				bookCache.Format(int(i))
			}
		default:
			fmt.Println("error param")

		}

	}
}

func readConsoleLine() string {
	reader := bufio.NewReader(os.Stdin)
	data, _, e := reader.ReadLine()
	if e != nil {
		return ""
	}
	regexStr := string(data)
	return regexStr
}
