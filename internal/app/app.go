// Package app provides the main application structure and HTTP handlers for the chess arbiter delegation generator.
// It manages the application state, handles HTTP requests, and coordinates between different packages.
package app

import (
	"eu.michalvalko.chess_arbiter_delegation_generator/internal/data"
)

// App represents the main application with all dependencies.
// It serves as the central coordinator for the application, managing storage and providing access to handlers.
type App struct {
	storage *data.SessionData // In-memory storage for session data (arbiters, leagues, etc.)
}

// New creates a new App instance with all dependencies initialized.
// It sets up a new SessionData storage instance ready for use.
// Returns a pointer to a new App instance.
func New() *App {
	return &App{
		storage: data.NewSessionData(),
	}
}

// GetStorage returns the storage instance for external access.
// This allows other packages to access the session storage if needed.
// Returns a pointer to the SessionData instance.
func (app *App) GetStorage() *data.SessionData {
	return app.storage
}
