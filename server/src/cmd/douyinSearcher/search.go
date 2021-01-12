package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"porter/api"
	"porter/define"
	"porter/util"
	"strings"
	"time"
)

func IsSimilarInQuanmin(nickname string) (bool, error) {
	link := "https://quanmin.baidu.com/mvideo/api?api_name=newTabSearch"
	postValues := url.Values{}

	postValues.Set("newTabSearch", fmt.Sprintf("type=author&pn=1&query_word=%s&times=2", util.FilterSpecial(nickname)))

	postStr := postValues.Encode()
	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(postStr))
	if err != nil {
		return false, err
	}
	req.Header.Add("User-Agent", define.MobileUserAgent)
	cookie := http.Cookie{Name: "BDUSS", Value: "h5ZHVRT0FCbE5KbTlzcE84YlBDYVdLLUcxNks0LUFhZkdjcWR0QkFFWUlSUVZnRVFBQUFBJCQAAAAAAAAAAAEAAADRJ-omx7HLrtepvNIyMDEyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAi43V8IuN1fYW", Expires: time.Now().Add(180 * 24 * time.Hour)}
	req.AddCookie(&cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	user := &api.APIQuanminSearch{}
	if err := json.Unmarshal(data, user); err != nil {
		return false, err
	}

	for _, user := range user.NewTabSearch.Data.AuthorList {
		// if strings.Contains(user.Content.Nickname, nickname) {
		// 	return true, nil
		// }
		// if SimilarText(user.Content.Nickname, nickname) >= 70 {
		// 	return true, nil
		// }
		if user.Content.Daren != 0 {
			return true, nil
		}
	}

	return false, nil
}
