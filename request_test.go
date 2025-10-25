package scuter

import (
	"net/http"
	"testing"
	"time"

	"github.com/mdw-go/scuter/internal/should"
)

func TestReadUint64Header(t *testing.T) {
	should.So(t, ReadUint64Header(http.Header{"a": []string{"1"}}, "a"), should.Equal, uint64(1))
	should.So(t, ReadUint64Header(http.Header{"a": []string{"NaN"}}, "a"), should.Equal, uint64(0))
	should.So(t, ReadUint64Header(http.Header{"a": []string{"1"}}, "nope"), should.Equal, uint64(0))
}
func TestReadTimeHeader(t *testing.T) {
	date := time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "2006-01-02", "a"), should.Equal, date)
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"nope"}}, "2006-01-02", "a"), should.Equal, time.Time{})
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "nope", "a"), should.Equal, time.Time{})
	should.So(t, ReadTimeHeader(http.Header{"a": []string{"2025-01-02"}}, "2006-01-02", "nope"), should.Equal, time.Time{})
}

func TestReadPathElement(t *testing.T) {
	assertPathElement(t, "/path/users/value", "users", "value")
	assertPathElement(t, "/path/users/value/", "users", "value")
	assertPathElement(t, "/path/users/value/other/stuff", "users", "value")
	assertPathElement(t, "/not/found/", "users", "")
}
func assertPathElement(t *testing.T, raw, element, expected string) {
	should.So(t, ReadPathElement(raw, element), should.Equal, expected)
}

func TestReadNumericPathElement(t *testing.T) {
	assertNumericPathElement(t, "/path/users/123", "users", 123)
	assertNumericPathElement(t, "/path/users/123/", "users", 123)
	assertNumericPathElement(t, "/path/users/123/other/stuff", "users", 123)
	assertNumericPathElement(t, "/path/users/-123/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/abc/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/abc123/other/stuff", "users", 0)
	assertNumericPathElement(t, "/path/users/_123/other/stuff", "users", 0)
}
func assertNumericPathElement(t *testing.T, raw, element string, expected uint64) {
	should.So(t, ReadNumericPathElement(raw, element), should.Equal, expected)
}
