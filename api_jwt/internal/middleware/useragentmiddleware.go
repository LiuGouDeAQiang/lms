package middleware

import "net/http"

type UserAgentMiddleware struct {
	JWTAuth string
}

func NewUserAgentMiddleware() *UserAgentMiddleware {
	return &UserAgentMiddleware{}
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		token1 := r.Header.Get("Authorization")
		m.JWTAuth = token1
		// Passthrough to next handler if need
		next(w, r)
	}
}
