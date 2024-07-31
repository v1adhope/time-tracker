package entities

type User struct {
	ID             string `json:"id" example:"1ef4f189-7b2a-6740-a609-370ed63a9fc7"`
	Surname        string `json:"surname" example:"Funk"`
	Name           string `json:"name" example:"Theresia"`
	Patronymic     string `json:"patronymic" example:"Cummerata-Thompson"`
	Address        string `json:"address" example:"53636 Gabrielle Mount"`
	PassportNumber string `json:"passportNumber" example:"3333 333333"`
}

type UserPagination struct {
	Limit  string
	Offset string
}

type UserFilter struct {
	ByID             string
	BySurname        string
	ByName           string
	ByPatronymic     string
	ByAddress        string
	ByPassportNumber string
}

type UserRepresentation struct {
	Pagination UserPagination
	Filter     UserFilter
}
