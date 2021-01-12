package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type APIWebDouyinUser struct {
	UID            string `json:"uid" gorm:"primaryKey"`
	UniqueUID      string `json:"unique_id" gorm:"primaryKey"` // 抖音号
	Nickname       string `json:"nickname"`
	AwemeCount     int    `json:"aweme_count"`
	FollowerCount  int    `json:"follower_count"`
	Signature      string `json:"signature" gorm:"-"` // 个人签名
	TotalFavorited string `json:"total_favorited"`
	AvatarMedium   struct {
		URLList []string `json:"url_list"`
	} `json:"avatar_medium"`

	SecUID string `json:"-"`
}

func NewAPIWebDouyinUser(shareURL string) (*APIWebDouyinUser, error) {
	if shareURL == "" {
		return nil, errors.New("shareURL不能为空")
	}
	secUID := GetSecID(shareURL)
	// 基本数据
	if secUID == "" {
		return nil, errors.New("secUID为空,不能获取抖音用户数据")
	}

	data := struct {
		StatusCode int               `json:"status_code"`
		UserInfo   *APIWebDouyinUser `json:"user_info"`
	}{}

	infoReq := fmt.Sprintf("%s?sec_uid=%s", define.GetUserInfo, secUID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, infoReq, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}
	req.Header.Add("User-Agent", requester.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取用户数据时失败:%s %s", infoReq, err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取response内容失败: %s", err)
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("解析用户json数据失败: %s", err)
	}

	if data.StatusCode != 0 {
		return nil, fmt.Errorf("返回的状态代码未成功: %d", data.StatusCode)
	}

	data.UserInfo.SecUID = secUID

	return data.UserInfo, nil
}

func GetSecID(shareURL string) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, shareURL, nil)
	if err != nil {
		wlog.Error("获取secid时创建请求失败:", err)
		return ""
	}
	req.Header.Add("User-Agent", requester.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("获取secid失败: %s %s", shareURL, err)
		return ""
	}
	defer resp.Body.Close()

	return resp.Request.URL.Query().Get("sec_uid")
}

type APIDouyinVideo struct {
	AwemeID    string `json:"aweme_id"`
	Desc       string `json:"desc"`
	CreateTime int64  `json:"create_time"`
	Video      struct {
		Cover struct {
			URLList []string `json:"url_list"`
		} `json:"cover"`
		Duration int `json:"duration"`
		Height   int `json:"height"`
		Width    int `json:"width"`
		PlayAddr struct {
			URI     string   `json:"uri"`
			URLList []string `json:"url_list"`
		} `json:"play_addr"`
	} `json:"video"`
}

func GetDouyinVideo(secUID, secSig string, cursor int64) ([]*APIDouyinVideo, int64, error) {
	if secUID == "" {
		return nil, 0, errors.New("secUID不能为空")
	}

	data := struct {
		AwemeList  []*APIDouyinVideo `json:"aweme_list"`
		HasMore    interface{}       `json:"has_more"`
		MaxCursor  int64             `json:"max_cursor"`
		MinCursor  int64             `json:"min_cursor"`
		StatusCode int               `json:"status_code"`
	}{}

	tryTimes := 0
	url := fmt.Sprintf("%s?sec_uid=%s&count=21&max_cursor=%d&aid=1128&_signature=%s&dytk=", define.GetVideoList, secUID, cursor, secSig)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", requester.UserAgent)
	client := &http.Client{}

	for {
		time.Sleep(100 * time.Millisecond)

		if tryTimes > 500 {
			wlog.Infof("[警告]获取视频列表尝试超过%d次仍然没有获得数据", tryTimes)
			tryTimes = 0
		}
		resp, err := client.Do(req)
		if err != nil {
			tryTimes++
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return nil, 0, err
		}

		if len(data.AwemeList) == 0 && data.StatusCode == 0 && data.MaxCursor == 0 {
			tryTimes++
			continue
		}

		break
	}

	maxCursor := int64(0)
	hasmore := false
	switch data.HasMore.(type) {
	case bool:
		hasmore = data.HasMore.(bool)
	}

	//抖音那边 这几个字段都正常时都有可能是已经读完了 has_more: true max_cursor: 1598057927000 min_cursor: 1600132116000 status_code: 0
	if hasmore && data.MaxCursor != 0 && len(data.AwemeList) != 0 {
		maxCursor = data.MaxCursor
	}

	return data.AwemeList, maxCursor, nil

}

type APIDouyinVideoExraInfo struct {
	ItemList []struct {
		CreateTime int64 `json:"create_time"`
		Video      struct {
			Vid string `json:"vid"`
		} `json:"video"`
	} `json:"item_list"`
}

func GetVideoExtraInfo(awemeid string) (*APIDouyinVideoExraInfo, error) {
	url := fmt.Sprintf("%s?item_ids=%s", define.GetVideoURI, awemeid)

	tryTimes := 10
	var resp *http.Response
	var err error
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("获取secid时创建请求失败:%s", err)
	}
	req.Header.Add("User-Agent", requester.UserAgent)

	for tryTimes > 0 {
		resp, err = client.Do(req)
		if err != nil {
			tryTimes--
			// 设置间隔为了防止两次调用时间间隔过短导致握手失败
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}
	if resp == nil {
		return nil, fmt.Errorf("[%s]请求视频信息失败: 超过重试次数", awemeid)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &APIDouyinVideoExraInfo{}
	if err := json.Unmarshal(body, info); err != nil {
		return info, fmt.Errorf("[%s]数据解析失败:%s", awemeid, err)
	}

	return info, nil
}

func GetSecSignature(shareURL string) string {
	sigChan := make(chan string, 1)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.UserAgent(requester.UserAgent),
	}

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocatorCancel()

	optsctx, optCancel := chromedp.NewContext(
		allocatorCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer optCancel()

	ctx, cancel := context.WithTimeout(optsctx, 10*time.Second)
	defer cancel()

	listenForNetworkEvent := func(ctx context.Context) {
		chromedp.ListenTarget(ctx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *network.EventRequestWillBeSent:
				req := ev.Request

				u, err := url.Parse(req.URL)
				if err != nil {
					wlog.Info("解析域名失败:", req.URL)
				}
				if u.Path == "/web/api/v2/aweme/post/" {
					sigChan <- u.Query().Get("_signature")
				}
			}
		})
	}
	listenForNetworkEvent(ctx)

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(shareURL),
	)
	if err != nil {
		wlog.Info("获取_signature失败:", err)
		return ""
	}

	sig := ""
	select {
	case <-ctx.Done():
		wlog.Info("获取signature超时了:", shareURL)
	case sig = <-sigChan:
	}

	return sig
}
