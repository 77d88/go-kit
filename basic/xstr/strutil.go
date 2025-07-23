package xstr

import (
	"bytes"
	"errors"
	"github.com/77d88/go-kit/basic/xerror"
	"regexp"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// CamelCase 驼峰命名法转换字符串为camelCase形式。非字母和数字字符将被忽略。
func CamelCase(s string) string {
	var builder strings.Builder

	strs := splitIntoStrings(s, false)
	for i, str := range strs {
		if i == 0 {
			builder.WriteString(strings.ToLower(str))
		} else {
			builder.WriteString(Capitalize(str))
		}
	}

	return builder.String()
}

// Capitalize 首字母大写转换，将字符串的首字符转换为大写，其余转为小写。
func Capitalize(s string) string {
	result := make([]rune, len(s))
	for i, v := range s {
		if i == 0 {
			result[i] = unicode.ToUpper(v)
		} else {
			result[i] = unicode.ToLower(v)
		}
	}

	return string(result)
}

// UpperFirst 首字母大写转换，仅将字符串的首字符转换为大写。
func UpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	r = unicode.ToUpper(r)

	return string(r) + s[size:]
}

// LowerFirst 首字母小写转换，仅将字符串的首字符转换为小写。
func LowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	r = unicode.ToLower(r)

	return string(r) + s[size:]
}

// Pad 填充字符串，如果字符串长度小于指定大小，则在左侧、右侧或两侧填充指定字符
func Pad(source string, size int, padStr string) string {
	return padAtPosition(source, size, padStr, 0)
}

// PadStart 左侧填充字符串，如果字符串长度小于指定大小，则在左侧填充指定字符。
func PadStart(source string, size int, padStr string) string {
	return padAtPosition(source, size, padStr, 1)
}

// PadEnd 右侧填充字符串，如果字符串长度小于指定大小，则在右侧填充指定字符。
func PadEnd(source string, size int, padStr string) string {
	return padAtPosition(source, size, padStr, 2)
}

// KebabCase 转换字符串为kebab-case形式，非字母和数字字符将被忽略。
func KebabCase(s string) string {
	result := splitIntoStrings(s, false)
	return strings.Join(result, "-")
}

// UpperKebabCase 转换字符串为大写的KEBAB-CASE形式，非字母和数字字符将被忽略。
func UpperKebabCase(s string) string {
	result := splitIntoStrings(s, true)
	return strings.Join(result, "-")
}

// SnakeCase 转换字符串为snake_case形式，非字母和数字字符将被忽略。
func SnakeCase(s string) string {
	result := splitIntoStrings(s, false)
	return strings.Join(result, "_")
}

// UpperSnakeCase 转换字符串为大写的SNAKE_CASE形式，非字母和数字字符将被忽略。
func UpperSnakeCase(s string) string {
	result := splitIntoStrings(s, true)
	return strings.Join(result, "_")
}

// Before 返回源字符串中直到指定字符第一次出现的子串。
func Before(s, char string) string {
	i := strings.Index(s, char)

	if s == "" || char == "" || i == -1 {
		return s
	}

	return s[0:i]
}

// BeforeLast 返回源字符串中直到指定字符最后一次出现的子串。
func BeforeLast(s, char string) string {
	i := strings.LastIndex(s, char)

	if s == "" || char == "" || i == -1 {
		return s
	}

	return s[0:i]
}

// After 返回源字符串中指定字符之后的子串。
func After(s, char string) string {
	i := strings.Index(s, char)

	if s == "" || char == "" || i == -1 {
		return s
	}

	return s[i+len(char):]
}

// AfterLast 返回源字符串中指定字符最后一次出现之后的子串。
func AfterLast(s, char string) string {
	i := strings.LastIndex(s, char)

	if s == "" || char == "" || i == -1 {
		return s
	}

	return s[i+len(char):]
}

// IsString 检查给定值是否为字符串类型。
func IsString(v any) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

// Reverse 反转字符串中的字符顺序。
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Wrap 使用指定字符串包裹原字符串。
func Wrap(str string, wrapWith string) string {
	if str == "" || wrapWith == "" {
		return str
	}
	var sb strings.Builder
	sb.WriteString(wrapWith)
	sb.WriteString(str)
	sb.WriteString(wrapWith)

	return sb.String()
}

// Unwrap 移除原字符串中的包裹字符串。
func Unwrap(str string, wrapToken string) string {
	if str == "" || wrapToken == "" {
		return str
	}

	firstIndex := strings.Index(str, wrapToken)
	lastIndex := strings.LastIndex(str, wrapToken)

	if firstIndex == 0 && lastIndex > 0 && lastIndex <= len(str)-1 {
		if len(wrapToken) <= lastIndex {
			str = str[len(wrapToken):lastIndex]
		}
	}

	return str
}

