package mdcon

import (
	"math"
	"net/url"
	"regexp"
	"strings"
)

type title struct {
	level int
	text  string
}

const contentRegex = "-----\r?\n(\\* 目录)\r?\n[\\s\\S]*?-----"

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

func SetContents(text string) string {
	lines := strings.Split(text, "\n")
	titles := findTitle(lines)
	if len(titles) == 0 {
		return text
	}

	contents := buildContents(titles)

	if existContents(text) {
		return replaceContents(text, contents)
	} else {
		return mergeContents(lines, contents)
	}
}

func existContents(text string) bool {
	r, _ := regexp.Compile(contentRegex)
	if r.FindString(text) != "" {
		return true
	}
	return false
}

func replaceContents(text, contents string) string {
	r, _ := regexp.Compile(contentRegex)
	return r.ReplaceAllString(text, contents)
}

func mergeContents(lines []string, contents string) string {
	sb := strings.Builder{}
	insert := true
	blockFlag := false
	for _, l := range lines {
		if strings.HasPrefix(l, "`") {
			if blockFlag {
				blockFlag = false
			} else {
				blockFlag = true
			}
		}
		if !blockFlag && insert {
			sb.WriteString(l)
			sb.WriteRune('\n')
			if strings.HasPrefix(l, "# ") {
				sb.WriteRune('\n')
				sb.WriteString(contents)
				sb.WriteRune('\n')
				insert = false
			}
		} else {
			sb.WriteString(l)
			sb.WriteRune('\n')
		}
	}
	if insert {
		return "\n" + contents + "\n" + sb.String()
	}
	return sb.String()
}

func buildContents(titles []string) string {
	sb := strings.Builder{}
	pre := title{level: math.MaxInt32}
	preSpaceCount := 1

	for _, t := range titles {
		title := parseTitle(t)
		if title.level < pre.level {
			preSpaceCount = preSpaceCount - 1
			sb.WriteString(buildTitleLink(title, preSpaceCount))
		} else if title.level == pre.level {
			sb.WriteString(buildTitleLink(title, preSpaceCount))
		} else {
			preSpaceCount = preSpaceCount + 1
			sb.WriteString(buildTitleLink(title, preSpaceCount))
		}
		pre = title
	}
	return "-----\n* 目录\n" + sb.String() + "-----\n"
}

func buildTitleLink(t title, spaceCount int) string {
	return fourSpace(spaceCount) + "- [" + t.text + "](#" + url.QueryEscape(t.text) + ")\n"
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

	blockFlag := false

	for _, l := range lines {
		if strings.HasPrefix(l, "`") {
			if blockFlag {
				blockFlag = false
			} else {
				blockFlag = true
			}
		}
		if blockFlag {
			continue
		}
		if strings.HasPrefix(l, "#") {
			titles = append(titles, strings.TrimSpace(l))
		}
	}
	return titles
}
