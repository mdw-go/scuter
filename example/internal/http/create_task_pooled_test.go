package http

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mdw-go/scuter"
	"github.com/mdw-go/scuter/example/internal/app"
	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
)

func TestCreateTaskFixture(t *testing.T) {
	gunit.Run(new(CreateTaskFixture), t)
}

type CreateTaskFixture struct {
	*gunit.Fixture
	*HTTPFixture
}

func (this *CreateTaskFixture) Setup() {
	this.HTTPFixture = &HTTPFixture{Fixture: this.Fixture}
	this.HTTPFixture.Setup()
}

func (this *CreateTaskFixture) TestInvalidJSONRequestBody() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.Body(strings.NewReader("invalid json")),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusBadRequest),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"body"},
				Name:    "malformed-request-payload",
				Message: "The body did not contain well-formed data and could not be properly deserialized.",
			}),
		),
	)
}
func (this *CreateTaskFixture) TestInvalidFields() {
	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(nil),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusUnprocessableEntity),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"due_date"},
				Name:    "missing-due-date",
				Message: "The due date is required.",
			}),
			scuter.Response.JSONError(scuter.Error{
				Fields:  []string{"details"},
				Name:    "missing-details",
				Message: "The details of the task are required.",
			}),
		),
	)
}
func (this *CreateTaskFixture) TestNoID() {
	this.app = func(v any) { v.(*app.CreateTaskCommand).Result.ID = 0 }

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusInternalServerError),
			scuter.Response.JSONError(scuter.Error{
				ID:      54321,
				Name:    "internal-server-error",
				Message: "Internal Server Error",
			}),
		),
	)
}
func (this *CreateTaskFixture) TestTaskTooHard() {
	this.app = func(v any) { v.(*app.CreateTaskCommand).Result.Error = app.ErrTaskTooHard }

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusTeapot),
			scuter.Response.JSONError(scuter.Error{
				ID:      12345,
				Fields:  []string{"details"},
				Name:    "task-too-hard",
				Message: "the specified task was deemed overly difficult",
			}),
		),
	)
}
func (this *CreateTaskFixture) TestHappyPath() {
	this.app = func(v any) {
		command := v.(*app.CreateTaskCommand)
		this.So(command.Details, should.Equal, "Details")
		command.Result.ID = 42
	}

	this.assertFullHTTP("PUT /tasks",
		scuter.Request.With(
			scuter.Request.JSONBody(map[string]any{
				"details":  "Details",
				"due_date": this.now,
			}),
		),
		scuter.Response.With(
			scuter.Response.StatusCode(http.StatusCreated),
			scuter.Response.JSONBody(map[string]any{
				"id":      42,
				"details": "Details",
			}),
		),
	)
}
