package songhunter

import (
	"regexp"
	"strconv"
	"strings"
)

func SearchTrackPattern(text string) []string {
	// Convert newlines to HTML <br> tags
	text = strings.ReplaceAll(text, "\n", "<br />")

	// Split text into lines
	lines := strings.Split(text, "<br />")

	// Clean lines
	var cleanedLines []string
	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip lines without a "-" character
		if !strings.Contains(line, " - ") {
			continue
		}

		// Strip HTML tags
		line = stripTags(line)

		// Remove HTML special characters
		line = htmlSpecialCharsDecode(line)

		// Remove certain words
		line = strings.ReplaceAll(line, "track", "")
		line = strings.ReplaceAll(line, "Track", "")
		line = strings.ReplaceAll(line, "ID", "")

		// Remove everything before the "-" character
		if pos := strings.Index(line, "-&gt;"); pos != -1 {
			line = line[pos+5:]
		}

		// Remove ">"" characters
		line = strings.ReplaceAll(line, "&gt;", "")

		// Remove "?" and "!" characters
		line = strings.ReplaceAll(line, "?", "")
		line = strings.ReplaceAll(line, "!", "")

		// Remove "/w" and "w/" characters
		line = strings.ReplaceAll(line, "/w", "")
		line = strings.ReplaceAll(line, "w/", "")

		// Remove "." characters
		line = strings.ReplaceAll(line, ".", "")

		// Remove time frame
		timeFrame := getIn(line, "[", "]")
		if len(timeFrame) > 0 {
			timeFrameString := "[" + timeFrame[0] + "]"
			line = strings.ReplaceAll(line, timeFrameString, "")
		}
		line = stripTime(line)
		secondsRegex := regexp.MustCompile(`^\d+:\d+:\d+\s+`)
		line = secondsRegex.ReplaceAllString(line, "")
		sRegex := regexp.MustCompile(`^\d+:\d+:\d+\s+-\s+`)
		line = sRegex.ReplaceAllString(line, "")
		tRegex := regexp.MustCompile(`^\d+:\d+`)
		line = tRegex.ReplaceAllString(line, "")

		// Trim whitespace and "-" characters
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "-")
		line = strings.TrimSpace(line)

		// Skip lines that contain URLs
		if strings.Contains(line, "http") || strings.Contains(line, "www.") {
			continue
		}

		// Remove numbers in front of the line
		prefixNumberRegex := regexp.MustCompile(`^\d+\s+`)
		line = prefixNumberRegex.ReplaceAllString(line, "")

		// Remove emojis
		emojiRegex := regexp.MustCompile(`[\p{So}\p{Sk}]+`)
		line = emojiRegex.ReplaceAllString(line, "")

		// Skip empty lines
		if line == "" {
			continue
		}

		cleanedLines = append(cleanedLines, line)
	}

	return cleanedLines
}

func stripTags(html string) string {
	return tagRegex.ReplaceAllString(html, "")
}

func htmlSpecialCharsDecode(s string) string {
	return htmlEntityRegex.ReplaceAllStringFunc(s, func(m string) string {
		if len(m) > 3 && m[1] == '#' {
			if num, err := strconv.Atoi(m[2 : len(m)-1]); err == nil {
				return string(rune(num))
			}
		} else {
			if entity, ok := htmlEntities[m[1:len(m)-1]]; ok {
				return entity
			}
		}
		return m
	})
}

func stripTime(s string) string {
	s = timeRegex.ReplaceAllString(s, " ")
	return s
}

func getIn(s, start, end string) []string {
	var result []string
	startPos := strings.Index(s, start)
	if startPos == -1 {
		return result
	}
	endPos := strings.Index(s[startPos+len(start):], end)
	if endPos == -1 {
		return result
	}
	endPos += startPos + len(start)
	result = append(result, s[startPos+len(start):endPos])
	return result
}

var tagRegex = regexp.MustCompile("<[^>]*>")
var htmlEntityRegex = regexp.MustCompile(`&[a-zA-Z0-9#]+;`)
var htmlEntities = map[string]string{
	"&quot;":   "\"",
	"&amp;":    "&",
	"&apos;":   "'",
	"&lt;":     "<",
	"&gt;":     ">",
	"&nbsp;":   " ",
	"&iexcl;":  "¡",
	"&cent;":   "¢",
	"&pound;":  "£",
	"&curren;": "¤",
	"&yen;":    "¥",
	"&brvbar;": "¦",
	"&sect;":   "§",
	"&uml;":    "¨",
	"&copy;":   "©",
	"&ordf;":   "ª",
	"&laquo;":  "«",
	"&not;":    "¬",
	"&shy;":    "\u00AD",
	"&reg;":    "®",
	"&macr;":   "¯",
	"&deg;":    "°",
	"&plusmn;": "±",
}
var timeRegex = regexp.MustCompile(`\[[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}\]`)
