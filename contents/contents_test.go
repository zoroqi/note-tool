package contents

import (
	"fmt"
	"testing"
)

func TestScanFiles(t *testing.T) {
	path := "E:\\doc\\Langeweile"
	tree, err := ScanFiles(path)
	if err != nil {
		t.Error(err)
		return
	}
	printTree(tree, " ")

}
func printTree(r *fileNode, head string) {
	fmt.Println(head, r.File.Name())
	if len(r.Child) > 0 {
		for _, c := range r.Child {
			printTree(c, head+head)
		}
	}
}
