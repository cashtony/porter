package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"porter/wlog"
	"unicode"
)

func PrintBody(body io.Reader) {
	newBuff := make([]byte, 0)
	b := bytes.NewBuffer(newBuff)
	io.Copy(b, body)
	wlog.Info("body:", (string)(b.Bytes())) //在返回页面中显示内容。
}
func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

// 过滤掉特殊字符
func FilterSpecial(content string) string {
	var buffer bytes.Buffer
	for _, v := range content {
		if unicode.Is(unicode.Han, v) || unicode.IsLetter(v) || unicode.IsDigit(v) || unicode.IsPunct(v) {
			buffer.WriteRune(v)
			continue
		}
	}

	return buffer.String()
}
