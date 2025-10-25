package scuter

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mdw-go/scuter/internal/should"
)

func TestQuery(t *testing.T) {
	request := NewTestRequest(t.Context(), http.MethodGet, "/target?a=1",
		Request.Query("a", "2"),
		Request.Query("b", "1"),
	)
	should.So(t, request.URL.Query()["a"], should.Equal, []string{"1", "2"})
	should.So(t, request.URL.Query()["b"], should.Equal, []string{"1"})
}
func TestHeader(t *testing.T) {
	request := NewTestRequest(t.Context(), http.MethodGet, "/target",
		Request.Header("a", "1"),
		Request.Header("a", "2"),
		Request.Header("b", "2"),
	)
	should.So(t, request.Header["A"], should.Equal, []string{"1", "2"})
	should.So(t, request.Header["B"], should.Equal, []string{"2"})
}
func TestBody(t *testing.T) {
	request := NewTestRequest(t.Context(), http.MethodGet, "/target",
		Request.Body(strings.NewReader("Hello, world!")),
	)
	buffer := bytes.NewBuffer(nil)
	_, _ = io.Copy(buffer, request.Body)
	should.So(t, buffer.String(), should.Equal, "Hello, world!")
}
func TestBodyErrorPanic(t *testing.T) {
	readCloser := &Closer{
		Reader: &Reader{
			readErr: errors.New("read error"),
			Reader:  strings.NewReader("Hello, world!"),
		},
	}
	defer func() {
		r, _ := recover().(error)
		should.So(t, r, should.NOT.BeNil)
	}()
	_ = NewTestRequest(t.Context(), http.MethodGet, "/target",
		Request.Body(readCloser),
	)
}
func TestJSONBody(t *testing.T) {
	request := NewTestRequest(t.Context(), http.MethodGet, "/target",
		Request.JSONBody(map[string]string{"a": "1"}),
	)
	buffer := bytes.NewBuffer(nil)
	_, _ = io.Copy(buffer, request.Body)
	should.So(t, buffer.String(), should.Equal, `{"a":"1"}`)
}
func TestJSONBodyErrorPanic(t *testing.T) {
	defer func() {
		r, _ := recover().(error)
		should.So(t, r, should.NOT.BeNil)
	}()
	_ = NewTestRequest(t.Context(), http.MethodGet, "/target",
		Request.JSONBody(make(chan int)),
	)
}
