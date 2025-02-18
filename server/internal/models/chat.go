package models

type ChatRoom struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ProfilePic string `json:"profilePic"`
	UserId     int    `json:"userId"`
}
