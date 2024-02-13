package entity

type Client struct {
	ID      int    `json:"id"`
	Name    string `json:"nome"`
	Limit   int    `json:"limite"`
	Balance int    `json:"saldo"`
}
