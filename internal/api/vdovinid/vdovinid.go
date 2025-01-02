package vdovinidapi

import (
	"github.com/ubahwin/edu/server/internal/api"
	"github.com/ubahwin/edu/server/internal/core/model"
	"net/http"
)

type TokenReq struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

type TokenResp struct {
	Success bool   `json:"success"`
	Comment string `json:"comment,omitempty"`
}

func (req TokenReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) Token(_ *api.Context, req *TokenReq) (*TokenResp, int) {
	scope, err := model.ParseSessionScope(req.Scope)
	if err != nil {
		return &TokenResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusOK
	}

	err = g.authorizer.VdovinIDAccessToken(req.AccessToken, scope)
	if err != nil {
		return &TokenResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusOK
	}

	return &TokenResp{
		Success: true,
	}, http.StatusOK
}
