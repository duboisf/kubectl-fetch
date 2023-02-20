package terminal_test

import (
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type MockFunc2[A any] struct {
	calls int
	output A
	err error
}

type mockTermInfo struct {
	mockQuery MockFunc2[string]
	mockQueryInt MockFunc2[int]
}

func (m *mockTermInfo) Query(capname ...string) (string, error) {
	m.mockQuery.calls++
	return m.mockQuery.output, m.mockQuery.err
}

func (m *mockTermInfo) QueryInt(capname string) (int, error) {
	m.mockQueryInt.calls++
	return m.mockQueryInt.output, m.mockQueryInt.err
}

func TestUI_NewWriter(t *testing.T) {
	assert.NotNil(t, terminal.NewUI(nil, nil, nil))
}
