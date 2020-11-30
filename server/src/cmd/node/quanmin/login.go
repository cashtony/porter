package quanmin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/iikira/baidu-tools/randominfo"
)

func ClientSignature(post map[string]string) {
	if post == nil {
		post = map[string]string{}
	}

	// 已经签名, 则重新签名
	if _, ok := post["sign"]; ok {
		delete(post, "sign")
	}

	var (
		bduss        = post["BDUSS"]
		model        = randominfo.GetPhoneModel(bduss)
		phoneIMEIStr = strconv.FormatUint(randominfo.SumIMEI(model+"_"+bduss), 10)
		m            = md5.New()
	)

	// 预设
	post["_client_type"] = "2"
	post["_client_version"] = "7.0.0.0"
	post["_phone_imei"] = phoneIMEIStr
	post["from"] = "mini_ad_wandoujia"
	post["model"] = model
	m.Write([]byte(bduss + "_" + post["_client_version"] + "_" + post["_phone_imei"] + "_" + post["from"]))
	post["cuid"] = strings.ToUpper(hex.EncodeToString(m.Sum(nil))) + "|" + StringReverse(phoneIMEIStr)

	keys := make([]string, 0, len(post))
	for key := range post {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))

	m.Reset()
	for _, key := range keys {
		m.Write([]byte(key + "=" + post[key]))
	}
	m.Write([]byte("tiebaclient!!!"))

	post["sign"] = strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

func ClientRawQuerySignature(rawQuery string) (signedRawQuery string) {
	filterString := fmt.Sprintf("%stiebaclient!!!", strings.Replace(rawQuery, "&", "", -1))
	m := md5.New()
	m.Write([]byte(filterString))
	signedRawQuery = rawQuery + "&sign=" + strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	return
}

// StringReverse 反转字符串, 此操作不会修改原值
func StringReverse(s string) string {
	newBytes := make([]byte, len(s))
	copy(newBytes, s)
	b := BytesReverse(newBytes)
	return *(*string)(unsafe.Pointer(&b))
}

// BytesReverse 反转字节数组, 此操作会修改原值
func BytesReverse(b []byte) []byte {
	length := len(b)
	for i := 0; i < length/2; i++ {
		b[i], b[length-i-1] = b[length-i-1], b[i]
	}
	return b
}
