package authapi

type SessionManager interface {
	GenerateNewSession() (*SessionData, error)
}
