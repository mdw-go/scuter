package http

import "github.com/mdw-go/scuter"

func errResponse(code int, err scuter.Error) scuter.ResponseOption {
	return scuter.Response.With(
		scuter.Response.StatusCode(code),
		scuter.Response.JSONError(err),
	)
}

var (
	errBadRequestInvalidID = scuter.Error{
		Fields:  []string{"id"},
		Name:    "invalid-id",
		Message: "The id was invalid or not supplied.",
	}
	errBadRequestInvalidJSON = scuter.Error{
		Fields:  []string{"body"},
		Name:    "malformed-request-payload",
		Message: "The body did not contain well-formed data and could not be properly deserialized.",
	}
	errMissingDueDate = scuter.Error{
		Fields:  []string{"due_date"},
		Name:    "missing-due-date",
		Message: "The due date is required.",
	}
	errMissingDetails = scuter.Error{
		Fields:  []string{"details"},
		Name:    "missing-details",
		Message: "The details of the task are required.",
	}
	errTaskTooHard = scuter.Error{
		Fields:  []string{"details"},
		ID:      12345,
		Name:    "task-too-hard",
		Message: "the specified task was deemed overly difficult",
	}
	errInternalServerError = scuter.Error{
		ID:      54321,
		Name:    "internal-server-error",
		Message: "Internal Server Error",
	}
)
