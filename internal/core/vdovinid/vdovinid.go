package vdovinid

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/edu/server/internal/core/model"
)

type SessionStorage interface {
	Create(accessTokenOfVdovinID string, scope model.SessionScope) (*model.Session, error)
	Get(accessToken string) (*model.Session, error)
	Refresh(refreshToken string) (*model.Session, error)
	Delete(accessToken string) error
}

type WebsocketManager interface {
	SendMessage(id uuid.UUID, message string) error
	CloseConnection(id uuid.UUID) error
}

type Authorizer struct {
	sessionStorage   SessionStorage
	websocketManager WebsocketManager
}

func NewAuthorizer(ss SessionStorage, wm WebsocketManager) *Authorizer {
	return &Authorizer{
		sessionStorage:   ss,
		websocketManager: wm,
	}
}

func (m *Authorizer) VdovinIDAccessToken(accessTokenOfVdovinID string, scope model.SessionScope) error {
	// Создаём сессию с edu_access_token и access_token
	session, err := m.sessionStorage.Create(accessTokenOfVdovinID, scope)
	if err != nil {
		return err
	}

	// Отправляем edu_access_token
	err = m.websocketManager.SendMessage(uuid.New(), session.AccessToken)
	if err != nil {
		return err
	}

	// Закрываем Websocket соединение
	err = m.websocketManager.CloseConnection(uuid.New())
	if err != nil {
		return err
	}

	return nil
}

var ErrInvalidPassword = errors.New("invalid password")
