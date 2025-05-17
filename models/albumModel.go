package models

// struct that represents an album
type Album struct {
	Id     int    `json:"Id"`
	Title  string `json:"Title"`
	Artist string `json:"Artist"`
	Price  string `json:"Price"`
}
