package app

import (
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// App represents the main application with all dependencies
type App struct {
	storage *data.SessionData
}

// New creates a new App instance with all dependencies
func New() *App {
	return &App{
		storage: data.NewSessionData(),
	}
}

// GetStorage returns the storage instance
func (app *App) GetStorage() *data.SessionData {
	return app.storage
}
