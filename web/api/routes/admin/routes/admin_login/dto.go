package admin_login

import "time"

type (
	input struct {
		Password string `json:"password"`
	}
	output struct {
		Token  string    `json:"token"`
		Expiry time.Time `json:"expiry"`
	}
)
