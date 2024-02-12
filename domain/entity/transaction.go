package entity

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	ClientID    int       `json:"cliente_id"`
	Date        time.Time `json:"realizada_em"`
}
