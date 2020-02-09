package ember

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type clipping struct {
	id           int
	book         string
	offsetStart  int
	offsetEnd    int
	date         string
	time         time.Time
	text         string
	clippingType string
}

var dateFindRegex *regexp.Regexp
var offsetRegex *regexp.Regexp
var typeRegex *regexp.Regexp

func init() {
	date, _ := regexp.Compile("添加于\\s(.*)")
	offset, _ := regexp.Compile("位置 #(.*?)([）]*?|\\s*?)的")
	ctype, _ := regexp.Compile("位置 .*的(.*?) ")
	dateFindRegex = date
	offsetRegex = offset
	typeRegex = ctype
}

type parseIterator struct {
	lines  []string
	index  int
	length int
}

func (c *parseIterator) hasNext() bool {
	if c.index >= c.length {
		return false
	}
	return true
}

func (c *parseIterator) next() *clipping {
	start := c.index
	end := c.index
	for i := c.index; i < c.length; i++ {
		if strings.HasPrefix(c.lines[i], "=========") {
			end = i
			break
		}
	}
	c.index = end + 1
	return parseClippingBlock(c.lines[start:end])
}

func parseClippingBlock(block []string) *clipping {
	l := len(block)
	e := &clipping{}
	if l == 0 {
		return e
	}

	if l >= 1 {
		book := strings.Trim(block[0], string('\uFEFF'))
		e.book = strings.TrimSpace(book)
	}

	if l >= 2 {
		t := block[1]
		date := dateFindRegex.FindStringSubmatch(t)
		if len(date) > 0 {
			d := date[1]
			e.date = d
			e.time = parseTime(d)
		}

		offset := offsetRegex.FindStringSubmatch(t)
		if len(offset) > 0 {
			o := offset[1]
			os := strings.SplitN(o, "-", 2)
			start, _ := strconv.ParseInt(os[0], 10, 32)
			end := start
			if len(os) > 1 {
				end, _ = strconv.ParseInt(os[1], 10, 32)
			}
			e.offsetStart = int(start)
			e.offsetEnd = int(end)
		}

		ctype := typeRegex.FindStringSubmatch(t)
		if len(ctype) > 0 {
			e.clippingType = ctype[1]
		}
	}
	if l >= 4 {
		e.text = block[3]
	}
	return e

}

// 解析中文版蛋疼的时间格式, 不知道其他版本. 更加蛋疼的golang的时间解析
//  中文的时间格式 2015年10月28日星期三 上午8:40:11
//  对应format yyyy年M月d日EEE ah:m:s
func parseTime(date string) time.Time {
	regex, _ := regexp.Compile("(\\d{1,4})年(\\d{1,2})月(\\d{1,2})日.* (.*?)(\\d{1,2}):(\\d{1,2}):(\\d{1,2})")

	matchs := regex.FindStringSubmatch(date)
	r := "am"
	if matchs[4] == "下午" {
		r = "pm"
	}

	timeText := fmt.Sprintf("%s-%s-%s %s%s:%s:%s", matchs[1], matchs[2], matchs[3], r, matchs[5], matchs[6], matchs[7])

	const timeFormat = "2006-1-2 pm3:4:5"

	t, _ := time.Parse(timeFormat, timeText)
	return t
}

func ParseClippings(clippingsText string) []Book {

	lines := strings.Split(strings.ReplaceAll(clippingsText, "\r", ""), "\n")

	iterator := &parseIterator{lines: lines, index: 0, length: len(lines)}
	clippings := make([]*clipping, 0, len(clippingsText)/100)
	for iterator.hasNext() {
		ec := iterator.next()
		clippings = append(clippings, ec)
	}

	books := make([]Book, 0, len(clippings))
	notes := make([]*clipping, 0, len(clippings)/50)
	offsetMapping := make(map[string]map[int]int)

	for i, c := range clippings {
		switch c.clippingType {
		case "标注":
			books = append(books, Book{Id: i, Name: c.book, Date: c.date, Time: c.time, Text: c.text, Note: ""})
			if _, exist := offsetMapping[c.book]; !exist {
				offsetMapping[c.book] = make(map[int]int)
			}
			offsetMapping[c.book][c.offsetEnd] = i
		case "笔记":
			notes = append(notes, c)
		}
	}

	for _, n := range notes {
		if o, exist := offsetMapping[n.book]; exist {
			id := o[n.offsetEnd]
			if id > 0 {
				index := sort.Search(len(books), func(i int) bool {
					return books[i].Id >= id
				})
				if index >= 0 && index < len(books) {
					if books[index].Id == id {
						books[index].Note = n.text
					}
				}
			}
		}
	}

	return books
}
