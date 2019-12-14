package main

import (
	"flag"
	"github.com/zoroqi/note-tool/contents"
)

func main() {
	name := flag.String("n", "contents.md", "contents name")
	level := flag.Int("l", 3, "output file level")

	contentsRemove := contents.ContentsRemove(*name)

	remove := func(s string) bool {
		return contents.HiddenRemove(s) || contents.ImgRemove(s) || contentsRemove(s)
	}

	config := contents.Config{Name: *name,
		Level:  *level,
		Remove: remove,
		Space:  contents.PreSpace}

	contents.CreateContents("./", config)
}
