package terminal_test

import (
	"testing"
	"time"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

// Test all methods of the Spinner struct
func TestSpinner(t *testing.T) {
	s := terminal.NewSpinner(100 * time.Millisecond)
	s.Spin()
	assert.Equals(t, s.String(), "⣻")
	s.Spin()
	assert.Equals(t, s.String(), "⣽")
	s.Spin()
	assert.Equals(t, s.String(), "⣾")
	s.Spin()
	assert.Equals(t, s.String(), "⣷")
	s.Spin()
	assert.Equals(t, s.String(), "⣯")
	s.Spin()
	assert.Equals(t, s.String(), "⣟")
	s.Spin()
	assert.Equals(t, s.String(), "⡿")
	s.Spin()
	assert.Equals(t, s.String(), "⢿")
}
