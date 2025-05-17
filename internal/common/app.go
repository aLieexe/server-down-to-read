package common

import (
	"go-template/internal/models"
	"log/slog"
)

type Application struct {
	Logger     *slog.Logger
	Books      *models.BookModel
	Repository *models.RepositoryModel
}

var App *Application
