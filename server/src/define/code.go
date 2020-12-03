package define

const (
	Success = iota + 1000
	WrongPassword
	UpdateAccountErr
	IllegalToken
)

// 业务
const (
	ParamErr = iota + 2001
	QueryErr
	WrongContent // 绑定用户时的内容有错误
)
