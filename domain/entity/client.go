package entity

type Client struct {
	ID      int    `json:"id"`
	Name    string `json:"nome"`
	Limit   int64  `json:"limite"`
	Balance int64  `json:"saldo"`
}
