package model

import "slices"

// Session is who made the request: fields from the access token payload. Id == 0 means anonymous.
type Session struct {
	Id   int64
	Role string
}

func (s *Session) IsAuthed() bool {
	return s != nil && s.Id > 0
}

func (s *Session) HasRole(roles ...string) bool {
	return s.IsAuthed() && slices.Contains(roles, s.Role)
}

// Tokens is a token pair issued on sign-in and refresh
type Tokens struct {
	Access  string
	Refresh string
}
