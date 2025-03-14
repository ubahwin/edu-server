package vdovinid

import (
	"errors"
	"github.com/ubahwin/edu/server/internal/core/model"
)

type SessionStorage interface {
	Create(accessTokenOfVdovinID string, scope model.SessionScope) (*model.Session, error)
	Get(accessToken string) (*model.Session, error)
	Refresh(refreshToken string) (*model.Session, error)
	Delete(accessToken string) error
}

type WebsocketManager interface {
	SendMessage(id string, message string) error
	CloseConnection(id string) error
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

func (m *Authorizer) VdovinIDAccessToken(authID, accessTokenOfVdovinID string, scope model.SessionScope) error {
	// Создаём сессию с edu_access_token и access_token
	session, err := m.sessionStorage.Create(accessTokenOfVdovinID, scope)
	if err != nil {
		return err
	}

	// Отправляем edu_access_token
	err = m.websocketManager.SendMessage(authID, session.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

var ErrInvalidPassword = errors.New("invalid password")
