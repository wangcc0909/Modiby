package models

//easyjson
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//easyjson
type Address struct {
	City   string `json:"city"`
	Zone   string `json:"zone"`
	Number string `json:"number"`
}
