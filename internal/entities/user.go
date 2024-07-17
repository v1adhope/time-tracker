package entities

type User struct {
	ID             string `json:"id,omitempty"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	PassportNumber string `json:"passportNumber,omitempty"`
}

type UserPagination struct {
	Limit  string
	Offset string
}

type UserFilter struct {
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
