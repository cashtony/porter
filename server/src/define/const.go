package define

const (
	// 按awemeid获取视频信息 https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=6834090710124236043
	GetVideoURI = "https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/"
	// 按uri下载无水印视频 https://aweme.snssdk.com/aweme/v1/play/?video_id=v0200f7b0000brbp0h64tqbtfrkkjqlg&ratio=(720p,540p)&line=0
	GetVideoDownload = "https://aweme.snssdk.com/aweme/v1/play/"
	// 获取视频列表 https://www.iesdouyin.com/web/api/v2/aweme/post/?user_id=%s&sec_uid=&count=20&max_cursor=0&aid=1128&_signature=&dytk=
	GetVideoList = "https://www.iesdouyin.com/web/api/v2/aweme/post/"
	// 获取个人信息 https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=MS4wLjABAAAAgxcLUz9MzZW1VzK4Kt61HD1TKghYPLQwzGpYDKJvRwg
	// 访问用户的shareURL时会进行跳转,跳转的链接中带有sec_uid
	GetUserInfo = "https://www.iesdouyin.com/web/api/v2/user/info/"
)

// 百度接口, cookie中带有bduss即可
const (
	GetBaiduUserInfo = "https://quanmin.baidu.com/wise/video/pcpub/userinfo?isAjax=1"
)

var (
	TaskPushTopic     = "TaskPush"
	TaskFinishedTopic = "TaskFinished"
)
