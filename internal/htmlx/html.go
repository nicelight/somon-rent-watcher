package htmlx

import (
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Node is a deliberately small DOM representation sufficient for scraping.
type Node struct {
	Tag      string
	TextData string
	Attr     map[string]string
	Parent   *Node
	Children []*Node
}

var removableBlocks = []*regexp.Regexp{
	regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script\s*>`),
	regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style\s*>`),
	regexp.MustCompile(`(?is)<noscript\b[^>]*>.*?</noscript\s*>`),
	regexp.MustCompile(`(?is)<template\b[^>]*>.*?</template\s*>`),
	regexp.MustCompile(`(?is)<!--.*?-->`),
}

var voidTags = map[string]struct{}{
	"area": {}, "base": {}, "br": {}, "col": {}, "embed": {}, "hr": {},
	"img": {}, "input": {}, "link": {}, "meta": {}, "param": {},
	"source": {}, "track": {}, "wbr": {},
}

var blockTags = map[string]struct{}{
	"address": {}, "article": {}, "aside": {}, "blockquote": {}, "body": {},
	"br": {}, "button": {}, "dd": {}, "details": {}, "div": {}, "dl": {},
	"dt": {}, "fieldset": {}, "figcaption": {}, "figure": {}, "footer": {},
	"form": {}, "h1": {}, "h2": {}, "h3": {}, "h4": {}, "h5": {}, "h6": {},
	"header": {}, "hr": {}, "li": {}, "main": {}, "nav": {}, "ol": {},
	"p": {}, "pre": {}, "section": {}, "summary": {}, "table": {}, "tbody": {},
	"td": {}, "tfoot": {}, "th": {}, "thead": {}, "tr": {}, "ul": {},
}

// Parse builds a deliberately small, tolerant HTML tree. It is not a full
// HTML5 parser, but it accepts the forms relevant to scraping: quoted and
// unquoted attributes, boolean attributes, void elements, optional closing
// tags, comments, and partially malformed markup. Unlike encoding/xml it does
// not stop at the first non-XML construct.
func Parse(data []byte) (*Node, error) {
	clean := append([]byte(nil), data...)
	for _, re := range removableBlocks {
		clean = re.ReplaceAll(clean, nil)
	}

	root := &Node{Tag: "#document", Attr: map[string]string{}}
	stack := []*Node{root}
	for i := 0; i < len(clean); {
		rel := bytesIndexByte(clean[i:], '<')
		if rel < 0 {
			appendTextNode(stack[len(stack)-1], clean[i:])
			break
		}
		lt := i + rel
		if lt > i {
			appendTextNode(stack[len(stack)-1], clean[i:lt])
		}

		if hasPrefixFold(clean[lt:], "<!--") {
			if end := bytesIndex(clean[lt+4:], []byte("-->")); end >= 0 {
				i = lt + 4 + end + 3
				continue
			}
			break
		}
		if lt+1 >= len(clean) {
			appendTextNode(stack[len(stack)-1], clean[lt:])
			break
		}
		if clean[lt+1] == '!' || clean[lt+1] == '?' {
			end := findTagEnd(clean, lt+2)
			if end < 0 {
				break
			}
			i = end + 1
			continue
		}

		end := findTagEnd(clean, lt+1)
		if end < 0 {
			appendTextNode(stack[len(stack)-1], clean[lt:])
			break
		}
		inside := strings.TrimSpace(string(clean[lt+1 : end]))
		if inside == "" {
			i = end + 1
			continue
		}
		if strings.HasPrefix(inside, "/") {
			tag := readName(strings.TrimSpace(strings.TrimPrefix(inside, "/")))
			if tag != "" {
				for j := len(stack) - 1; j > 0; j-- {
					if stack[j].Tag == tag {
						stack = stack[:j]
						break
					}
				}
			}
			i = end + 1
			continue
		}

		selfClosing := strings.HasSuffix(strings.TrimSpace(inside), "/")
		if selfClosing {
			inside = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(inside), "/"))
		}
		tag, attrs := parseStartTag(inside)
		if tag == "" {
			i = end + 1
			continue
		}
		autoCloseForStart(&stack, tag)
		parent := stack[len(stack)-1]
		n := &Node{Tag: tag, Attr: attrs, Parent: parent}
		parent.Children = append(parent.Children, n)
		if _, isVoid := voidTags[tag]; !isVoid && !selfClosing {
			stack = append(stack, n)
		}
		i = end + 1
	}
	return root, nil
}

func appendTextNode(parent *Node, raw []byte) {
	if parent == nil || len(raw) == 0 {
		return
	}
	text := html.UnescapeString(string(raw))
	if text == "" {
		return
	}
	n := &Node{Tag: "#text", TextData: text, Attr: map[string]string{}, Parent: parent}
	parent.Children = append(parent.Children, n)
}

func findTagEnd(data []byte, start int) int {
	var quote byte
	for i := start; i < len(data); i++ {
		c := data[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			continue
		}
		if c == '>' {
			return i
		}
	}
	return -1
}

