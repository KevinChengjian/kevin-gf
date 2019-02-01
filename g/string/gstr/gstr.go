// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gstr provides useful API for string handling.
// 
// 字符串操作.
package gstr

import (
    "bytes"
    "fmt"
    "gitee.com/johng/gf/g/util/grand"
    "math"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"
)

// 字符串替换(大小写敏感)
func Replace(origin, search, replace string, count...int) string {
    n := -1
    if len(count) > 0 {
        n = count[0]
    }
    return strings.Replace(origin, search, replace, n)
}

// 使用map进行字符串替换(大小写敏感)
func ReplaceByMap(origin string, replaces map[string]string) string {
    result := origin
    for k, v := range replaces {
        result = strings.Replace(result, k, v, -1)
    }
    return result
}

// 字符串转换为小写
func ToLower(s string) string {
    return strings.ToLower(s)
}

// 字符串转换为大写
func ToUpper(s string) string {
    return strings.ToUpper(s)
}

// 字符串首字母转换为大写
func UcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterLower(s[0]) {
        return string(s[0] - 32) + s[1 :]
    }
    return s
}

// 字符串首字母转换为小写
func LcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterUpper(s[0]) {
        return string(s[0] + 32) + s[1 :]
    }
    return s
}

// Uppercase the first character of each word in a string.
//
// 大写字符串中每个单词的第一个字符。
func UcWords(str string) string {
    return strings.Title(str)
}

// Traverse the array to find the string index position, if not exist, return-1.
//
// 遍历数组查找字符串索引位置，如果不存在则返回-1，使用完整遍历查找.
func SearchArray (a []string, s string) int {
    for i, v := range a {
        if s == v {
            return i
        }
    }
    return -1
}

// 判断字符串是否在数组中
func InArray (a []string, s string) bool {
    return SearchArray(a, s) != -1
}

// 判断给定字符是否小写
func IsLetterLower(b byte) bool {
    if b >= byte('a') && b <= byte('z') {
        return true
    }
    return false
}

// 判断给定字符是否大写
func IsLetterUpper(b byte) bool {
    if b >= byte('A') && b <= byte('Z') {
        return true
    }
    return false
}

// 判断锁给字符串是否为数字
func IsNumeric(s string) bool {
    length := len(s)
    if length == 0 {
        return false
    }
    for i := 0; i < len(s); i++ {
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}

// 字符串截取，支持中文
func SubStr(str string, start int, length...int) (substr string) {
    // 将字符串的转换成[]rune
    rs  := []rune(str)
    lth := len(rs)
    // 简单的越界判断
    if start < 0 {
        start = 0
    }
    if start >= lth {
        start = lth
    }
    end := lth
    if len(length) > 0 {
        end = start + length[0]
        if end < start {
            end = lth
        }
    }
    if end > lth {
        end = lth
    }
    // 返回子串
    return string(rs[start : end])
}

// 字符串长度截取限制，超过长度限制被截取并在字符串末尾追加指定的内容，支持中文
func StrLimit(str string, length int, suffix...string) (string) {
    rs := []rune(str)
    if len(str) < length {
        return str
    }
    addStr := "..."
    if len(suffix) > 0 {
        addStr = suffix[0]
    }
    return string(rs[0 : length]) + addStr
}

// Reverse a string.
func Reverse(str string) string {
    runes := []rune(str)
    for i, j := 0, len(runes) - 1; i < j; i, j = i + 1, j - 1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// Format a number with grouped thousands.
// decimals: Sets the number of decimal points.
// decPoint: Sets the separator for the decimal point.
// thousandsSep: Sets the thousands separator.
// See http://php.net/manual/en/function.number-format.php.
//
// 以千位分隔符方式格式化一个数字.
func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
    neg := false
    if number < 0 {
        number = -number
        neg = true
    }
    // Will round off
    str := fmt.Sprintf("%."+strconv.Itoa(decimals)+"F", number)
    prefix, suffix := "", ""
    if decimals > 0 {
        prefix = str[ : len(str) - (decimals + 1)]
        suffix = str[len(str) - decimals : ]
    } else {
        prefix = str
    }
    sep := []byte(thousandsSep)
    n, l1, l2 := 0, len(prefix), len(sep)
    // thousands sep num
    c := (l1 - 1) / 3
    tmp := make([]byte, l2*c+l1)
    pos := len(tmp) - 1
    for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
        if l2 > 0 && n > 0 && n%3 == 0 {
            for j := range sep {
                tmp[pos] = sep[l2-j-1]
                pos--
            }
        }
        tmp[pos] = prefix[i]
    }
    s := string(tmp)
    if decimals > 0 {
        s += decPoint + suffix
    }
    if neg {
        s = "-" + s
    }

    return s
}

