package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"porter/wlog"
	"strings"

	"github.com/tealeg/xlsx"
)

type Item struct {
	Baijia         string
	BaijiaPassword string
	Email          string
	Passwrod       string
	BAIDUID        string
	BDUSS          string
	// PANPSC          string
	// PANWEB          string
	// STOKEN_BFESS    string
	// PTOKEN          string
	// pplogid         string
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		wlog.Fatal("读取文件错误:", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		wlog.Fatal("读取文件内容失败:", err)
	}

	// 先提取出账号密码
	lines := strings.Split(string(content), "\n")
	items := make([]*Item, 0)

	for _, one := range lines {
		if one == "" {
			continue
		}
		i := &Item{}

		segments := strings.Split(one, "----")
		if len(segments) != 5 {
			fmt.Printf("[警告]数据有误,将会跳过这行:\n [%s] \n", one)
			continue
		}

		i.Baijia = segments[0]
		i.BaijiaPassword = segments[1]
		i.Email = segments[2]
		i.Passwrod = segments[3]

		orther := segments[4]

		bdussIndex := strings.Index(orther, "BDUSS=")
		if bdussIndex == -1 {
			fmt.Printf("这一行没有找到bduss,将会跳过:\n [%s] \n", one)
			continue
		}
		bdussIndex += len("BDUSS=")
		i.BDUSS = orther[bdussIndex : bdussIndex+192]

		items = append(items, i)
	}

	writeXLSX(items)
}

func writeXLSX(items []*Item) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("数据")

	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度

	cell := row.AddCell()
	cell.Value = "百家号账号"
	cell = row.AddCell()
	cell.Value = "百家号密码"

	cell = row.AddCell()
	cell.Value = "邮箱"
	cell = row.AddCell()
	cell.Value = "邮箱密码"

	cell = row.AddCell()
	cell.Value = "BDUSS"

	for _, item := range items {

		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		cell := row.AddCell()
		cell.Value = item.Baijia
		cell = row.AddCell()
		cell.Value = item.BaijiaPassword

		cell = row.AddCell()
		cell.Value = item.Email
		cell = row.AddCell()
		cell.Value = item.Passwrod

		cell = row.AddCell()
		cell.Value = item.BDUSS

	}

	err := file.Save("output.xlsx")
	if err != nil {
		panic(err)
	}
}
