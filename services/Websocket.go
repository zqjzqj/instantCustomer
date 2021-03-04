package services

import "time"

const(
	WsMessageTypeText = "text"
	WsMessageTypeImage = "image"
	WsMessageTypeVideo = "video"

	WsMessageStatusWaitSend = 0
	WsMessageStatusSend = 1
	WsMessageStatusFail = 2
	WsMessageStatusWaitRead = 3
	WsMessageStatusRead = 4

	EventOnChat = "OnChat"
)

type WsMessage struct {
	Id int64
	VisitorId uint64
	MAccountId uint64
	MchId uint64
	Data string
	Type string
	CreatedAt time.Time
	Status uint8
}