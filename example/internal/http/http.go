package http

import (
	"net/http"

	"github.com/mdw-go/scuter/example/internal/app"
)

func New(logger app.Logger, application app.Handler) http.Handler {
	router := http.NewServeMux()
	router.Handle("PUT    /tasks", NewCreateTaskShell(logger, application))
	router.Handle("DELETE /tasks", NewDeleteTaskShell(logger, application))
	return router
}
