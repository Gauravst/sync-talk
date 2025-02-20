package models

type AccessToken struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}
