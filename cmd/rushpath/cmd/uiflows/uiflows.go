package uiflows

import (
	"github.com/sveniu/rushpath/internal/service"
)

// UI simply wraps...
type UI struct {
	Service *service.Service
}

func New() *UI {
	return &UI{
		Service: service.New(),
	}
}
