package models

type Book struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    Title      string `json:"title"`
    Author     string `json:"author"`
    Quantity   int    `json:"quantity"`
    IsReserved bool   `json:"isReserved"`
    Reserver   string `json:"reserver"`
}