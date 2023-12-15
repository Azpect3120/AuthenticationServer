package model

type CreateApplicationRequest struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
}
