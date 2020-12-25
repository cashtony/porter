package define

const (
	Success = iota + 1000
	WrongPassword
	UpdateAccountErr
	IllegalToken
)

// 业务
const (
	ParamErr     = iota + 2001
	QueryDataErr // 数据库中查找数据失败
	WrongContent // 绑定用户时的内容有错误
	CannotBind
	AlreadyUpdating
	AlreadyBind
	WrongDouyinShareURL
	CannotDelete
)
