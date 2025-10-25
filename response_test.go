package scuter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mdw-go/scuter/internal/should"
)

func TestResponseHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.Header("Content-Type", "testing-content-type"))
	should.So(t, recorder.Header().Get("Content-Type"), should.Equal, "testing-content-type")
}
func TestResponseContentType(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.ContentType("testing-content-type"))
	should.So(t, recorder.Header().Get("Content-Type"), should.Equal, "testing-content-type")
}
func TestResponseJSONContentType(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.JSONContentType())
	should.So(t, recorder.Header().Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
}
func TestResponseStatusCode(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.StatusCode(http.StatusTeapot))
	should.So(t, recorder.Code, should.Equal, http.StatusTeapot)
}
func TestResponseBytesBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.BytesBody([]byte("Hello, world!")))
	should.So(t, recorder.Body.String(), should.Equal, "Hello, world!")
}
func TestResponseJSONBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.JSONBody([]string{"a", "b", "c"}))
	should.So(t, recorder.Body.String(), should.Equal, `["a","b","c"]`)
}
func TestResponseJSONError(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.JSONError(Error{
		Fields:  []string{"field-1", "field-2"},
		ID:      42,
		Name:    "testing-error",
		Message: "testing error message",
	}))
	should.So(t, recorder.Body.String(), should.Equal,
		`{"errors":[{"fields":["field-1","field-2"],"id":42,"name":"testing-error","message":"testing error message"}]}`,
	)
}
func TestResponseBodyFromReader(t *testing.T) {
	recorder := httptest.NewRecorder()
	Flush(recorder, Response.BodyFromReader(strings.NewReader("Hello, world!")))
	should.So(t, recorder.Body.String(), should.Equal, "Hello, world!")
}
func TestResponseBodyWithAttachment(t *testing.T) {
	recorder := httptest.NewRecorder()

	Flush(recorder, Response.BodyWithAttachment("filename.txt", strings.NewReader("Hello, world!")))

	should.So(t, recorder.Code, should.Equal, http.StatusOK)
	should.So(t, recorder.Header().Get("Content-Type"), should.Equal, "text/plain; charset=utf-8")
	should.So(t, recorder.Body.String(), should.Equal, "Hello, world!")
}
func TestResponseWith(t *testing.T) {
	recorder := httptest.NewRecorder()

	Flush(recorder, Response.With(Response.StatusCode(http.StatusTeapot), Response.JSONContentType()))

	should.So(t, recorder.Code, should.Equal, http.StatusTeapot)
	should.So(t, recorder.Header().Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
}
func TestResponseBodyFromReadCloser_CloseCalled(t *testing.T) {
	recorder := httptest.NewRecorder()
	readCloser := &Closer{
		Reader: &Reader{
			Reader: strings.NewReader("Hello, world!"),
		},
	}

	Flush(recorder, Response.BodyFromReader(readCloser))

	should.So(t, recorder.Body.String(), should.Equal, "Hello, world!")
	should.So(t, readCloser.closed, should.Equal, 1)
}

type Reader struct {
	io.Reader
	readErr error
}

func (this *Reader) Read(p []byte) (n int, err error) {
	n, err = this.Reader.Read(p)
	if this.readErr != nil {
		err = this.readErr
	}
	return n, err
}

type Closer struct {
	*Reader
	readErr  error
	closeErr error
	closed   int
}

func (this *Closer) Close() error {
	this.closed++
	return this.closeErr
}
