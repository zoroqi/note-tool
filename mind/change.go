package mind

import (
	"github.com/beevik/etree"
	"strings"
)

const (
	POSITION    = "POSITION"
	right       = "right"
	left        = "left"
	version     = "version"
	V     = "1.0.1"
	CREATED     = "CREATED"
	MODIFIED    = "MODIFIED"
	ID          = "ID"
	ID_         = "ID_"
	TEXT        = "TEXT"
	LINK        = "LINK"
	richcontent = "richcontent"
	TYPE        = "TYPE"
	NOTE        = "NOTE"
	html        = "html"
	body        = "body"
	p           = "p"
	NODE        = "node"
	MAP         = "map"
)

func Md2Mm(tree *MdTree, balanceLevel int) *etree.Document {
	doc := etree.NewDocument()
	md2mm(tree, tree.Root.children[0], doc.CreateElement(MAP))
	mind := doc.SelectElement(MAP)
	mind.CreateAttr(version, V)
	root := mind.ChildElements()
	balance(tree, root[0], balanceLevel)
	return doc
}

func balance(tree *MdTree, element *etree.Element, balance int) {
	cs := element.ChildElements()
	if cs == nil {
		return
	}
	if balance > 0 {
		count := make(map[string]int)
		total := 0
		for _, v := range cs {
			s := attrFind(v.Attr, CREATED)
			if s == "" {
				continue
			}
			count[s] = tree.idMapper[s].descendantCountByLevel(balance)
			total += count[s]
		}
		b := total >> 1
		for _, v := range cs {
			if total > b {
				v.CreateAttr(POSITION, right)
			} else {
				v.CreateAttr(POSITION, left)
			}
			s := attrFind(v.Attr, CREATED)
			total = total - count[s]
		}
	} else {
		b := len(cs) / 2
		for i, v := range cs {
			if i <= b {
				v.CreateAttr(POSITION, right)
			} else {
				v.CreateAttr(POSITION, left)
			}
		}
	}
}

func md2mm(tree *MdTree, node *MdNode, doc *etree.Element) {
	if node == nil {
		return
	}
	n := doc.CreateElement(NODE)
	n.CreateAttr(CREATED, node.id)
	n.CreateAttr(ID, ID_+node.id)
	n.CreateAttr(MODIFIED, node.id)
	n.CreateAttr(TEXT, node.subject)
	if node.link != "" {
		n.CreateAttr(LINK, node.link)
	}

	if node.richContent != "" {
		rc := n.CreateElement(richcontent)
		rc.CreateAttr(TYPE, NOTE)
		b := rc.CreateElement(html).CreateElement(body)
		for _, v := range strings.Split(node.richContent, "\n") {
			if v != "" {
				p := b.CreateElement(p)
				p.SetText(v)
			}
		}
	}
	for _, v := range node.children {
		md2mm(tree, v, n)
	}
}

func attrFind(attr []etree.Attr, key string) string {
	for _, v := range attr {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
