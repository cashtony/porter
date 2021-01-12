package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"porter/wlog"
	"strconv"
	"strings"
	"time"
)

var byteTable = "D6 28 3B 71 70 76 BE 1B A4 FE 19 57 5E 6C BC 21 B2 14 37 7D 8C A2 FA 67 55 6A 95 E3 FA 67 78 ED 8E 55 33 89 A8 CE 36 B3 5C D6 B2 6F 96 C4 34 B9 6A EC 34 95 C4 FA 72 FF B8 42 8D FB EC 70 F0 85 46 D8 B2 A1 E0 CE AE 4B 7D AE A4 87 CE E3 AC 51 55 C4 36 AD FC C4 EA 97 70 6A 85 37 6A C8 68 FA FE B0 33 B9 67 7E CE E3 CC 86 D6 9F 76 74 89 E9 DA 9C 78 C5 95 AA B0 34 B3 F2 7D B2 A2 ED E0 B5 B6 88 95 D1 51 D6 9E 7D D1 C8 F9 B7 70 CC 9C B6 92 C5 FA DD 9F 28 DA C7 E0 CA 95 B2 DA 34 97 CE 74 FA 37 E9 7D C4 A2 37 FB FA F1 CF AA 89 7D 55 AE 87 BC F5 E9 6A C4 68 C7 FA 76 85 14 D0 D0 E5 CE FF 19 D6 E5 D6 CC F1 F4 6C E9 E7 89 B2 B7 AE 28 89 BE 5E DC 87 6C F7 51 F2 67 78 AE B3 4B A2 B3 21 3B 55 F8 B3 76 B2 CF B3 B3 FF B3 5E 71 7D FA FC FF A8 7D FE D8 9C 1B C4 6A F9 88 B5 E5"

var cookies = map[string]string{
	"tt_webid":       "9efd6f2d6eb04a4ff16cb7bd34b995ba",
	"d_ticket":       "bbb35c121072ddb074b7dcf21e106c882c755",
	"multi_sids":     "111559650225%3Aa7206f8d4ae6ebbb711bd1166e808d7e",
	"odin_tt":        "d242373cdacb84292db6d5edab79a6001c6dd272650af09200fa2cd56bfa0fc518a3670015cca1e5f0815b87abf40dfa2070793b9346edf4c8de358d4162cda3",
	"n_mh":           "61zgAI6WwPy65JeI47NfRA6NPI6Ya_b--OtEUei81HQ",
	"MONITOR_WEB_ID": "1b53683b-59d1-4cbe-a14a-1e5e946dc113",
	"install_id":     "3254196246160136",
	"sid_guard":      "a7206f8d4ae6ebbb711bd1166e808d7e%7C1609137319%7C5184000%7CFri%2C+26-Feb-2021+06%3A35%3A19+GMT",
	"ttreq":          "1$1d4540d6189d3afdaed5e7686c1e94cc9293fc4d",
	"sid_tt":         "a7206f8d4ae6ebbb711bd1166e808d7e",
	"sessionid":      "a7206f8d4ae6ebbb711bd1166e808d7e",
	"uid_tt":         "d3c3cdaff70798f2f8cee3335d57b49e",
}

type APIPhoneUserBaseInfo struct {
	UID            string `json:"uid"`
	Nickname       string `json:"nickname"`
	Gender         int    `json:"gender"`
	Signature      string `json:"signature"`
	Birthday       string `json:"birthday"`
	IsVerified     bool   `json:"is_verified"`
	AwemeCount     int    `json:"aweme_count"`
	FollowerCount  int    `json:"follower_count"`
	TotalFavorited int    `json:"total_favorited"`
	UniqueID       string `json:"unique_id"`
	SecUID         string `json:"sec_uid"`
}

type APIDouyinSearch struct {
	Type         int    `json:"type"`
	HasMore      int    `json:"has_more"`
	InputKeyword string `json:"input_keyword"`
	UserList     []struct {
		UserInfo *APIPhoneUserBaseInfo `json:"user_info"`
	} `json:"user_list"`
}
type APIPhoneDouyinUser struct {
	User struct {
		*APIPhoneUserBaseInfo
		Avatar struct {
			URLList []string `json:"url_list"`
		} `json:"avatar_larger"`
		City     string `json:"city"`
		Location string `json:"location"`
		Province string `json:"province"`
	} `json:"user"`
}

func DouyinSearchKeyword(keyword string, cursor int) (*APIDouyinSearch, error) {
	ts := time.Now().Unix()
	_rticket := ts * 1000
	link := fmt.Sprintf("https://aweme.snssdk.com/aweme/v1/discover/search/?os_api=23&device_type=MI+5s&ssmix=a&manifest_version_code=130701&dpi=270&uuid=540000000264074&app_name=aweme&version_name=13.7.0&ts=%d&cpu_support64=false&app_type=normal&appTheme=dark&ac=wifi&host_abi=armeabi-v7a&update_version_code=13709900&channel=aweGW&_rticket=%d&device_platform=android&iid=3254196246160136&version_code=130700&cdid=4c41b427-81a9-4c81-a6a1-05f5d4bacaf4&openudid=f549f16478802fb0&device_id=70970912735&resolution=810*1440&os_version=6.0.1&language=zh&device_brand=Xiaomi&aid=1128", ts, _rticket)

	postValues := url.Values{}
	postValues.Set("cursor", strconv.Itoa(cursor))
	postValues.Set("keyword", keyword)
	postValues.Set("count", "20")
	postValues.Set("type", "1")
	postValues.Set("is_pull_refresh", "1")
	postValues.Set("query_correct_type", "1")

	postStr := postValues.Encode()
	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(postStr))
	if err != nil {
		return nil, err
	}
	addSecinfo(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	search := &APIDouyinSearch{}
	if err := json.Unmarshal(data, search); err != nil {
		return nil, err
	}

	return search, nil
}

