package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"porter/api"
	"porter/requester"
	"porter/wlog"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

const url = "https://quanmin.baidu.com/mvideo/api"

type Order struct {
	CharmPoints string `json:"charmPoints"`
	CreateTime  string `json:"createTime"`
	StatusTips  string `json:"statusTips"`
	AppendInfo  struct {
		MainTitle []struct {
			Title string `json:"title"`
		} `json:"mainTitle"`

		SlaveTitle []struct {
			Title string `json:"title"`
		} `json:"slaveTitle"`
	} `json:"appendInfo"`
}
type Record struct {
	OrderList struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			OrderList  []Order     `json:"orderList"`
			WithDrawDi interface{} `json:"withdrawDi"`
			IssueDi    interface{} `json:"issueDi"`
			HasMore    int         `json:"hasMore"`
		} `json:"data"`
	} `json:"orderlist"`
}

func main() {
	file, err := os.Open("./records.txt")
	if err != nil {
		wlog.Fatal("读取文件错误:", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		wlog.Fatal("读取文件内容失败:", err)
	}

	excelFile := xlsx.NewFile()

	// 先提取出账号密码
	lines := strings.Split(string(content), "\n")
	for _, bduss := range lines {
		if bduss == "" {
			continue
		}
		bduss = strings.TrimSpace(bduss)
		apiQuanminUser, err := api.GetQuanminInfo(bduss)
		if err != nil {
			wlog.Info("获取全民用户数据失败:", err)
			continue
		}
		orders, err := GetUserRecords(bduss)
		if err != nil {
			wlog.Info("错误:", err)
			continue
		}

		sheet, _ := excelFile.AddSheet(apiQuanminUser.Mine.Data.User.UserName)
		row := sheet.AddRow()

		cell := row.AddCell()
		cell.Value = "类型"

		cell = row.AddCell()
		cell.Value = "钻石数量"

		cell = row.AddCell()
		cell.Value = "获取时间"

		cell = row.AddCell()
		cell.Value = "状态"

		WriteToExcel(sheet, orders)
	}
	excelFile.Save("records.xlsx")
}

func WriteToExcel(sheet *xlsx.Sheet, orders []Order) {
	for _, order := range orders {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = order.AppendInfo.MainTitle[0].Title

		cell = row.AddCell()
		cell.Value = order.CharmPoints

		cell = row.AddCell()
		cell.Value = order.CreateTime

		cell = row.AddCell()
		cell.Value = order.StatusTips
	}
}
func GetUserRecords(bduss string) ([]Order, error) {
	orders := make([]Order, 0)
	hasmore := 1
	withdrawDi := ""
	issueDi := ""

	for hasmore == 1 {
		record, err := GetRecord(bduss, withdrawDi, issueDi)
		if err != nil {
			return orders, err
		}
		orders = append(orders, record.OrderList.Data.OrderList...)
		hasmore = record.OrderList.Data.HasMore
		withdrawDi = cover(record.OrderList.Data.WithDrawDi)

		issueDi = cover(record.OrderList.Data.IssueDi)
	}

	return orders, nil
}

func cover(arg interface{}) string {
	result := ""
	switch arg.(type) {
	case json.Number:
		result = arg.(json.Number).String()
		wlog.Info("json number:", arg)
	case int:
		result = strconv.Itoa(arg.(int))
		wlog.Info("json int:", arg)
	case string:
		result = arg.(string)
		wlog.Info("json string:", arg)
	}

	return result
}
func GetRecord(bduss string, withdrawDi string, issueDi string) (*Record, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", requester.UserAgent)
	cookie := http.Cookie{Name: "BDUSS", Value: bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	req.AddCookie(&cookie)

	q := req.URL.Query()
	q.Add("api_name", "orderlist")
	q.Add("tab", "all")
	if withdrawDi != "" {
		q.Add("withdrawDi", withdrawDi)
	}
	if issueDi != "" {
		q.Add("issueDi", issueDi)
	}

	req.URL.RawQuery = q.Encode()

	wlog.Info("raw query:", req.URL.RawQuery)

	resq, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		return nil, err
	}

	record := &Record{}
	if err := json.Unmarshal(data, record); err != nil {
		return nil, err
	}

	return record, nil
}
