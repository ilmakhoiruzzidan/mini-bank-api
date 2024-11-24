package models

type Customer struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsLoggedOut bool   `json:"is_logged_out"`
}
