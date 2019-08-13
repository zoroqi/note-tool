package mind

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MdNode struct {
	id              string    // 序列id
	subject         string    // 主题
	children        []*MdNode // 孩子
	richContent     string    // 其他描述
	parent          *MdNode   // 父节点
	link            string    // 链接
	depth           int       // 当前节点深度
	childCount      int       // 子孩节点数量
	descendantCount int       // 子孙节点数量
}

func (m *MdNode) AddChild(c *MdNode) {
	m.children = append(m.children, c)
	c.parent = m
	m.childCount++
	m.descendantCount++
	n := m.parent
	for n != nil {
		n.descendantCount++
		n = n.parent
	}
}

func (m *MdNode) descendantCountByLevel(level int) int {
	if len(m.children) == 0 {
		return 0
	}
	q := make([]*MdNode, 0, 20)
	q = append(q, m.children...)
	c := 0
	maxDepth := m.depth + level
	for len(q) > 0 {
		n := q[0]
		if n.depth <= maxDepth {
			q = append(q, n.children...)
			c++
		}
		q = q[1:]
	}

	return c
}

type MdTree struct {
	Root          *MdNode
	Topic         string
	idMapper      map[string]*MdNode
	subjectMapper map[string][]*MdNode
}

func (t *MdTree) Print() {
	if t.Root == nil {
		return
	}
	for _, v := range t.Root.children {
		DF(v, "", "    ")
	}
}
func (t *MdTree) AddNode(node *MdNode) {
	t.idMapper[node.id] = node
	if t.subjectMapper[node.subject] == nil {
		t.subjectMapper[node.subject] = make([]*MdNode, 0, 1)
	}
	t.subjectMapper[node.subject] = append(t.subjectMapper[node.subject], node)
}

func DF(root *MdNode, prefix, indent string) {
	if root == nil {
		return
	}
	s := prefix + indent
	r := root.richContent
	if root.link != "" {
		fmt.Printf("%s* [%s](%s)\n", prefix, root.subject, root.link)
	} else {
		fmt.Printf("%s* %s %s\n", prefix, root.subject, root.link)
	}
	if r != "" {
		rl := strings.Split(r, "\n")
		for _, v := range rl {
			if v != "" {
				fmt.Printf("%s%s\n", s, v)
			}
		}
	}
	for _, v := range root.children {
		DF(v, s, indent)
	}
}

func ParseMd(markdown string, indent string) (tree *MdTree, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	err = nil

	md := strings.ReplaceAll(markdown, "\r", "")
	lines := strings.Split(md, "\n")
	topic := ""
	for _, v := range lines {
		if strings.TrimSpace(v) != "" {
			topic = strings.TrimSpace(v)[2:]
			break
		}
	}

	tree = &MdTree{Topic: topic, idMapper: make(map[string]*MdNode), subjectMapper: make(map[string][]*MdNode)}
	root := newMdNode(topic, "\\0", -1)
	tree.AddNode(root)
	tree.Root = root

	stackOffset := 0
	parentStack := make([]*MdNode, 30, 40)
	parentStack = append(parentStack, root)
	parentStack[stackOffset] = root
	nodeId := 0
outer:
	for _, v := range lines {
		if v == "" || strings.Trim(v, indent) == "" {
			continue
		}
		d := depth(v, indent)
		subject := trimIndent(v, indent)
		var node *MdNode
		if isSubject(subject) {
			nodeId++
			subject = subject[2:]
			node = newMdNode(strconv.Itoa(nodeId), subject, d)
			pd := parentStack[stackOffset].depth
			parent := parentStack[stackOffset]
			if d <= pd {
				// 相同和父节点, 需要寻找最近的父节点
				for {
					stackOffset--
					if stackOffset < 0 {
						break outer
					}
					if parentStack[stackOffset].depth < d {
						parent = parentStack[stackOffset]
						stackOffset++
						parentStack[stackOffset] = node
						break
					}
				}
			} else if d > pd {
				// 变深, 直接加一
				stackOffset++
				parentStack[stackOffset] = node
			}
			// 修改id的编码方式, 改为x.x.x这种方式
			parent.AddChild(node)
			node.id = parent.id + "." + strconv.Itoa(parent.childCount)
			tree.AddNode(node)
		} else {
			parentStack[stackOffset].richContent += subject + "\n"
		}
	}

	return tree, err
}

func isSubject(str string) bool {
	return strings.HasPrefix(str, "*")
}

func depth(str string, indent string) int {
	inLen := len(indent)
	l := len(str)/inLen + 1
	for i := 0; i < l; i++ {
		if str[i*inLen:(i+1)*inLen] != indent {
			return i
		}
	}
	return 0
}

func trimIndent(str string, indent string) string {
	return strings.TrimLeft(str, indent)
}

// 进行, 将subject转换成对应数据
// 需要转换link内容
func newMdNode(id string, subject string, depth int) *MdNode {
	regex, _ := regexp.Compile("\\[(.*?)\\]\\((.*?)\\)")
	s := regex.FindStringSubmatch(subject)
	sub := subject
	link := ""
	if s != nil {
		sub = s[1]
		link = s[2]
	}
	return &MdNode{id: id, subject: sub, children: make([]*MdNode, 0), depth: depth, link: link}
}
