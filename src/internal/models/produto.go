package models

type Produto struct {
	Id         uint    `gorm:"primaryKey" json:"id"`
	Nome       string	`json:"nome"`
	Descricao  string	`json:"descricao"`
	Preco      float64	`json:"preco"`
	Quantidade int		`json:"Quantidade"`
}
