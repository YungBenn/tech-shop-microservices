package entity

type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}