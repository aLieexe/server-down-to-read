package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() http.Handler {
	router := chi.NewRouter()

	router.Post("/repo", postRepositoryHandler)
	router.Get("/repo/{repoId}", getRepoByIdHandler)
	router.Get("/repo/{repoId}/books", getBooksByRepoIdHandler)
	router.Post("/repo/{repoId}/books", postBookToRepoHandler)

	return router
}
