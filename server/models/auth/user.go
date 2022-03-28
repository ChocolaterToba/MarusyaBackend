package auth

type User struct {
	UserID        uint64
	ApplicationID string
	SessionID     string
	Username      string
}
