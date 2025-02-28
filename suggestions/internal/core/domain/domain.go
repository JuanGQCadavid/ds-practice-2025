package domain

type Book struct {
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	Id     string `json:"maximunCapacty,omitempty"`
}
