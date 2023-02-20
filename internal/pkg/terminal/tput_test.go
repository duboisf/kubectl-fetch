package terminal_test

import (
	"io"
	"strings"
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type mockCmd struct {
	start struct {
		actualInput strings.Builder
		err         error
	}
	output struct {
		output string
		err    error
	}
}

func (m *mockCmd) StdinPipe() (io.WriteCloser, error) {
	return nopCloser{&m.start.actualInput}, m.start.err
}

func (m *mockCmd) Output() ([]byte, error) {
	return []byte(m.output.output), m.output.err
}

type tputFixture struct {
	actualName string
	actualArgs []string
	cmd     *mockCmd
}

func (t *tputFixture) newCommand(name string, args ...string) terminal.Cmd {
	t.actualName = name
	t.actualArgs = args
	return t.cmd
}

func newTPutFixture() *tputFixture {
	mockTPutRunner := &mockCmd{}
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

func TestTPut_Query(t *testing.T) {
	f := newTPutFixture()
	f.cmd.output.output = "erasing!\n"
	tput := terminal.NewTPut(f.newCommand)
	output, err := tput.Query("el")
	assert.Nil(t, err)
	assert.Equals(t, "erasing!", output)
	assert.Equals(t, "tput", f.actualName)
	assert.SliceEquals(t, []string{"-S"}, f.actualArgs)
	assert.Equals(t, "el\n", f.cmd.start.actualInput.String())
}

func TestTPut_QueryInt(t *testing.T) {
	f := newTPutFixture()
	f.cmd.output.output = "162\n"
	tput := terminal.NewTPut(f.newCommand)
	output, err := tput.QueryInt("cols")
	assert.Nil(t, err)
	assert.Equals(t, 162, output)
}
