package event

type WelcomeUser struct {
	UserId string `json:"user_id"`
	Email string `json:"email"`
}