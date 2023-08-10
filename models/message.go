package models

import (
	"time"
)

type Message struct {
	Id         uint      `json:"id,omitempty"`
	FromUserId uint      `json:"from_user_id,omitempty"`
	ToUserId   uint      `json:"to_user_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
}

type MessageChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
