package main

type UpdateType int

const (
	UpdateTypeDaily UpdateType = iota + 1
	UpdateTypeNewly
)

type GetMode int

const (
	GetModeNewly GetMode = iota + 1
	GetModeOlder
)
