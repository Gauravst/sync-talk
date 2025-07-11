package models

type ChatRoom struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Members     int    `json:"members"`
	Private     bool   `json:"private"`
	Code        string `json:"code"`
	Description string `json:"description"`
	UserId      int    `json:"userId"`
}
