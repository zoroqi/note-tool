package main

import (
	"flag"
	"fmt"
)

func main()  {
	file := flag.String("f", "", "markdown file")
	flag.Parse()
	if *file == "" {
		fmt.Errorf("f is empty")
	}


}