func addSecinfo(req *http.Request) {
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v, Expires: time.Now().Add(365 * 24 * time.Hour)})
	}

	cookieStr := req.Header.Get("cookie")
	tsStr := req.URL.Query().Get("ts")
	ts, err := strconv.Atoi(tsStr)
	if err != nil {
		wlog.Error("时间转换错误:", tsStr)
		return
	}
	x := getXGon(req.URL.RawQuery, "", cookieStr, cookies["sessionid"])
	xgorgon := xGorgon(ts, strToByte(x))

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("X-Tt-Token", "004f1b1d23fdca8b781f865c86d7d35ef0016fd3a54652e4a1fcb9f3ab41affaaf16a83cf11ae28872ff40b100bab6a2af8f3a99a08938a82c14df51cb245319dac8450ceeb7d388be3e551c929f783ce89cf-1.0.0")
	req.Header.Add("X-Khronos", tsStr)
	req.Header.Add("X-Gorgon", xgorgon)
	req.Header.Add("User-Agent", "okhttp/3.10.0.1")
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
}

// UserInfoDetail 比搜索获得的数据中多出城市和省份
func NewPhoneDouyinUser(secUID string) (*APIPhoneDouyinUser, error) {
	ts := time.Now().Unix()
	_rticket := ts * 1000
	link := fmt.Sprintf("https://aweme.snssdk.com/aweme/v1/user/profile/other/?sec_user_id=%s&address_book_access=1&from=0&publish_video_strategy_type=2&user_avatar_shrink=188_188&user_cover_shrink=750_422&os_api=23&device_type=MI+5s&ssmix=a&manifest_version_code=130701&dpi=270&uuid=540000000264074&app_name=aweme&version_name=13.7.0&ts=%d&cpu_support64=false&app_type=normal&appTheme=dark&ac=wifi&host_abi=armeabi-v7a&update_version_code=13709900&channel=aweGW&_rticket=%d&device_platform=android&iid=3254196246160136&version_code=130700&cdid=4c41b427-81a9-4c81-a6a1-05f5d4bacaf4&openudid=f549f16478802fb0&device_id=70970912735&resolution=810*1440&os_version=6.0.1&language=zh&device_brand=Xiaomi&aid=1128", secUID, ts, _rticket)
	req, err := http.NewRequest(http.MethodPost, link, nil)
	if err != nil {
		return nil, err
	}
	addSecinfo(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	user := &APIPhoneDouyinUser{}
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

func MobileDouyinVideo(secUID string, cursor int64) ([]*APIDouyinVideo, int64, error) {
	if secUID == "" {
		return nil, 0, errors.New("secUID不能为空")
	}

	ts := time.Now().Unix()
	_rticket := ts * 1000
	link := fmt.Sprintf("https://aweme.snssdk.com/aweme/v1/aweme/post/?source=0&user_avatar_shrink=96_96&video_cover_shrink=248_330&publish_video_strategy_type=2&max_cursor=%d&sec_user_id=%s&count=30&is_order_flow=0&os_api=23&device_type=MI+5s&ssmix=a&manifest_version_code=130701&dpi=270&uuid=540000000264074&app_name=aweme&version_name=13.7.0&ts=%d&cpu_support64=false&app_type=normal&appTheme=dark&ac=wifi&host_abi=armeabi-v7a&update_version_code=13709900&channel=aweGW&_rticket=%d&device_platform=android&iid=3254196246160136&version_code=130700&cdid=4c41b427-81a9-4c81-a6a1-05f5d4bacaf4&openudid=f549f16478802fb0&device_id=70970912735&resolution=810*1440&os_version=6.0.1&language=zh&device_brand=Xiaomi&aid=1128", cursor, secUID, ts, _rticket)

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, 0, err
	}
	addSecinfo(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, 0, err
	}

	awemeResult := &struct {
		AwemeList  []*APIDouyinVideo `json:"aweme_list"`
		HasMore    int               `json:"has_more"`
		MaxCursor  int64             `json:"max_cursor"`
		MinCursor  int64             `json:"min_cursor"`
		StatusCode int               `json:"status_code"`
	}{}

	if err := json.Unmarshal(data, awemeResult); err != nil {
		return nil, 0, err
	}

	maxCursor := int64(0)

	if awemeResult.HasMore == 1 && awemeResult.MaxCursor != 0 && len(awemeResult.AwemeList) != 0 {
		maxCursor = awemeResult.MaxCursor
	}

	return awemeResult.AwemeList, maxCursor, nil
}
