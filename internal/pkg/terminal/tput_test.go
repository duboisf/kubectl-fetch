package terminal_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type mockCmd struct {
	stdinPipe struct {
		actualInput strings.Builder
		nopCloser   nopCloser
		err         error
	}
	output struct {
		output string
		err    error
	}
}

func (m *mockCmd) StdinPipe() (io.WriteCloser, error) {
	return &m.stdinPipe.nopCloser, m.stdinPipe.err
}

func (m *mockCmd) Output() ([]byte, error) {
	return []byte(m.output.output), m.output.err
}

type tputFixture struct {
	actualName string
	actualArgs []string
	cmd        *mockCmd
}

func (t *tputFixture) newCommand(name string, args ...string) terminal.Cmd {
	t.actualName = name
	t.actualArgs = args
	return t.cmd
}

func newTPutFixture() *tputFixture {
	mockTPutRunner := &mockCmd{}
	mockTPutRunner.stdinPipe.nopCloser = nopCloser{
		Writer: &mockTPutRunner.stdinPipe.actualInput,
	}
	return &tputFixture{
		cmd: mockTPutRunner,
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

type mockWriter struct {
	err error
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return len(p), m.err
}

func TestTPut_Query(t *testing.T) {
	t.Parallel()
	t.Run("works", func(t *testing.T) {
		f := newTPutFixture()
		f.cmd.output.output = "erasing!\n"
		tput := terminal.NewTPut(f.newCommand)
		output, err := tput.Query("el")
		assert.Nil(t, err)
		assert.Equals(t, "erasing!", output)
		assert.Equals(t, "tput", f.actualName)
		assert.SliceEquals(t, []string{"-S"}, f.actualArgs)
		assert.Equals(t, "el\n", f.cmd.stdinPipe.actualInput.String())
	})

	t.Run("handles errors", func(t *testing.T) {
		f := newTPutFixture()
		f.cmd.stdinPipe.err = errors.New("e1")
		tput := terminal.NewTPut(f.newCommand)
		_, err := tput.Query("hi1")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "e1")
		f.cmd.stdinPipe.err = nil

		f.cmd.output.err = errors.New("e2")
		_, err = tput.Query("hi2")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "e2")
		f.cmd.output.err = nil

		f.cmd.stdinPipe.nopCloser.Writer = &mockWriter{err: errors.New("e3")}
		_, err = tput.Query("hi2")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "e3")
		f.cmd.stdinPipe.nopCloser.Writer = &strings.Builder{}

		f.cmd.stdinPipe.nopCloser.Writer = &mockWriter{err: errors.New("e4")}
		_, err = tput.Query(strings.Repeat(".", 4097))
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "e4")
	})
}

func TestTPut_QueryInt(t *testing.T) {
	t.Parallel()
	t.Run("works", func(t *testing.T) {
		f := newTPutFixture()
		f.cmd.output.output = "162\n"
		tput := terminal.NewTPut(f.newCommand)
		output, err := tput.QueryInt("cols")
		assert.Nil(t, err)
		assert.Equals(t, 162, output)
	})

	t.Run("handles errors", func(t *testing.T) {
		f := newTPutFixture()
		f.cmd.output.err = errors.New("an error")
		tput := terminal.NewTPut(f.newCommand)
		_, err := tput.QueryInt("i")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "an error")
	})
}
