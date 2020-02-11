package mdcon

import (
	"math"
	"regexp"
	"strings"
)

type title struct {
	level int
	text  string
}

var spaces = make([]string, 6)

func init() {
	for i := 0; i < 6; i++ {
		spaces[i] = fourSpaceText(i)
	}
}

func fourSpaceText(count int) string {
	s := ""
	for i := 0; i < count; i++ {
		s += "    "
	}
	return s
}

func fourSpace(count int) string {
	if count < 0 {
		return ""
	}
	if count >= 6 {
		return fourSpaceText(count)
	}
	return spaces[count]
}

func BuildMdContents(text string) string {
	lines := strings.Split(text, "\n")
	titles := findTitle(lines)

	contents := buildContents(titles)

	if existContents(text) {
		return replaceContents(text, contents)
	} else {
		return text
	}
}

func existContents(text string) bool {
	r, _ := regexp.Compile("-----\r?\n(\\* 目录)\r?\n[\\s\\S]*?-----")
	if r.FindString(text) != "" {
		return true
	}
	return false
}

func replaceContents(text, contents string) string {
	r, _ := regexp.Compile("-----\r?\n(\\* 目录)\r?\n[\\s\\S]*?-----")
	return r.ReplaceAllString(text, contents)
}

func buildContents(titles []string) string {
	sb := strings.Builder{}
	next := title{level: math.MaxInt32}

	for _, t := range titles {
		title := parseTitle(t)
		if title.level < next.level {
			sb.WriteString(buildTitleLink(title))
		} else if title.level == next.level {
			sb.WriteString(buildTitleLink(title))
		} else {
			sb.WriteString(buildTitleLink(title))
		}
		next = title
	}
	return "-----\n* 目录\n" + sb.String() + "-----\n"
}

func buildTitleLink(t title) string {
	return fourSpace(t.level-1) + "- [" + t.text + "](#" + t.text + ")\n"
}

func parseTitle(t string) title {
	r, _ := regexp.Compile("(#+)\\s+?(.+)")
	sp := r.FindAllStringSubmatch(t, -1)
	well := sp[0][1]
	text := sp[0][2]
	return title{level: len(well), text: text}
}

func findTitle(lines []string) []string {
	titles := make([]string, 0, len(lines)/20)

	blackFlag := false

	for _, l := range lines {
		if strings.HasPrefix(l, "`") {
			if blackFlag {
				blackFlag = false
			} else {
				blackFlag = true
			}
		}
		if blackFlag {
			continue
		}
		if strings.HasPrefix(l, "#") {
			titles = append(titles, strings.TrimSpace(l))
		}
	}
	return titles
}
