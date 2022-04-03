package ember

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ClippingType string
type BookDesc string

const ClippingMark ClippingType = "mark"
const ClippingNote ClippingType = "note"
const ClippingBookmark ClippingType = "bookmark"

const BookMobiOrAzw BookDesc = "mobiOrAzw"
const BookPdf BookDesc = "pdf"

type clipping struct {
	id           int
	book         string
	offsetStart  int
	offsetEnd    int
	date         string
	time         time.Time
	text         string
	clippingType ClippingType
}

var dateFindRegex = regexp.MustCompile("添加于\\s(.*)")

// - 您在第 126-126 页的标注 | 添加于 2022年3月28日星期一 上午8:35:18
// - 您在第 57 页的笔记 | 添加于 2022年3月25日星期五 下午7:20:52
// - 您在位置 #1680 的书签 | 添加于 2016年1月27日星期三 上午9:08:40
// - 您在位置 #1103-1103的标注 | 添加于 2016年1月26日星期二 下午8:05:45
// - 您在位置 #1385 的笔记 | 添加于 2015年12月16日星期三 下午7:55:38

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

	// - 您在第 126-126 页的标注 | 添加于 2022年3月28日星期一 上午8:35:18
	// - 您在第 57 页的笔记 | 添加于 2022年3月25日星期五 下午7:20:52
	// - 您在位置 #1680 的书签 | 添加于 2016年1月27日星期三 上午9:08:40
	// - 您在位置 #1103-1103的标注 | 添加于 2016年1月26日星期二 下午8:05:45
	// - 您在位置 #1385 的笔记 | 添加于 2015年12月16日星期三 下午7:55:38
	if l >= 2 {
		t := block[1]
		filedesc := bookDesc(t)
		start, end, ct := filedesc.indexParse(t)
		date := dateFindRegex.FindStringSubmatch(t)
		if len(date) > 0 {
			d := date[1]
			e.date = d
			e.time = parseTime(d)
		}
		e.offsetStart = int(start)
		e.offsetEnd = int(end)
		e.clippingType = ct
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
		case ClippingMark:
			books = append(books, Book{Id: i, Name: c.book, Date: c.date, Time: c.time, Text: c.text, Note: ""})
			if _, exist := offsetMapping[c.book]; !exist {
				offsetMapping[c.book] = make(map[int]int)
			}
			offsetMapping[c.book][c.offsetEnd] = i
		case ClippingNote:
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

// - 您在第 57 页的笔记 | 添加于 2022年3月25日星期五 下午7:20:52 -> pdf
// - 您在位置 #1680 的书签 | 添加于 2016年1月27日星期三 上午9:08:40 -> mobi
func bookDesc(s string) BookDesc {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	if strings.HasPrefix(s, "您在第") {
		return BookPdf
	}
	return BookMobiOrAzw
}

func (c BookDesc) indexParse(s string) (start, end int64, t ClippingType) {
	switch c {
	case BookPdf:
		return pdfIndexParse(s)
	default:
		return mobiIndexParse(s)
	}
}

// - 您在位置 #1680 的书签 | 添加于 2016年1月27日星期三 上午9:08:40
// - 您在位置 #1103-1103的标注 | 添加于 2016年1月26日星期二 下午8:05:45
var mobiOffsetRegex = regexp.MustCompile("位置 #(.*?)([）]*?|\\s*?)的")
var mobiTypeRegex = regexp.MustCompile("位置 .*的(.*?) ")

func mobiIndexParse(s string) (start, end int64, t ClippingType) {
	offset := mobiOffsetRegex.FindStringSubmatch(s)
	if len(offset) > 0 {
		start, end = offsetHandler(offset)
	}
	t = ClippingMark
	ctype := mobiTypeRegex.FindStringSubmatch(s)
	if len(ctype) > 0 {
		if ctype[1] == "标注" {
			t = ClippingMark
		} else if ctype[1] == "笔记" {
			t = ClippingNote
		} else {
			t = ClippingBookmark
		}
	}
	return
}

func offsetHandler(offset []string) (start, end int64) {
	o := offset[1]
	os := strings.SplitN(o, "-", 2)
	start, _ = strconv.ParseInt(os[0], 10, 32)
	end = start
	if len(os) > 1 {
		end, _ = strconv.ParseInt(os[1], 10, 32)
	}
	return start, end
}

// - 您在第 57 页的笔记 | 添加于 2022年3月25日星期五 下午7:20:52
// - 您在第 126-126 页的标注 | 添加于 2022年3月28日星期一 上午8:35:18
var pdfOffsetRegex = regexp.MustCompile("第 (.*?)([）]*?|\\s*?)页的")
var pdfTypeRegex = regexp.MustCompile("第 .*页的(.*?) ")

func pdfIndexParse(s string) (start, end int64, t ClippingType) {
	offset := pdfOffsetRegex.FindStringSubmatch(s)
	if len(offset) > 0 {
		start, end = offsetHandler(offset)
	}
	t = ClippingMark
	ctype := pdfTypeRegex.FindStringSubmatch(s)
	if len(ctype) > 0 {
		if ctype[1] == "标注" {
			t = ClippingMark
		} else if ctype[1] == "笔记" {
			t = ClippingNote
		} else {
			t = ClippingBookmark
		}
	}
	return
}
