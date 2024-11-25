package models

import "time"

type RevokedToken struct {
	Token      string    `json:"token"`
	LogoutTime time.Time `json:"logout_time"`
}
