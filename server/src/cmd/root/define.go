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
	Disable BaiduUserStatus = iota
	Normal
)
