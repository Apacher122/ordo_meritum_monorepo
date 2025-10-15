package contexts

import "firebase.google.com/go/v4/auth"

type userContextKey string

type UserContext struct {
	Token  *auth.Token
	UID    string
	ApiKey string
}

const UserContextKey userContextKey = "userContextContext"
