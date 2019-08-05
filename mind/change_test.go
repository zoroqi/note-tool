package mind

import (
	"fmt"
	"testing"
)

func TestMd2Mm(t *testing.T) {
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
	md, _ := ParseMd(str, "    ")
	doc := Md2Mm(md, 1)
	doc.Indent(4)
	doc.Indent(2)
	xml, _ := doc.WriteToString()
	fmt.Println(xml)
}
