package entity

type Client struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"nome,omitempty"`
	Limit   int    `json:"limite,omitempty"`
	Balance int    `json:"saldo,omitempty"`
}
