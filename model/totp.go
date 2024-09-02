package model

type WebUser struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}
