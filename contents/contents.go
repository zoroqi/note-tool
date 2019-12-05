package contents

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Config struct {
	Name     string // contents file name
	Relative bool   // the link is relative
	Level    int    // file level
}

type fileNode struct {
	File   os.FileInfo
	Child  []*fileNode
	Parent *fileNode
}

func BuildContents(root string, config Config) {
	tree, err := ScanFiles(root)
	if err != nil {
		fmt.Println(err)
		return
	}
	buildContent(tree, config)
}

func buildContent(node *fileNode, config Config) {

}

func buildText(node *fileNode, config Config) string {
	if node.File.IsDir() {

	} else {

	}
	return ""
}

func ScanFiles(path string) (*fileNode, error) {
	root := &fileNode{File: nil, Child: make([]*fileNode, 0), Parent: nil}
	err := scanFiles(path, root)
	if err != nil {
		return nil, err
	}
	if len(root.Child) > 0 {
		root.Child[0].Parent = nil
		return root.Child[0], nil
	} else {
		return nil, errors.New("no child")
	}
}

func scanFiles(path string, pn *fileNode) error {
	if cs, err := os.Lstat(path); err == nil {
		node := &fileNode{
			File:   cs,
			Child:  make([]*fileNode, 0),
			Parent: pn,
		}
		pn.Child = append(pn.Child, node)
		if cs.IsDir() {
			childNames, err := readDirNames(path)
			if err != nil {
				return err
			}
			for _, name := range childNames {
				if strings.HasPrefix(name, ".") {
					continue
				}
				cp := path + "/" + name
				if err := scanFiles(cp, node); err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return err
	}
}

// 蛋疼的读取子文件的方法
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