// SplitEx 分割字符串，并可控制结果数组是否包含空字符串。
func SplitEx(s, sep string, removeEmptyString bool) []string {
	if sep == "" {
		return []string{}
	}

	n := strings.Count(s, sep) + 1
	a := make([]string, n)
	n--
	i := 0
	sepSave := 0
	ignore := false

	for i < n {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		ignore = false
		if removeEmptyString {
			if s[:m+sepSave] == "" {
				ignore = true
			}
		}
		if !ignore {
			a[i] = s[:m+sepSave]
			s = s[m+len(sep):]
			i++
		} else {
			s = s[m+len(sep):]
		}
	}

	var ret []string
	if removeEmptyString {
		if s != "" {
			a[i] = s
			ret = a[:i+1]
		} else {
			ret = a[:i]
		}
	} else {
		a[i] = s
		ret = a[:i+1]
	}

	return ret
}

// Substring 提取字符串中从指定偏移量开始的指定长度的子串
func Substring(s string, offset int, length uint) string {
	rs := []rune(s)
	size := len(rs)

	if offset < 0 {
		offset = size + offset
		if offset < 0 {
			offset = 0
		}
	}
	if offset > size {
		return ""
	}

	if length > uint(size)-uint(offset) {
		length = uint(size - offset)
	}

	str := string(rs[offset : offset+int(length)])

	return strings.Replace(str, "\x00", "", -1)
}

// SplitWords 将字符串按单词分割，单词仅包含字母字符。
func SplitWords(s string) []string {
	var word string
	var words []string
	var r rune
	var size, pos int

	isWord := false

	for len(s) > 0 {
		r, size = utf8.DecodeRuneInString(s)

		switch {
		case isLetter(r):
			if !isWord {
				isWord = true
				word = s
				pos = 0
			}

		case isWord && (r == '\'' || r == '-'):
			// is word

		default:
			if isWord {
				isWord = false
				words = append(words, word[:pos])
			}
		}

		pos += size
		s = s[size:]
	}

	if isWord {
		words = append(words, word[:pos])
	}

	return words
}

// WordCount 计算字符串中有效单词的数量，单词仅包含字母字符。
func WordCount(s string) int {
	var r rune
	var size, count int

	isWord := false

	for len(s) > 0 {
		r, size = utf8.DecodeRuneInString(s)

		switch {
		case isLetter(r):
			if !isWord {
				isWord = true
				count++
			}

		case isWord && (r == '\'' || r == '-'):
			// is word

		default:
			isWord = false
		}

		s = s[size:]
	}

	return count
}

// RemoveNonPrintable 移除字符串中的不可打印字符。
func RemoveNonPrintable(str string) string {
	result := strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, str)

	return result
}

// StringToBytes 无内存分配地将字符串转换为字节切片。
func StringToBytes(str string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&str))
}

// BytesToString 无内存分配地将字节切片转换为字符串。
func BytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// IsBlank 检查字符串是否为空白，即是否为空或仅包含空白字符。
func IsBlank(str string) bool {
	if len(str) == 0 {
		return true
	}
	// memory copies will occur here, but UTF8 will be compatible
	runes := []rune(str)
	for _, r := range runes {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// IsNotBlank 检查字符串是否不为空白，即是否非空且不只包含空白字符。
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// IndexOffset 在字符串中查找子串substr的位置，从索引`idxFrom`处开始搜索。
func IndexOffset(str string, substr string, idxFrom int) int {
	if idxFrom > len(str)-1 || idxFrom < 0 {
		return -1
	}

	return strings.Index(str[idxFrom:], substr) + idxFrom
}

// ReplaceWithMap 使用映射替换字符串中的特定子串。
func ReplaceWithMap(str string, replaces map[string]string) string {
	for k, v := range replaces {
		str = strings.ReplaceAll(str, k, v)
	}

	return str
}

// Split 分割字符串。
func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

// SplitAndTrim 分割字符串并去除每部分的前导和尾随空白字符，然后过滤掉空字符串。
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	result := make([]string, 0)

	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

// DefaultTrimChars 默认的去除字符集合，用于Trim系列函数。
var DefaultTrimChars = string([]byte{
	'\t', // 制表符
	'\v', // 垂直制表符
	'\n', // 换行符
	'\r', // 回车符
	'\f', // 换页符
	' ',  // 空格
	0x00, // 空字符
	0x85, // 删除字符
	0xA0, // 不间断空格
})

// Trim 去除字符串两端的空白字符，可自定义额外去除的字符。
func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars

	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}

	return strings.Trim(str, trimChars)
}

// HideString 隐藏字符串中的部分内容，用指定字符替换指定范围内的字符。
func HideString(origin string, start, end int, replaceChar string) string {
	size := len(origin)

	if start > size-1 || start < 0 || end < 0 || start > end {
		return origin
	}

	if end > size {
		end = size
	}

	if replaceChar == "" {
		return origin
	}

	startStr := origin[0:start]
	endStr := origin[end:size]

	replaceSize := end - start
	replaceStr := strings.Repeat(replaceChar, replaceSize)

	return startStr + replaceStr + endStr
}

