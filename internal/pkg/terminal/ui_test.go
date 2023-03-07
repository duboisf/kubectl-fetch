package terminal_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type TestFunc struct {
	calls int
}

type MockFunc1[A any] struct {
	TestFunc
	output A
}

type MockFunc2[A any, B any] struct {
	TestFunc
	a A
	b B
}

type mockProgressBar struct {
	increment          TestFunc
	setTotalIncrements TestFunc
	setWidth           TestFunc
	string             MockFunc1[string]
}

func (m *mockProgressBar) Increment(i int) {
	m.increment.calls++
}

func (m *mockProgressBar) SetTotalIncrements(i int) {
	m.setTotalIncrements.calls++
}

func (m *mockProgressBar) SetWidth(width int) {
	m.setWidth.calls++
}

func (m *mockProgressBar) String() string {
	m.string.calls++
	return m.string.output
}

type mockTermInfo struct {
	mockQuery    MockFunc2[string, error]
	mockQueryInt MockFunc2[int, error]
}

func (m *mockTermInfo) Query(capname ...string) (string, error) {
	m.mockQuery.calls++
	return m.mockQuery.a, m.mockQuery.b
}

func (m *mockTermInfo) QueryInt(capname string) (int, error) {
	m.mockQueryInt.calls++
	return m.mockQueryInt.a, m.mockQueryInt.b
}

func TestUI_NewUI(t *testing.T) {
	assert.NotNil(t, terminal.NewUI(nil, nil, nil, nil))
}

func TestUI_SetTotalKinds(t *testing.T) {
	pbar := &mockProgressBar{}
	termInfo := &mockTermInfo{}
	var stderr strings.Builder
	ui := terminal.NewUI(pbar, nil, termInfo, &stderr)
	_ = ui.SetTotalKinds(0)
}

func TestUI_Start(t *testing.T) {
	t.Parallel()
	t.Run("works", func(t *testing.T) {
		pbar := &mockProgressBar{}
		termInfo := &mockTermInfo{}
		var stderr strings.Builder
		spinner := terminal.NewSpinner(1 * time.Millisecond)
		ui := terminal.NewUI(pbar, spinner, termInfo, &stderr)
		var waitGroup sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		waitGroup.Add(1)
		go ui.Start(ctx, &waitGroup)
		updates := ui.SetTotalKinds(2)
		updates <- &terminal.GetResourcesUpdate{"deployment", 5}
		updates <- &terminal.GetResourcesUpdate{"services", 1}
		time.Sleep(10 * time.Millisecond)
		close(updates)
		waitGroup.Wait()
		t.Log(stderr.String())
		assert.Contains(t, stderr.String(), "Discovering kinds... found 2.")
		assert.Contains(t, stderr.String(), "Total resources found:    6")
	})
}