// Split a string into smaller chunks.
// Can be used to split a string into smaller chunks which is useful for
// e.g. converting BASE64 string output to match RFC 2045 semantics.
// It inserts end every chunkLen characters.
//
// 将字符串分割成小块。使用此函数将字符串分割成小块非常有用。
// 例如将BASE64的输出转换成符合RFC2045语义的字符串。
// 它会在每 chunkLen 个字符后边插入 end。
func ChunkSplit(body string, chunkLen int, end string) string {
    if end == "" {
        end = "\r\n"
    }
    runes, endRunes := []rune(body), []rune(end)
    l := len(runes)
    if l <= 1 || l < chunkLen {
        return body + end
    }
    ns := make([]rune, 0, len(runes) + len(endRunes))
    for i := 0; i < l; i += chunkLen {
        if i + chunkLen > l {
            ns = append(ns, runes[i : ]...)
        } else {
            ns = append(ns, runes[i : i + chunkLen]...)
        }
        ns = append(ns, endRunes...)
    }
    return string(ns)
}

//  Return information about words used in a string.
func WordCount(str string) []string {
    return strings.Fields(str)
}

// Wraps a string to a given number of characters.
// TODO: Enable cut param, see http://php.net/manual/en/function.wordwrap.php.
func WordWrap(str string, width int, br string) string {
    if br == "" {
        br = "\n"
    }
    init := make([]byte, 0, len(str))
    buf := bytes.NewBuffer(init)
    var current int
    var wordBuf, spaceBuf bytes.Buffer
    for _, char := range str {
        if char == '\n' {
            if wordBuf.Len() == 0 {
                if current + spaceBuf.Len() > width {
                    current = 0
                } else {
                    current += spaceBuf.Len()
                    spaceBuf.WriteTo(buf)
                }
                spaceBuf.Reset()
            } else {
                current += spaceBuf.Len() + wordBuf.Len()
                spaceBuf.WriteTo(buf)
                spaceBuf.Reset()
                wordBuf.WriteTo(buf)
                wordBuf.Reset()
            }
            buf.WriteRune(char)
            current = 0
        } else if unicode.IsSpace(char) {
            if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
                current += spaceBuf.Len() + wordBuf.Len()
                spaceBuf.WriteTo(buf)
                spaceBuf.Reset()
                wordBuf.WriteTo(buf)
                wordBuf.Reset()
            }
            spaceBuf.WriteRune(char)
        } else {
            wordBuf.WriteRune(char)
            if current + spaceBuf.Len()+wordBuf.Len() > width && wordBuf.Len() < width {
                buf.WriteString(br)
                current = 0
                spaceBuf.Reset()
            }
        }
    }

    if wordBuf.Len() == 0 {
        if current + spaceBuf.Len() <= width {
            spaceBuf.WriteTo(buf)
        }
    } else {
        spaceBuf.WriteTo(buf)
        wordBuf.WriteTo(buf)
    }
    return buf.String()
}

// Get string length of unicode.
func RuneLen(str string) int {
    return utf8.RuneCountInString(str)
}

// Repeat a string.
func Repeat(input string, multiplier int) string {
    return strings.Repeat(input, multiplier)
}

// Returns part of haystack string starting from and including the first occurrence of needle to the end of haystack.
// See http://php.net/manual/en/function.strstr.php.
//
// 查找字符串的首次出现。返回 haystack 字符串从 needle 第一次出现的位置开始到 haystack 结尾的字符串。
func Str(haystack string, needle string) string {
    if needle == "" {
        return ""
    }
    idx := strings.Index(haystack, needle)
    if idx == -1 {
        return ""
    }
    return haystack[idx + len([]byte(needle)) - 1 : ]
}

// Randomly shuffles a string.
//
// 将字符串打乱。
func Shuffle(str string) string {
    runes := []rune(str)
    s     := make([]rune, len(runes))
    for i, v := range grand.Perm(len(runes)) {
        s[i] = runes[v]
    }
    return string(s)
}

// Strip whitespace (or other characters) from the beginning and end of a string.
//
// 去除字符串首尾处的空白字符（或者其他字符）。
func Trim(str string, characterMask ...string) string {
    if len(characterMask) > 0 {
        return strings.Trim(str, characterMask[0])
    } else {
        return strings.TrimSpace(str)
    }
}

// Strip whitespace (or other characters) from the beginning of a string.
//
// 去除字符串首的空白字符（或者其他字符）。
func TrimLeft(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimLeft(str, mask)
}

// Strip whitespace (or other characters) from the end of a string.
//
// 去除字符串尾的空白字符（或者其他字符）。
func TrimRight(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimRight(str, mask)
}

