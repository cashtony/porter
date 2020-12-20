package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"path"
	"porter/define"
	"porter/wlog"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

func (b *BaiduClient) Upload(filePath, desc string) error {
	// 生成缩略图
	tsPath := fmt.Sprintf("./thumbsnails/%s/%s.jpg", b.Nickname, path.Base(filePath))
	err := genThumbnails(filePath, tsPath)
	if err != nil {
		return fmt.Errorf("缩略图生成失败:%s", err)
	}
	// 获取视频长度
	duration, err := getVideoDuration(filePath)
	if err != nil {
		return fmt.Errorf("获取视频长度失败:%s", err)
	}
	// 如果大于5分钟那么需要剪辑视频
	second := strings.Split(duration, ".")
	d, err := strconv.Atoi(second[0])
	if err != nil {
		return fmt.Errorf("视频长度转换错误: %s", err)
	}
	if int(d/60) >= 5 {
		wlog.Infof("视频长度[%s]超过全民小视频规定,将进行剪辑: \n", duration)
		dir := path.Dir(filePath)
		filenameall := path.Base(filePath)
		filesuffix := path.Ext(filePath)
		fileprefix := filenameall[0 : len(filenameall)-len(filesuffix)]
		newfilepath := fmt.Sprintf("%s/%s2%s", dir, fileprefix, filesuffix)
		err := cutVideoLength(filePath, newfilepath)
		if err != nil {
			return fmt.Errorf("视频剪辑发生错误: %s", err)
		}

		filePath = newfilepath
		duration, err = getVideoDuration(filePath)
		if err != nil {
			return fmt.Errorf("获取视频长度失败:%s", err)
		}
	}

	// 开始解密
	sv := fmt.Sprintf("%s%s", duration, b.sv)
	tv := Ase256(sv, b.k, b.v)

	return b.streamUpload(filePath, tsPath, tv, desc)
}

// 生成解密数据
func (b *BaiduClient) FetchSecretInfo() error {
	resp, err := b.client.Req("GET", define.GetQuanminInfo, nil, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return fmt.Errorf("请求全民用户数据失败, %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("请求上传的json解析出错: %s", err)
	}
	errCode := json.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := json.Get("errmsg").MustString()
		return fmt.Errorf("获取用户加密信息失败: %d, 消息: %s", errCode, errMsg)
	}

	dataJ := json.Get("data")
	b.k, err = dataJ.Get("k").String()
	if err != nil {
		return fmt.Errorf("获取用户加密k信息失败")
	}
	b.v, err = dataJ.Get("v").String()
	if err != nil {
		return fmt.Errorf("获取用户加密v信息失败")
	}
	b.sv, err = dataJ.Get("sv").String()
	if err != nil {
		return fmt.Errorf("获取用户加密sv信息失败")
	}
	b.Nickname, _ = dataJ.Get("name").String()

	return nil
}

func (b *BaiduClient) streamUpload(filepath, thunbsnailpath, tv, desc string) error {
	// 请求空间
	uploadID, mediaID, err := b.requestSpace()
	if err != nil {
		return fmt.Errorf("申请空间时出错:%s", err)
	}

	time.Sleep(1 * time.Second)
	// 发送视频
	tags, err := b.uploadPart(filepath, uploadID, mediaID)
	if err != nil {
		return fmt.Errorf("上传文件时出错:%s", err)
	}

	time.Sleep(1 * time.Second)

	// 发送完毕
	if err := b.uploadFinished(tags, uploadID, mediaID); err != nil {
		return fmt.Errorf("请求上传完毕消息时出错:%s", err)
	}

	time.Sleep(1 * time.Second)

	// 上传封面
	if err := b.uploadPoster(thunbsnailpath, uploadID, mediaID); err != nil {
		return fmt.Errorf("请求上传完毕消息时出错:%s", err)
	}

	time.Sleep(1 * time.Second)

	// 结束上传
	if err := b.videoPublish(mediaID, tv, desc); err != nil {
		return fmt.Errorf("请求发布视频错误:%s", err)
	}
	time.Sleep(2 * time.Second)

	return nil
}

