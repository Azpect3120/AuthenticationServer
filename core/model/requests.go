package model

type CreateApplicationRequest struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
}

type ModifyApplicationRequest struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
}

type ValidateUserRequest struct {
  Columns []string `json:"columns"`
  User UserData `json:"user"`
}
