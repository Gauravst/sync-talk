package models

type ChatRoom struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	ProfilePic string `json:"profilePic"`
}
