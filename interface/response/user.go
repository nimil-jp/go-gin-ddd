package response

type UserLogin struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
