package cmd_test

import (
	"context"
	"io/fs"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/duboisf/kubectl-fetch/internal/cmd"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

var cmdNamespacedResources string

type mockFetcher struct{
	err error
	resources []string
}

func (m *mockFetcher) Fetch(ctx context.Context) ([]string, error) {
	return m.resources, m.err
}

type mockStarter struct {
	calls              int
	ungracefulShutdown bool
}

func (m *mockStarter) Start(ctx context.Context, wg *sync.WaitGroup) {
	m.calls++
	if !m.ungracefulShutdown {
		wg.Done()
	}
}

type mockFileInfo struct {
	mode fs.FileMode
}

func (f *mockFileInfo) Name() string {
	panic("not implemented") // TODO: Implement
}

func (f *mockFileInfo) Size() int64 {
	panic("not implemented") // TODO: Implement
}

func (f *mockFileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f *mockFileInfo) ModTime() time.Time {
	panic("not implemented") // TODO: Implement
}

func (f *mockFileInfo) IsDir() bool {
	panic("not implemented") // TODO: Implement
}

func (f *mockFileInfo) Sys() any {
	panic("not implemented") // TODO: Implement
}

type mockStdout struct {
	builder  strings.Builder
	fileInfo mockFileInfo
}

func (m *mockStdout) Write(p []byte) (n int, err error) {
	return m.builder.Write(p)
}

func (m *mockStdout) Stat() (fs.FileInfo, error) {
	return &m.fileInfo, nil
}

func TestCmd_Run(t *testing.T) {
	t.Parallel()
	t.Run("works", func(t *testing.T) {
		plugin := &mockFetcher{resources: []string{"deployment/foo"}}
		ui := &mockStarter{}
		var stderr strings.Builder
		stdout := &mockStdout{}
		stdout.fileInfo.mode = fs.ModeCharDevice
		cmd, err := cmd.NewCmd(plugin, stdout, &stderr, ui)
		assert.Nil(t, err)
		err = cmd.Run(context.Background())
		assert.Nil(t, err)
		assert.Equals(t, 1, ui.calls)
		assert.Contains(t, stdout.builder.String(), "deployment/foo")
	})

	t.Run("returns after a timeout period if the ui doesn't close", func(t *testing.T) {
		plugin := &mockFetcher{}
		ui := &mockStarter{ungracefulShutdown: true}
		var stderr strings.Builder
		stdout := &mockStdout{}
		stdout.fileInfo.mode = fs.ModeCharDevice
		cmd, err := cmd.NewCmd(plugin, stdout, &stderr, ui)
		assert.Nil(t, err)
		cmd.UIStopTimeout = 1 * time.Millisecond
		err = cmd.Run(context.Background())
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("displays a message to stderr when no resources were found", func(t *testing.T) {
		plugin := &mockFetcher{resources: nil}
		ui := &mockStarter{}
		var stderr strings.Builder
		stdout := &mockStdout{}
		cmd, err := cmd.NewCmd(plugin, stdout, &stderr, ui)
		assert.Nil(t, err)
		err = cmd.Run(context.Background())
		assert.Nil(t, err)
		assert.Contains(t, stderr.String(), "No resources found.")
	})
}