// Split a string by a string, to an array.
//
// 此函数返回由字符串组成的数组，每个元素都是 str 的一个子串，它们被字符串 delimiter 作为边界点分割出来。
func Split(str, delimiter string) []string {
    return strings.Split(str, delimiter)
}

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
//
// 用sep将字符串数组array连接为一个字符串。
func Join(array []string, sep string) string {
    return strings.Join(array, sep)
}

// Split a string by a string, to an array.
// See http://php.net/manual/en/function.explode.php.
//
// 此函数返回由字符串组成的数组，每个元素都是 str 的一个子串，它们被字符串 delimiter 作为边界点分割出来。
func Explode(delimiter, str string) []string {
    return Split(str, delimiter)
}

// Join array elements with a string.
// http://php.net/manual/en/function.implode.php
//
// 用glue将字符串数组pieces连接为一个字符串。
func Implode(glue string, pieces []string) string {
    var buf bytes.Buffer
    l := len(pieces)
    for _, str := range pieces {
        buf.WriteString(str)
        if l--; l > 0 {
            buf.WriteString(glue)
        }
    }
    return buf.String()
}

// Generate a single-byte string from a number.
//
// 返回相对应于 ascii 所指定的单个字符。
func Chr(ascii int) string {
    return string(ascii)
}

// Convert the first byte of a string to a value between 0 and 255.
//
// 解析 char 二进制值第一个字节为 0 到 255 范围的无符号整型类型。
func Ord(char string) int {
    return int(char[0])
}

// 按照百分比从字符串中间向两边隐藏字符(主要用于姓名、手机号、邮箱地址、身份证号等的隐藏)，支持utf-8中文，支持email格式。
func HideStr(str string, percent int, hide string) string {
    array := strings.Split(str, "@")
    if len(array) > 1 {
        str = array[0]
    }
    rs       := []rune(str)
    length   := len(rs)
    mid      := math.Floor(float64(length/2))
    hideLen  := int(math.Floor(float64(length) * (float64(percent)/100)))
    start    := int(mid - math.Floor(float64(hideLen) / 2))
    hideStr  := []rune("")
    hideRune := []rune(hide)
    for i := 0; i < int(hideLen); i++ {
        hideStr = append(hideStr, hideRune...)
    }
    buffer := bytes.NewBuffer(nil)
    buffer.WriteString(string(rs[0 : start]))
    buffer.WriteString(string(hideStr))
    buffer.WriteString(string(rs[start + hideLen : ]))
    if len(array) > 1 {
        buffer.WriteString(array[1])
    }
    return buffer.String()
}

// Inserts HTML line breaks before all newlines in a string.
// \n\r, \r\n, \r, \n
//
// 在字符串 string 所有新行之前插入 '<br />' 或 '<br>'，并返回。
func Nl2Br(str string, isXhtml...bool) string {
    r, n, runes := '\r', '\n', []rune(str)
    var br []byte
    if len(isXhtml) > 0 && isXhtml[0] {
        br = []byte("<br />")
    } else {
        br = []byte("<br>")
    }
    skip := false
    length := len(runes)
    var buf bytes.Buffer
    for i, v := range runes {
        if skip {
            skip = false
            continue
        }
        switch v {
            case n, r:
                if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
                    buf.Write(br)
                    skip = true
                    continue
                }
                buf.Write(br)
            default:
                buf.WriteRune(v)
        }
    }
    return buf.String()
}

// Quote string with slashes.
//
// 转义字符串中的单引号（'）、双引号（"）、反斜线（\）与 NUL（NULL 字符）。
func AddSlashes(str string) string {
    var buf bytes.Buffer
    for _, char := range str {
        switch char {
            case '\'', '"', '\\':
                buf.WriteRune('\\')
        }
        buf.WriteRune(char)
    }
    return buf.String()
}

// Un-quotes a quoted string.
//
// 反转义字符串。
func StripSlashes(str string) string {
    var buf bytes.Buffer
    l, skip := len(str), false
    for i, char := range str {
        if skip {
            skip = false
        } else if char == '\\' {
            if i + 1 < l && str[i + 1] == '\\' {
                skip = true
            }
            continue
        }
        buf.WriteRune(char)
    }
    return buf.String()
}

// Returns a version of str with a backslash character (\) before every character that is among:
// .\+*?[^]($)
//
// 转义字符串，转义的特殊字符包括：.\+*?[^]($)。
func QuoteMeta(str string) string {
    var buf bytes.Buffer
    for _, char := range str {
        switch char {
        case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
            buf.WriteRune('\\')
        }
        buf.WriteRune(char)
    }
    return buf.String()
}