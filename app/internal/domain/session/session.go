package session

import (
	"time"

	"github.com/google/uuid"
)

const (
	sessionTTl = time.Hour
)

type Session struct {
	id         uuid.UUID
	login      string
	token      string
	createdAt  time.Time
	modifiedAt time.Time
}

func (s Session) Id() uuid.UUID {
	return s.id
}

func (s Session) Login() string {
	return s.login
}

func (s Session) Token() string {
	return s.token
}

func (s Session) IsExpired() bool {
	return s.createdAt.Add(sessionTTl).Before(time.Now())
}

func (s Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s Session) ModifiedAt() time.Time {
	return s.modifiedAt
}

func NewSession(
	login string,
	token string,
) (*Session, error) {
	return &Session{
		id:         uuid.New(),
		login:      login,
		token:      token,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
	}, nil
}

func NewSessionWithId(
	id uuid.UUID,
	login string,
	token string,
	createdAt time.Time,
	modifiedAt time.Time,
) (*Session, error) {
	return &Session{
		id:         id,
		login:      login,
		token:      token,
		createdAt:  createdAt,
		modifiedAt: modifiedAt,
	}, nil
}
