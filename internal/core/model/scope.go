package model

import (
	"errors"
	"strings"
)

type SessionScope int64

func (scope SessionScope) IsAllowed(scopeAccess SessionScope) bool {
	return int64(scope)&int64(scopeAccess) != 0
}

func ParseSessionScope(scopeStr string) (SessionScope, error) {
	if strings.Contains(scopeStr, " ") {
		return 0, ErrInvalidScopeString
	}

	var sessionScope SessionScope

	for _, entry := range strings.Split(scopeStr, ",") {
		if entry, ok := sessionScopeEntryNames[entry]; !ok {
			return 0, ErrInvalidScopeString
		} else {
			sessionScope |= SessionScope(entry)
		}
	}

	return sessionScope, nil
}

type SessionScopeEntry int64

// Usage example: SessionScopeBasicInfoAccess | SessionScopePhoneAccess | SessionScopeEmailAccess
const (
	// SessionScopeAuthenticatorEntry - особое поле для приложения-аутентификатора
	SessionScopeAuthenticatorEntry SessionScopeEntry = 1 << 0

	// SessionScopeBasicInfoEntry - доступ к имени и фамилии
	SessionScopeBasicInfoEntry SessionScopeEntry = 1 << 1

	// SessionScopePhoneEntry - доступ к телефону
	SessionScopePhoneEntry SessionScopeEntry = 1 << 2

	// SessionScopeEmailEntry - доступ к email
	SessionScopeEmailEntry SessionScopeEntry = 1 << 3
)

var sessionScopeEntryNames = map[string]SessionScopeEntry{
	"authenticator": SessionScopeAuthenticatorEntry,
	"basic":         SessionScopeBasicInfoEntry,
	"phone":         SessionScopePhoneEntry,
	"email":         SessionScopeEmailEntry,
}

var ErrInvalidScopeString = errors.New("invalid scope string")
