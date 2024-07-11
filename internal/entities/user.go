package entity

type User struct {
	ID             string `json:"-"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	PassportNumber string `json:"passportNumber"`
}
