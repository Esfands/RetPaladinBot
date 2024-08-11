package isbanned

type User struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Banned    bool   `json:"banned"`
	BanReason string `json:"banReason"`
}
