package kafka
type UserLoginEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Time   string `json:"time"`
}