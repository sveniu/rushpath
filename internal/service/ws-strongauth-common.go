package service

type totpPromptType int

const (
	newDevices totpPromptType = iota
	everyLogin
)
