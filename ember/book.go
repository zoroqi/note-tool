package ember

import (
	"fmt"
	"strings"
	"time"
)

const time_layout = "2006-01-02"

type Book struct {
	Id   int
	Name string
	Time time.Time
	Date string
	Text string
	Note string
}

type BookCache struct {
	BookArray []Book
	BookNames map[string][]Book
	BookIds   []string
}

func (es *BookCache) ListAllBooks() {
	for i, e := range es.BookIds {
		fmt.Println(i, e)
	}
}

func (es *BookCache) ListBookById(id int) ([]Book, bool) {
	if id > len(es.BookIds) {
		return nil, false
	}

	return es.ListBook(es.BookIds[id])
}

func (es *BookCache) ListBook(bookName string) ([]Book, bool) {
	books, exist := es.BookNames[bookName]
	if !exist {
		return nil, false
	}
	return books, true
}

func (es *BookCache) FindBookByWord(keyword string) []int {
	books := make([]int, 0, 10)
	for i, e := range es.BookIds {
		if strings.Contains(e, keyword) {
			fmt.Println(i, e)
			books = append(books, i)
		}
	}
	return books
}

// 格式化输出文本
//  layout 简单替换
//  -t 标注内容
//  -b 书名
//  -c 评论/标注
//  -t 时间
//  -i 编号, 固定顺序自增
//  针对15天以上的标注间隔, 插入一条时间范围说明
func (es *BookCache) Format(bookId int) string {
	books, exist := es.ListBookById(bookId)
	if !exist {
		return ""
	}
	sa := make([]string, 0, len(books))
	startReadIndex := 0
	saIndex := 0
	currentTime := books[0].Time
	sa = append(sa, "")
	const interval = 86400 * 15
	for i, v := range books {
		if (v.Time.Unix() - currentTime.Unix()) >= interval {
			sa[saIndex] = dateFormat(books[startReadIndex], books[i-1])
			startReadIndex = i
			sa = append(sa, "")
			saIndex = len(sa) - 1
		}
		sa = append(sa, format(i+1, v))
		currentTime = v.Time
	}
	sa[saIndex] = dateFormat(books[startReadIndex], books[len(books)-1])

	sb := strings.Builder{}
	sb.WriteString(books[0].Name)
	sb.WriteString("\n")
	for _, v := range sa {
		sb.WriteString(v)
	}
	fmt.Println(replaceCp(sb.String()))
	return replaceCp(sb.String())
}

func format(i int, e Book) string {
	str := fmt.Sprintf("%d. %s\n", i, e.Text)
	if e.Note != "" {
		str += fmt.Sprintf("    * %s\n", e.Note)
	}
	return str
}

func dateFormat(e1, e2 Book) string {
	return fmt.Sprintf("* %s ~ %s\n", e1.Time.Format(time_layout), e2.Time.Format(time_layout))
}

func BuildBookCache(es []Book) *BookCache {
	result := BookCache{BookArray: es}
	books := make(map[string][]Book)
	bookIds := make([]string, 0, len(es)/20)
	id := 1
	for _, e := range es {
		if _, exist := books[e.Name]; !exist {
			books[e.Name] = make([]Book, 0, 20)
			bookIds = append(bookIds, e.Name)
			id++
		}
		books[e.Name] = append(books[e.Name], e)
	}

	result.BookNames = books
	result.BookIds = bookIds
	return &result
}

// 替换一部分中文标点, 这个纯粹是个人习惯. 部分中文标点会被替换成英文标点.
func replaceCp(str string) string {

	cp := make(map[rune]rune)
	cp['，'] = ','
	cp['。'] = '.'
	cp['“'] = '"'
	cp['”'] = '"'
	cp['’'] = '\''
	cp['‘'] = '\''
	cp['《'] = '<'
	cp['》'] = '>'
	cp['（'] = '('
	cp['）'] = ')'
	cp['　'] = ' '
	cp['-'] = '-'
	cp['；'] = ';'
	cp['：'] = ':'
	cp['？'] = '?'
	cp['！'] = '!'
	cp['、'] = ','
	defaultMap := func(r rune) rune {
		if v, e := cp[r]; e {
			return v
		}
		return r
	}
	rs := []rune(str)

	for i, v := range rs {
		rs[i] = defaultMap(v)
	}

	return string(rs)
}