func parseStartTag(input string) (string, map[string]string) {
	attrs := make(map[string]string)
	i := 0
	skipSpaces := func() {
		for i < len(input) && isASCIISpace(input[i]) {
			i++
		}
	}
	skipSpaces()
	start := i
	for i < len(input) && isNameByte(input[i]) {
		i++
	}
	if start == i {
		return "", attrs
	}
	tag := strings.ToLower(input[start:i])
	for i < len(input) {
		skipSpaces()
		if i >= len(input) {
			break
		}
		nameStart := i
		for i < len(input) && isAttrNameByte(input[i]) {
			i++
		}
		if nameStart == i {
			i++
			continue
		}
		name := strings.ToLower(input[nameStart:i])
		skipSpaces()
		value := ""
		if i < len(input) && input[i] == '=' {
			i++
			skipSpaces()
			if i < len(input) && (input[i] == '\'' || input[i] == '"') {
				quote := input[i]
				i++
				valueStart := i
				for i < len(input) && input[i] != quote {
					i++
				}
				value = input[valueStart:i]
				if i < len(input) {
					i++
				}
			} else {
				valueStart := i
				for i < len(input) && !isASCIISpace(input[i]) {
					i++
				}
				value = strings.TrimSuffix(input[valueStart:i], "/")
			}
		}
		attrs[name] = html.UnescapeString(value)
	}
	return tag, attrs
}

func autoCloseForStart(stack *[]*Node, tag string) {
	if len(*stack) <= 1 {
		return
	}
	top := (*stack)[len(*stack)-1].Tag
	closeTop := false
	switch tag {
	case "li":
		closeTop = top == "li"
	case "dt", "dd":
		closeTop = top == "dt" || top == "dd"
	case "tr":
		closeTop = top == "tr"
	case "td", "th":
		closeTop = top == "td" || top == "th"
	case "option":
		closeTop = top == "option"
	case "p":
		closeTop = top == "p"
	default:
		if _, block := blockTags[tag]; block && top == "p" {
			closeTop = true
		}
	}
	if closeTop {
		*stack = (*stack)[:len(*stack)-1]
	}
}

func readName(input string) string {
	i := 0
	for i < len(input) && isNameByte(input[i]) {
		i++
	}
	return strings.ToLower(input[:i])
}

func isASCIISpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\f'
}

func isNameByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == ':' || c == '_'
}

func isAttrNameByte(c byte) bool {
	return !isASCIISpace(c) && c != '=' && c != '>' && c != '/' && c != '\'' && c != '"'
}

func bytesIndexByte(data []byte, value byte) int {
	for i, c := range data {
		if c == value {
			return i
		}
	}
	return -1
}

func bytesIndex(data, needle []byte) int {
	if len(needle) == 0 {
		return 0
	}
	for i := 0; i+len(needle) <= len(data); i++ {
		match := true
		for j := range needle {
			if data[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func hasPrefixFold(data []byte, prefix string) bool {
	if len(data) < len(prefix) {
		return false
	}
	return strings.EqualFold(string(data[:len(prefix)]), prefix)
}

func (n *Node) GetAttr(name string) string {
	if n == nil {
		return ""
	}
	return n.Attr[strings.ToLower(name)]
}

func (n *Node) HasClassPart(part string) bool {
	part = strings.ToLower(strings.TrimSpace(part))
	if part == "" {
		return false
	}
	for _, c := range strings.Fields(strings.ToLower(n.GetAttr("class"))) {
		if strings.Contains(c, part) {
			return true
		}
	}
	return false
}

func Walk(n *Node, fn func(*Node) bool) {
	if n == nil {
		return
	}
	if !fn(n) {
		return
	}
	for _, child := range n.Children {
		Walk(child, fn)
	}
}

func FindAll(n *Node, pred func(*Node) bool) []*Node {
	var out []*Node
	Walk(n, func(cur *Node) bool {
		if pred(cur) {
			out = append(out, cur)
		}
		return true
	})
	return out
}

func FindFirst(n *Node, pred func(*Node) bool) *Node {
	var found *Node
	Walk(n, func(cur *Node) bool {
		if found != nil {
			return false
		}
		if pred(cur) {
			found = cur
			return false
		}
		return true
	})
	return found
}

func Text(n *Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	appendText(&b, n, false)
	return NormalizeSpace(b.String())
}

func TextLines(n *Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	appendText(&b, n, true)
	lines := strings.Split(strings.ReplaceAll(b.String(), "\r", ""), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = NormalizeSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func appendText(b *strings.Builder, n *Node, lines bool) {
	if n.Tag == "#text" {
		b.WriteString(n.TextData)
		b.WriteByte(' ')
		return
	}
	_, block := blockTags[n.Tag]
	if lines && block {
		b.WriteByte('\n')
	}
	for _, child := range n.Children {
		appendText(b, child, lines)
	}
	if lines && block {
		b.WriteByte('\n')
	}
}

func NormalizeSpace(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '\u00a0' || r == '\u202f' {
			return ' '
		}
		return r
	}, collapseWhitespace(s)))
}

func collapseWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	space := false
	for _, r := range s {
		if unicode.IsSpace(r) || r == '\u00a0' || r == '\u202f' {
			if !space {
				b.WriteByte(' ')
				space = true
			}
			continue
		}
		space = false
		b.WriteRune(r)
	}
	return b.String()
}

func DirectText(n *Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	for _, child := range n.Children {
		if child.Tag == "#text" {
			b.WriteString(child.TextData)
			b.WriteByte(' ')
		}
	}
	return NormalizeSpace(b.String())
}

func ResolveURL(baseURL, ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if strings.HasPrefix(ref, "//") {
		return "https:" + ref
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ref
	}
	u, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return base.ResolveReference(u).String()
}

func FirstMetaContent(root *Node, keys map[string]string) string {
	if root == nil {
		return ""
	}
	var value string
	Walk(root, func(n *Node) bool {
		if value != "" {
			return false
		}
		if n.Tag != "meta" {
			return true
		}
		for attr, expected := range keys {
			if strings.EqualFold(n.GetAttr(attr), expected) {
				value = n.GetAttr("content")
				return false
			}
		}
		return true
	})
	return strings.TrimSpace(value)
}

func UniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func SortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
