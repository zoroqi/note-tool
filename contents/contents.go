package contents

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Predicate func(string) bool

type PreText func(int) string

type Config struct {
	Name   string    // contents file name
	Level  int       // file level
	Remove Predicate // remove config file
	Space  PreText   // levelSpace
}

type fileNode struct {
	FilePath string
	File     os.FileInfo
	Child    []*fileNode
	Parent   *fileNode
}

type FileLink struct {
	Name  string
	Link  string
	Level int
}

func (f FileLink) buildLink(rootLength int) string {
	return fmt.Sprintf("[%s](./%s)", f.Name, strings.ReplaceAll(f.Link[rootLength:], " ", "%20"))
}

func CreateContents(root string, config Config) {
	tree, err := ScanFiles(root, config.Remove)
	if err != nil {
		fmt.Println(err)
		return
	}
	createContents(tree, 0, config)
}

func createContents(node *fileNode, level int, c Config) {
	if level >= c.Level {
		return
	}
	if node.File.IsDir() {
		contents := buildContents(node, c)
		if err := ioutil.WriteFile(node.FilePath+"/"+c.Name, []byte(contents), 0666); err != nil {
			fmt.Print(err)
		}
		for _, child := range node.Child {
			createContents(child, level+1, c)
		}
	}

}

func buildContents(node *fileNode, config Config) string {
	fl := buildText(node, 0, config)
	sb := strings.Builder{}
	sb.WriteString("# contents\n\n")
	rootFileNameLength := len(node.File.Name())
	for i, l := range fl {
		if i == 0 {
			l.Link = l.Link + "/" + config.Name
		}
		sb.WriteString(fmt.Sprintf("%s- %s\n", config.Space(l.Level), l.buildLink(rootFileNameLength+1)))
	}
	return sb.String()
}

func buildText(node *fileNode, level int, c Config) []FileLink {
	if level >= c.Level {
		return nil
	}
	if node.File.IsDir() {
		lines := make([]FileLink, 0, len(node.Child))
		lines = append(lines, FileLink{Name: node.File.Name(), Link: node.File.Name(), Level: level})
		for _, child := range node.Child {
			cs := buildText(child, level+1, c)
			for _, cc := range cs {
				fl := cc
				fl.Link = node.File.Name() + "/" + fl.Link
				lines = append(lines, fl)
			}
		}
		return lines
	} else {
		return []FileLink{FileLink{Name: node.File.Name(), Link: node.File.Name(), Level: level}}
	}
	return nil
}

func ScanFiles(path string, remove Predicate) (*fileNode, error) {
	root := &fileNode{File: nil, Child: make([]*fileNode, 0), Parent: nil}
	err := scanFiles(path, root, remove)
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

func scanFiles(path string, pn *fileNode, remove Predicate) error {
	if cs, err := os.Lstat(path); err == nil {
		node := &fileNode{
			FilePath: path,
			File:     cs,
			Child:    make([]*fileNode, 0),
			Parent:   pn,
		}
		pn.Child = append(pn.Child, node)
		if cs.IsDir() {
			childNames, err := readDirNames(path)
			if err != nil {
				return err
			}
			for _, name := range childNames {
				if remove(name) {
					continue
				}
				cp := path + "/" + name
				if err := scanFiles(cp, node, remove); err != nil {
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
