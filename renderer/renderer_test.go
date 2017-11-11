package renderer_test

import (
	"testing"
	"github.com/donutloop/xserverkit/renderer"
)

func TestNew(t *testing.T) {
	go func() {
		defer func() {
			if v := recover(); v != nil {
				t.Error(v)
			}
		}()

		renderer.New(nil)
		renderer.New(&renderer.Options{})
	}()
}
