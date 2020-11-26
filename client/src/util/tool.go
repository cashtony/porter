package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"porter/wlog"
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
