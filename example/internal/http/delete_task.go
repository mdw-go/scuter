package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mdw-go/scuter"
	"github.com/mdw-go/scuter/example/internal/app"
)

type DeleteTaskShell struct {
	logger  app.Logger
	handler app.Handler
}

func NewDeleteTaskShell(logger app.Logger, handler app.Handler) *DeleteTaskShell {
	return &DeleteTaskShell{
		logger:  logger,
		handler: handler,
	}
}

func (this *DeleteTaskShell) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	scuter.Flush(response, this.serveHTTP(request))
}
func (this *DeleteTaskShell) serveHTTP(request *http.Request) scuter.ResponseOption {
	query := request.URL.Query()
	id, err := strconv.ParseUint(query.Get("id"), 10, 64)
	if err != nil {
		return errResponse(http.StatusBadRequest, errBadRequestInvalidID)
	}

	command := app.DeleteTaskCommand{ID: id}
	this.handler.Handle(request.Context(), &command)

	switch {
	case command.Result.Error == nil:
		return nil
	case errors.Is(command.Result.Error, app.ErrTaskNotFound):
		return nil
	default:
		return errResponse(http.StatusInternalServerError, errInternalServerError)
	}
}
