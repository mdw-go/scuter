package scuter

import (
	"testing"

	"github.com/mdw-go/scuter/internal/should"
)

func TestError(t *testing.T) {
	err := Error{
		Fields:  []string{"hi"},
		ID:      1,
		Name:    "name",
		Message: "message",
	}
	should.So(t, err.Error(), should.Equal, "message")
}
