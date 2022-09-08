/*
 */
package common

import (
	"bytes"
	"compress/zlib"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	// 可用字符数组的长度，为26个小写英文字母+10个数字
	charArrayLen = 36
)

func Atoi(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func IsValidString(source string) bool {
	return source != "" && len(strings.TrimSpace(source)) > 0
}

//随机字符串
func RandomString(length int) string {
	str := "0123456789"
	strLen := len(str)
	strBytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, strBytes[r.Intn(strLen)])
	}
	return string(result)
}

func BufferJoin(elems []string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}

	var buffer bytes.Buffer
	for _, s := range elems {
		buffer.WriteString(s)
	}
	return buffer.String()
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func UInt64ArrayToString(array []uint64) string {
	s := ""
	if len(array) == 0 {
		return ""
	}
	for _, i := range array {
		s += strconv.FormatUint(i, 10) + ","
	}

	return s[:len(s)-1]
}

func InArrayUint64(need uint64, needArr []uint64) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}