// ContainsAll 检查目标字符串是否包含所有指定的子串。
func ContainsAll(str string, substrs []string) bool {
	for _, v := range substrs {
		if !strings.Contains(str, v) {
			return false
		}
	}

	return true
}

// ContainsAny 包含任意检查目标字符串是否包含子串列表中的任意子串。
func ContainsAny(str string, substrs ...string) bool {
	for _, v := range substrs {
		if strings.Contains(str, v) {
			return true
		}
	}

	return false
}

var (
	whitespaceRegexMatcher     *regexp.Regexp = regexp.MustCompile(`\s`)
	mutiWhitespaceRegexMatcher *regexp.Regexp = regexp.MustCompile(`[[:space:]]{2,}|[\s\p{Zs}]{2,}`)
)

// RemoveWhiteSpace // 移除空白字符从字符串中移除空白字符。replaceAll为真时移除所有空白字符，否则替换连续空白字符为单个空格。
func RemoveWhiteSpace(str string, repalceAll bool) string {
	if repalceAll && str != "" {
		return strings.Join(strings.Fields(str), "")
	} else if str != "" {
		str = mutiWhitespaceRegexMatcher.ReplaceAllString(str, " ")
		str = whitespaceRegexMatcher.ReplaceAllString(str, " ")
	}

	return strings.TrimSpace(str)
}

// SubInBetween 子串提取返回源字符串中start和end位置（不包括end）之间的子串。
func SubInBetween(str string, start string, end string) string {
	if _, after, ok := strings.Cut(str, start); ok {
		if before, _, ok := strings.Cut(after, end); ok {
			return before
		}
	}

	return ""
}

// HammingDistance Hamming距离计算两个字符串的Hamming距离，即对应位置上不同符号的数目。输入字符串长度不等时返回错误。
func HammingDistance(a, b string) (int, error) {
	if len(a) != len(b) {
		return -1, errors.New("a length and b length are unequal")
	}

	ar := []rune(a)
	br := []rune(b)

	var distance int
	for i, codepoint := range ar {
		if codepoint != br[i] {
			distance++
		}
	}

	return distance, nil
}

// Concat 连接字符串使用strings.Builder高效拼接输入的字符串，可选预计拼接长度。
func Concat(length int, str ...string) string {
	if len(str) == 0 {
		return ""
	}

	sb := strings.Builder{}
	if length <= 0 {
		sb.Grow(len(str[0]) * len(str))
	} else {
		sb.Grow(length)
	}

	for _, s := range str {
		sb.WriteString(s)
	}
	return sb.String()
}

// SplitTo 分割字符串并转为指定类型
// str: 待分割字符串
// delimiter: 分隔符
// parse: 类型转换函数
// ignoreErrors: true时忽略转换错误，false时遇到错误立即返回
func SplitTo[T any](str, delimiter string, parse func(v string) (T, error)) ([]T, error) {
	if str == "" {
		return nil, nil
	}

	ss := strings.Split(str, delimiter)
	result := make([]T, 0, len(ss))
	var err error
	for i, v := range ss {
		to, err := parse(v)
		if err != nil {
			err = errors.Join(err, xerror.Newf("parse[%d] %s error: %v", i, v, err))
			continue
		}
		result = append(result, to)
	}
	return result, err
}

// Repeat 重复字符串。
func Repeat(str string, count int) string {
	return strings.Repeat(str, count)
}

// EndsWith 判断字符串是否以给定的后缀结尾。
func EndsWith(str string, suffixes ...string) bool {
	if len(str) == 0 || len(suffixes) == 0 {
		return false
	}
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

// StartsWith 判断字符串是否以给定的前缀开头。
func StartsWith(str string, prefixes ...string) bool {
	if len(str) == 0 || len(prefixes) == 0 {
		return false
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

// TelNumberMask 手机号码掩码 最后4位数字不能打掩码
func TelNumberMask(value string, maskChar string) string {
	if len(value) <= 4 {
		return value
	}
	qh := ""
	tel := value
	// 如果有区号 则分别处理
	if ContainsAny(value, "-") {
		trim := SplitAndTrim(value, "-")
		qh = trim[0]
		tel = trim[1]
	}
	if len(qh) > 0 {
		// 区号只保留前后各一位
		if len(qh) >= 3 {
			qh = qh[:1] + Repeat(maskChar, len(qh)-2) + qh[len(qh)-1:]
		}
	}
	if len(tel) >= 11 {
		tel = tel[:3] + Repeat(maskChar, len(tel)-7) + tel[7:] // 保留前3位和后4位
	} else if len(tel) > 4 {
		// 保留前1后4
		tel = tel[:1] + Repeat(maskChar, len(tel)-5) + tel[len(tel)-4:]
	}
	if len(qh) > 0 {
		return qh + "-" + tel
	}
	return tel
}

// ParseWithStruct 使用struct数据解析模板
func ParseWithStruct(tplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("tpl").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParseWithMap 使用map数据解析模板
func ParseWithMap(tplStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("tpl").Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
