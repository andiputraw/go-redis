package model

type Car struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	MaxSpeed int    `json:"max_speed"`
	Color    string `json:"color"`
}