func (b *BaiduClient) videoPublish(mediaID, tv, desc string) error {
	params := []struct {
		Title             string `json:"title"`
		PosterWidth       int    `json:"poster_width"`
		PosterHeight      int    `json:"poster_height"`
		MediaID           string `json:"media_id"`
		CoverUploadType   string `json:"cover_upload_type"`
		CosswiseCoverType int    `json:"crosswise_cover_type"`
		TV                string `json:"tv"`
		PublishTimer      int    `json:"publish_timer"`
	}{{Title: desc, TV: tv, PosterWidth: 720, PosterHeight: 1280, MediaID: mediaID}}

	infoData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("发布时的参数解析错误:%s", err)
	}

	post := map[string]string{
		"video_info": string(infoData),
	}
	resp, err := b.client.Req("POST", define.VideoPushlish, post, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return fmt.Errorf("发布视频的请求错误, %s", err)
	}
	defer resp.Body.Close()

	// [{"title":"漂亮的车子啊啊啊","poster_width":576,"poster_height":1024,"media_id":"mda-km9kthn7twcbggu5","cover_upload_type":"","crosswise_cover_type":0,"tv":"CvkzS8OqGJucWt084d7i/E2Fe/78gj8tNZiScP6AMgYN1BFo8TsoIoM94sOTMuqHv2Gcgyehvv3kvhyqv8q/nw==","publish_timer":0}]
	upJSON, upErr := simplejson.NewFromReader(resp.Body)
	if upErr != nil {
		return fmt.Errorf("json解析失败: %s", err)
	}
	errCode := upJSON.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := upJSON.Get("errmsg").MustString()
		return fmt.Errorf("发布过程发生错误: %d, 消息: %s", errCode, errMsg)
	}

	return nil
}

func (b *BaiduClient) uploadPoster(thumbsnailpath string, uploadID, mediaID string) error {
	bodyBuf := bytes.NewBuffer(make([]byte, 0))
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	fileWriter, err := bodyWriter.CreateFormFile("file", thumbsnailpath)
	if err != nil {
		return fmt.Errorf("创建multipart失败:%s", err)
	}

	f, err := os.Open(thumbsnailpath)
	if err != nil {
		return fmt.Errorf("打开文件失败:%s", err)
	}

	if _, err := io.Copy(fileWriter, f); err != nil {
		return fmt.Errorf("io.copy出错:%s", err)
	}

	bodyWriter.WriteField("is_crosswise_cover", "0")
	bodyWriter.WriteField("media_id", mediaID)
	bodyWriter.WriteField("upload_id", uploadID)

	header := map[string]string{
		"Content-Type":   bodyWriter.FormDataContentType(),
		"content-length": strconv.Itoa(bodyBuf.Len()),
	}

	resp, err := b.client.Req("POST", define.UploadPoster, bodyBuf, header)
	if err != nil {
		return fmt.Errorf("上传请求出现错误 %s", err)
	}
	defer resp.Body.Close()

	json, upErr := simplejson.NewFromReader(resp.Body)
	if upErr != nil {
		return fmt.Errorf("json解析失败: %s", err)
	}
	errCode := json.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := json.Get("errmsg").MustString()
		return fmt.Errorf("上传视频信息失败: %d, 消息: %s", errCode, errMsg)
	}

	return nil
}

func (b *BaiduClient) uploadFinished(tags []*ETag, uploadID, mediaID string) error {
	// [{"eTag":"c1be3d7e5884306312f3abec844f4156","partNumber":1}]
	tagData, err := json.Marshal(tags)
	if err != nil {
		return fmt.Errorf("tag数据解析错误:%s", err)
	}
	formData := map[string]string{
		"media_id":  mediaID,
		"upload_id": uploadID,
		"parts":     string(tagData),
	}
	resp, err := b.client.Req("POST", define.UploadFinished, formData, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return fmt.Errorf("请求上传结束失败, %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("解析上传结束数据失败: %s", err)
	}

	errCode := json.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := json.Get("errmsg").MustString()
		return fmt.Errorf("上传结束过程发生错误: %d, 消息: %s", errCode, errMsg)
	}

	return nil
}

