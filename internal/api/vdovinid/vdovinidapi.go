package vdovinidapi

import (
	"github.com/ubahwin/edu/server/internal/core/vdovinid"
	"log"
)

type Group struct {
	authorizer *vdovinid.Authorizer
	log        *log.Logger
}

func New(log *log.Logger, authorizer *vdovinid.Authorizer) *Group {
	return &Group{
		authorizer: authorizer,
		log:        log,
	}
}
