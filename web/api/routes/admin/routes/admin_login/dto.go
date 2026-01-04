package admin_login

type (
	input struct {
		Password string `json:"password"`
	}
	output struct {
		Token string `json:"token"`
	}
)
