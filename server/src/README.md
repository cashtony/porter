# Porter后端简介
后端分为多个模块,各模块之间依赖NSQ来传递消息

cmd/root: 负责定时任务的发布,前端请求处理,数据库相关操作等工作.
cmd/node: 负责一些消耗流量的操作会在此模块进行,比如下载和上传等,可以在不同的机器上多开.
cmd/douyinSearcher:  负责抖音用户搜索相关任务,比如按关键字或分享链接搜索,得到的结果会发送到root进行处理

# 编译
```bash
cd cmd/xxx
go build .
```

# 运行参数
root模块需要指定数据库名和密码
```bash
.\root.exe --DBName=xxx --DBPassword=xxx
```

node模块需要指定NSQ地址,同时可以选择用户并发数量
```bash
.\node.exe --NSQLookupd=ip:port --NSQD=ip:port --thread=xx
```

其他参数请运行help进行查看