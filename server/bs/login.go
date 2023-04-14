package bs

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"pass"`
}

type LoginResponse struct{}

type UserInfoResponse struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"admin"`
}

type ChangePasswordRequest struct {
	OldPass string `json:"old_pass"`
	NewPass string `json:"new_pass"`
}
