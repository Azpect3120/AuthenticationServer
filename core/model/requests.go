package model

type CreateApplicationRequest struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
}

type ModifyApplicationRequest struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
}
