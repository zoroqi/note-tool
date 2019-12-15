package contents

import "strings"

var img = map[string]bool{"jpg": true, "img": true, "jpeg": true, "bmp": true, "gif": true}

var levelSpace []string

const FOUR_SPACE = "    "

func init() {
	for i := 0; i < 20; i++ {
		levelSpace = append(levelSpace, strings.Repeat(FOUR_SPACE, i))
	}

}

func ImgRemove(s string) bool {
	fileType := FileType(s)
	if fileType == "" {
		return false
	}
	return img[fileType]
}


func HiddenRemove(s string) bool {
	return strings.HasPrefix(s, ".")
}

func ContentsRemove(contentsName string) Predicate {
	c := strings.TrimSpace(strings.ToLower(contentsName))
	return func(s string) bool {
		return c == s
	}
}

func PreSpace(level int) string {
	if level < len(levelSpace) {
		return levelSpace[level]
	}
	return strings.Repeat(FOUR_SPACE, level)
}

func FileType(s string) string {
	i := strings.LastIndex(s, ".")
	if i < 0 {
		return ""
	}
	fileType := strings.ToLower(s[i+1:])
	return fileType
}
