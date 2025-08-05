package domain

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created_at"`
}

type FollowUser struct {
	FollowID   string `json:"follow_id"`
	FollowedID string `json:"followed_id"`
}

type Tweet struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
