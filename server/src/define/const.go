package define

const (
	// 按awemeid获取视频信息 https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=6834090710124236043
	GetVideoURI = "https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/"
	// 按uri下载无水印视频 https://aweme.snssdk.com/aweme/v1/play/?video_id=v0200f7b0000brbp0h64tqbtfrkkjqlg&ratio=(720p,540p)&line=0
	GetVideoDownload = "https://aweme.snssdk.com/aweme/v1/play/"
	// 获取视频列表 https://www.iesdouyin.com/web/api/v2/aweme/post/?user_id=%s&sec_uid=&count=20&max_cursor=0&aid=1128&_signature=&dytk=
	// 2020.12.8日抖音更新了api参数 https://www.iesdouyin.com/web/api/v2/aweme/post/?sec_uid=MS4wLjABAAAAgxcLUz9MzZW1VzK4Kt61HD1TKghYPLQwzGpYDKJvRwg&count=21&max_cursor=0&aid=1128&_signature=VzqsHwAACPa.4I5i5RgdPVc6rA&dytk=
	GetVideoList = "https://www.iesdouyin.com/web/api/v2/aweme/post/"
	// 获取个人信息 https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=MS4wLjABAAAAgxcLUz9MzZW1VzK4Kt61HD1TKghYPLQwzGpYDKJvRwg
	// 访问用户的shareURL时会进行跳转,跳转的链接中带有sec_uid
	GetUserInfo = "https://www.iesdouyin.com/web/api/v2/user/info/"
)

// 百度接口, cookie中带有bduss即可
const (
	GetBaiduBaseInfo = "https://pan.baidu.com/api/loginStatus?clienttype=5"
	// 加密信息获取
	GetQuanminInfo   = "https://quanmin.baidu.com/wise/video/pcpub/userinfo"
	GetQuanminInfoV2 = "https://quanmin.baidu.com/appui/user/mine?api_name=mine"
	// POST https://quanmin.baidu.com/mvideo/api?api_name=userprofilesubmit  api_name=userprofile
	// form  userprofilesubmit    nickname=超级码力366&user_type=ugc // 修改名称
	// form：userprofile method=get&user_type=ugc 获取用户是否能改名之类的信息
	QuanminAPI = "https://quanmin.baidu.com/mvideo/api"
	// 全民视频数据(包含剩余钻石数量) POST https://quanmin.baidu.com/appui/user/mine?api_name=mine
	// T豆查询 https://sv.baidu.com/liveserver/exchange/record?pn=1&rn=10&orderType=1&client_type=2
	UploadSpace    = "https://quanmin.baidu.com/wise/video/pcpub/getuploadid?video_num=1"
	UploadPart     = "https://quanmin.baidu.com/wise/video/pcpub/uploadvideopart"
	UploadFinished = "https://quanmin.baidu.com/wise/video/pcpub/finishupload"
	UploadPoster   = "https://quanmin.baidu.com/wise/video/pcpub/uploadposter"
	VideoPushlish  = "https://quanmin.baidu.com/wise/video/pcpub/publishvideo"

	// 设置头像
	SetPortrait = "https://passport.baidu.com/v2/sapi/center/setportrait"
)

var (
	TaskPushTopic     = "TaskPush"
	TaskFinishedTopic = "TaskFinished"

	TaskChangeInfoTopic = "TaskChangeInfo"

	TaskParseVideoTopic       = "TaskParseVideo"
	TaskParseVideoResultTopic = "TaskParseVideoResult"
)

const (
	// B byte
	B = (int64)(1 << (10 * iota))
	// KB kilobyte
	KB
	// MB megabyte
	MB
	// GB gigabyte
	GB
	// TB terabyte
	TB
	// PB petabyte
	PB
)

type ParseVideoType int

const (
	ParseVideoTypeOnePage ParseVideoType = iota + 1
	ParseVideoTypeAll
)
