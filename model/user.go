package model

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
