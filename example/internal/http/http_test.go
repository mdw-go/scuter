package http

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
	"time"

	"github.com/mdw-go/scuter"
	"github.com/smarty/gunit"
	"github.com/smarty/gunit/assert/should"
)

func TestHTTPFixture(t *testing.T) {
	gunit.Run(new(HTTPFixture), t)
}

type HTTPFixture struct {
	*gunit.Fixture
	now time.Time
	ctx context.Context
	app func(any)
}

func (this *HTTPFixture) Handle(ctx context.Context, messages ...any) {
	this.So(ctx.Value("testing"), should.Equal, this.Name())
	for _, msg := range messages {
		this.app(msg)
	}
}
func (this *HTTPFixture) Setup() {
	this.now = time.Now().Truncate(time.Second)
	this.ctx = context.WithValue(this.T().Context(), "testing", this.Name())
	this.app = func(any) {}
}

func (this *HTTPFixture) serve(request *http.Request) *httptest.ResponseRecorder {
	if testing.Verbose() {
		requestDump, err := httputil.DumpRequest(request, true)
		this.So(err, should.BeNil)
		for line := range strings.SplitSeq(string(requestDump), "\n") {
			this.Println("> ", line)
		}
	}

	actual := httptest.NewRecorder()
	router := New(this, this)
	router.ServeHTTP(actual, request)

	if testing.Verbose() {
		this.Println()
		responseDump, err := httputil.DumpResponse(actual.Result(), true)
		this.So(err, should.BeNil)
		for line := range strings.SplitSeq(string(responseDump), "\n") {
			this.Println("< ", line)
		}
	}
	return actual
}
func (this *HTTPFixture) assertFullHTTP(route string, req scuter.RequestOption, res scuter.ResponseOption) {
	fields := strings.Fields(route)
	method, target := fields[0], fields[1]
	actual := this.serve(scuter.NewTestRequest(this.ctx, method, target, req))
	expected := httptest.NewRecorder()
	scuter.Flush(expected, res)

	this.So(actual.Code, should.Equal, expected.Code)
	this.So(actual.Header(), should.Equal, expected.Header())
	if strings.Contains(actual.Header().Get("Content-Type"), "json") {
		var actualBody, expectedBody any
		_ = json.Unmarshal(actual.Body.Bytes(), &actualBody)
		_ = json.Unmarshal(expected.Body.Bytes(), &expectedBody)
		this.So(actualBody, should.Equal, expectedBody)
	} else {
		this.So(actual.Body.String(), should.Equal, expected.Body.String())
	}
}