func (b *BaiduClient) doUploadPart(data io.Reader, partNum int, mediaID, uploadID string) (string, error) {
	bodyBuf := bytes.NewBuffer(make([]byte, 0))
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%d.mp4"`, partNum))
	h.Set("Content-Type", "application/octet-stream")

	fileWriter, err := bodyWriter.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("创建multipart失败:%s", err)
	}

	if _, err := io.Copy(fileWriter, data); err != nil {
		return "", fmt.Errorf("io.copy出错:%s", err)
	}

	// 看是否要改成CreateFormField
	bodyWriter.WriteField("part_num", strconv.Itoa(partNum))
	bodyWriter.WriteField("media_id", mediaID)
	bodyWriter.WriteField("upload_id", uploadID)

	header := map[string]string{
		"Content-Type":   bodyWriter.FormDataContentType(),
		"content-length": strconv.Itoa(bodyBuf.Len()),
	}

	resp, err := b.client.Req("POST", define.UploadPart, bodyBuf, header)
	if err != nil {
		return "", fmt.Errorf("上传请求出现错误 %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		s, _ := ioutil.ReadAll(resp.Body) //把  body 内容读入字符串 s
		return "", fmt.Errorf("json解析失败 content:%s \n err: %s", string(s), err)
	}
	errCode := json.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := json.Get("errmsg").MustString()
		return "", fmt.Errorf("上传分片数据失败: %d, 消息: %s", errCode, errMsg)
	}

	etag, err := json.Get("data").Get("eTag").String()
	if err != nil {
		return "", fmt.Errorf("解析eTag时发生错误: %s", err)
	}

	return etag, nil
}

type ETag struct {
	Etag       string `json:"eTag"`
	PartNumber int    `json:"partNumber"`
}

func (b *BaiduClient) uploadPart(filepath string, uploadID, mediaID string) ([]*ETag, error) {
	bodyBuf := bytes.NewBuffer(make([]byte, 0))
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败:%s", err)
	}
	defer f.Close()

	tags := make([]*ETag, 0)
	// 每4M一个块
	s := make([]byte, 4*define.MB)

	partNum := 1

	for {
		switch nr, err := f.Read(s); true {

		case nr < 0:
			return nil, fmt.Errorf("从视频文件中读取数据出错:%s", err)
		case nr == 0: // EOF
			return tags, nil
		case nr > 0:
			blockReader := bytes.NewReader(s[:nr])
			etag, err := b.doUploadPart(blockReader, partNum, mediaID, uploadID)
			if err != nil {
				return nil, fmt.Errorf("上传第[%d]个分片数据时错误:%s", partNum, err)
			}
			tags = append(tags, &ETag{
				Etag:       etag,
				PartNumber: partNum,
			})

			partNum++
		}
	}

}
func (b *BaiduClient) requestSpace() (string, string, error) {
start:
	resp, err := b.client.Req("POST", define.UploadSpace, nil, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return "", "", fmt.Errorf("请求失败, %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("解析resp.Body失败: %s", err)
	}

	errCode := json.Get("errno").MustInt()
	if errCode != 0 {
		errMsg := json.Get("errmsg").MustString()
		return "", "", fmt.Errorf("获取上传空间失败: %d, 消息: %s", errCode, errMsg)
	}

	// 有时次数太多会触发保护机制
	need_pop, err := json.Get("data").Get("need_pop").Int()
	if err != nil {
		return "", "", fmt.Errorf("解析need_pop字段失败: %s", err)
	}
	if need_pop == 1 {
		wlog.Infof("[%s]当前操作触发了旋转验证,将弹窗手动验证", b.Nickname)
		if err := popCheack(b.BDUSS); err != nil {
			return "", "", fmt.Errorf("手动处理转验证失败:%s", err)
		}
		wlog.Infof("[%s]验证完成", b.Nickname)
		goto start
	}
	var uploadID, mediaID string
	spaceArray, err := json.Get("data").Get("upload").Array()
	if err != nil {
		return "", "", fmt.Errorf("解析upload字段失败: %s", err)
	}
	for _, single := range spaceArray {
		item := single.(map[string]interface{})
		// bucket = item["bucket"].(string)
		uploadID = item["upload_id"].(string)
		mediaID = item["media_id"].(string)
		// key = item["key"].(string)
	}

	return uploadID, mediaID, nil
}
