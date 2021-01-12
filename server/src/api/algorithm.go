package api

import (
	"crypto/md5"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func getXGon(params string, stub, cookie, sessionid string) string {
	NULL_MD5_STRING := "00000000000000000000000000000000"
	sb := ""
	if len(params) < 1 {
		sb = NULL_MD5_STRING
	} else {
		sb = encryption(params)
	}

	if len(stub) < 1 {
		sb += NULL_MD5_STRING
	} else {
		sb += stub
	}

	if len(cookie) < 1 {
		sb += NULL_MD5_STRING
	} else {
		sb += encryption(cookie)
	}

	if len(sessionid) < 1 {
		sb += NULL_MD5_STRING
	} else {
		sb += encryption(sessionid)
	}

	return sb
}

func xGorgon(timeMillis int, data []string) string {
	extraData := make([]string, 0)
	extraData = append(extraData, "3")
	extraData = append(extraData, "61")
	extraData = append(extraData, "41")
	extraData = append(extraData, "10")
	extraData = append(extraData, "80")
	extraData = append(extraData, "0")

	extraData2 := input(timeMillis, data)
	initialize(extraData2)
	handle(extraData2)
	for _, item := range extraData2 {
		extraData = append(extraData, item)
	}

	xGorgonStr := ""
	for _, item := range extraData {
		temp := item
		if len(temp) > 1 {
			xGorgonStr += temp
		} else {
			xGorgonStr += "0"
			xGorgonStr += temp
		}
	}

	return xGorgonStr
}

func initialize(data []string) {
	myhex := 0
	byteTable2 := strings.Split(byteTable, " ")
	for i := 0; i < len(data); i++ {
		hex1 := int64(0)
		if i == 0 {

			index, _ := strconv.ParseInt(byteTable2[0], 16, 64)
			hex1, _ = strconv.ParseInt(byteTable2[index-1], 16, 64)

			byteTable2[0] = strconv.FormatInt(hex1, 16)
		} else if i == 1 {
			arg1, _ := strconv.ParseInt("D6", 16, 64)
			arg2, _ := strconv.ParseInt("28", 16, 64)
			temp := arg1 + arg2
			if temp > 256 {
				temp -= 256
			}

			hex1, _ = strconv.ParseInt(byteTable2[temp-1], 16, 64)
			myhex = int(temp)
			byteTable2[i] = strconv.FormatInt(hex1, 16)
		} else {
			tempInt, _ := strconv.ParseInt(byteTable2[i], 16, 64)
			temp := myhex + int(tempInt)
			if temp > 256 {
				temp -= 256
			}
			hex1, _ = strconv.ParseInt(byteTable2[temp-1], 16, 64)
			myhex = temp
			byteTable2[i] = strconv.FormatInt(hex1, 16)
		}

		if hex1*2 > 256 {
			hex1 = hex1*2 - 256
		} else {
			hex1 = hex1 * 2
		}

		hex2 := byteTable2[int(hex1)-1]
		num1, _ := strconv.ParseInt(hex2, 16, 64)
		num2, _ := strconv.ParseInt(data[i], 16, 64)
		result := num1 ^ num2
		data[i] = strconv.FormatInt(result, 16)
	}

	for i := range data {
		data[i] = strings.ReplaceAll(data[i], "0x", "")
	}
}

func input(timeMillis int, input []string) []string {
	result := make([]string, 0)

	for i := 0; i < 4; i++ {
		num, _ := strconv.Atoi(input[i])
		temp := fmt.Sprintf("%x", num)
		if num < 0 {
			result = append(result, temp[6:])
		} else {
			result = append(result, temp)
		}
	}

	for i := 0; i < 4; i++ {
		result = append(result, "0")
	}

	for i := 0; i < 4; i++ {
		num, _ := strconv.Atoi(input[i+32])
		temp := fmt.Sprintf("%x", num)

		if num < 0 {
			result = append(result, temp[6:])
		} else {
			result = append(result, temp)
		}
	}

	for i := 0; i < 4; i++ {
		result = append(result, "0")
	}

	timeStr := fmt.Sprintf("%x", timeMillis)
	timeStr = strings.ReplaceAll(timeStr, "0x", "")
	for i := 0; i < 4; i++ {
		result = append(result, timeStr[i*2:i*2+2])
	}

	for i := 0; i < len(result); i++ {
		result[i] = strings.ReplaceAll(result[i], "0x", "")
	}

	return result
}

func encryption(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return strings.ToLower(md5str)
}

func strToByte(str string) []string {
	result := make([]string, 0)
	for i := 0; i < len(str); i += 2 {
		c1 := string(str[i])
		c2 := string(str[i+1])
		c := (str2Hex(c1) << 4) + str2Hex(c2)
		result = append(result, strconv.Itoa(c))
	}

	return result
}

func str2Hex(str string) int {
	cursor := 0
	upStr := strings.ToUpper(str)
	for _, c := range upStr {
		tmp := int(c)
		if tmp <= int('9') {
			cursor = cursor << 4
			cursor += tmp - int('0')
		} else if int('A') <= tmp && tmp <= int('F') {
			cursor = cursor << 4
			cursor += tmp - int('A') + 10
		}
	}

	return cursor
}

func handle(data []string) {
	for i, value := range data {
		byte1 := value
		if len(byte1) < 2 {
			byte1 += "0"
		} else {
			byte1 = string(data[i][1]) + string(data[i][0])
		}
		if i < len(data)-1 {
			num1, _ := strconv.ParseInt(byte1, 16, 64)
			num2, _ := strconv.ParseInt(data[i+1], 16, 64)
			byte1 = strconv.FormatInt(num1^num2, 16)
		} else {
			num1, _ := strconv.ParseInt(byte1, 16, 64)
			num2, _ := strconv.ParseInt(data[0], 16, 64)
			byte1 = strconv.FormatInt(num1^num2, 16)
		}

		clcu1, _ := strconv.ParseInt(byte1, 16, 64)
		clcu2, _ := strconv.ParseInt("AA", 16, 64)
		floata := float64((clcu1 & clcu2) / 2)
		a := int64(math.Abs(floata))

		clcu3, _ := strconv.ParseInt("55", 16, 64)
		byte2 := ((clcu1 & clcu3) * 2) | a

		clcu4, _ := strconv.ParseInt("33", 16, 64)
		clcu5, _ := strconv.ParseInt("cc", 16, 64)
		byte2 = ((byte2 & clcu4) * 4) | ((byte2 & clcu5) / 4)
		byte3 := strconv.FormatInt(byte2, 16)
		if len(byte3) > 1 {
			byte3 = string(byte3[1]) + string(byte3[0])
		} else {
			byte3 += "0"
		}

		clcu6, _ := strconv.ParseInt("FF", 16, 64)
		clcu7, _ := strconv.ParseInt("14", 16, 64)
		intByte3, _ := strconv.ParseInt(byte3, 16, 64)
		byte4 := (intByte3 ^ clcu6) ^ clcu7
		data[i] = strconv.FormatInt(byte4, 16)
	}
}
