package mind

import (
	"fmt"
	"testing"
)

func TestParseMd(t *testing.T) {
	str := `
* 根
    * 标1
        * [标1.1](http://www.baidu.com)
            * 标1.1.1
                1. 1.1.1.1
                2. 1.1.1.2
                3. 1.1.1.3
            * 标1.1.2
    * 标2
    * 标3
        1. 3.1.1.1
        2. 3.1.1.2
        * 标3.1
        * 标3.2
            1. 3.2.1.1
    * 标4
    * 标5
        * 标6
            * 标7
                * 标8
                    * 标9
    * 标10
`
	tree, err := ParseMd(str, "    ")
	fmt.Println(err)
	tree.Print()
	node6 := tree.subjectMapper["标6"][0]
	node7 := tree.subjectMapper["标7"][0]
	node8 := tree.subjectMapper["标8"][0]
	if node6 != node7.parent {
		t.Error("parent error")
	}
	if node7 != node8.parent {
		t.Error("parent error")
	}
	if node6.children[0] != node7 {
		t.Error("child error")
	}
	if node7.children[0] != node8 {
		t.Error("child error")
	}
	node3 := tree.subjectMapper["标3"][0]
	node32 := tree.subjectMapper["标3.2"][0]
	if node3.children[1] != node32 {
		t.Error("child error")
	}

	if node3 != node32.parent {
		t.Error("parent error")
	}
	fmt.Println(tree.idMapper)

}
