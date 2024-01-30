package authapi

import (
	"time"

	"github.com/gofrs/uuid"
)

func NewSessionManagerImpl(sessionDurationSecs int64) SessionManager {
	return &SessionManagerImpl{
		SessionDurationSecs: sessionDurationSecs,
	}
}

type SessionManagerImpl struct {
	SessionDurationSecs int64
}

func (this SessionManagerImpl) GenerateNewSession() (*SessionData, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UnixMilli()
	expireAt := createdAt + (this.SessionDurationSecs * 1000)
	return &SessionData{
		UserId:    uid.String(),
		CreatedAt: createdAt,
		ExpireAt:  expireAt,
	}, nil
}
