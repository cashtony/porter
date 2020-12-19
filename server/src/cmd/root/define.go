package main

type UploadType int

const (
	UploadTypeDaily UploadType = iota + 1
	UploadTypeNewly
)

type GetMode int

const (
	GetModeNewly GetMode = iota + 1
	GetModeOlder
)

type BaiduUserStatus int

const (
	BaiduUserStatusDisable BaiduUserStatus = iota
	BaiduUserStatusNormal
)

type BindErr int

const (
	BindErrDouyinUser = iota + 1
	BindErrBdussWrong
	BindErrAlreadyBind
	BindErrSqlQuery
)
