package models

type AccessToken struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}
