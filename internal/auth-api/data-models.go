package authapi

type SessionData struct {
	UserId    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
	ExpireAt  int64  `json:"expire_at"`
}
