package contents

import (
	"fmt"
	"testing"
)

func TestScanFiles(t *testing.T) {

	path := "F:\\doc\\github\\note"
	contentsRemove := ContentsRemove("contents.md")

	remove := func(s string) bool {
		return HiddenRemove(s) || ImgRemove(s) || contentsRemove(s)
	}

	tree, err := ScanFiles(path, remove)
	if err != nil {
		t.Fatal(err)
	}
	printTree(tree," ")
}

func printTree(r *fileNode, head string) {
	fmt.Println(head, r.File.Name())
	if len(r.Child) > 0 {
		for _, c := range r.Child {
			printTree(c, head+head)
		}
	}
}
