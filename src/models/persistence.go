package models

type DSN struct {
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
}
