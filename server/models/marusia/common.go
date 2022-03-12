package marusia

type Session struct {
	SessionID   string      `json:"session_id"`
	UserID      string      `json:"user_id"`
	SkillID     string      `json:"skill_id"`
	New         bool        `json:"new"`
	MessageID   int         `json:"message_id"`
	User        User        `json:"user"`
	Application Application `json:"application"`
}

type Meta struct {
	ClientID   string     `json:"client_id"`
	Locale     string     `json:"locale"`
	TimeZone   string     `json:"timezone"`
	Interfaces Interfaces `json:"interfaces"`
}

type Interfaces struct {
	Screen interface{} `json:"screen"`
}

type User struct {
	UserID string `json:"user_id"`
}

type Application struct {
	ApplicationType string `json:"application_type"`
	ApplicationID   string `json:"application_id"`
}
