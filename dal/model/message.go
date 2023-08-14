package model

type Message struct {
	ID         uint   `json:"id,omitempty"`
	FromUserId uint   `json:"from_user_id"`
	ToUserId   uint   `json:"to_user_id"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}
