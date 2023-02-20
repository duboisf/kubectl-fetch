package kubectl_test

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/pkg/kubectl"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type mockCmd struct {
	calls  int
	err    error
	output []string
}

func (m *mockCmd) Output() ([]byte, error) {
	m.calls++
	if m.err != nil {
		return nil, m.err
	}
	return []byte(m.output[m.calls-1]), nil
}

type fixture struct {
	actualName string
	actualArgs []string
	cmd        *mockCmd
	kubectl    *kubectl.Kubectl[*mockCmd]
}

func newFixture(cmd *mockCmd) *fixture {
	f := &fixture{}
	f.kubectl = kubectl.New(func(_ context.Context, name string, args ...string) *mockCmd {
		f.actualName = name
		f.actualArgs = args
		if cmd.output == nil {
			cmd.output = []string{}
		}
		return cmd
	})
	return f
}

func TestKubectl_GetApiResources(t *testing.T) {
	t.Parallel()
	t.Run("returns the list of namespaced api resources", func(t *testing.T) {
		cmd := &mockCmd{
			output: []string{
				"services\ndeployment\n", // returns the list unsorted
			},
		}
		f := newFixture(cmd)
		resources, err := f.kubectl.ListApiResources(context.Background(), true)
		if err != nil {
			t.Fatal(err)
		}
		expectedLines := []string{"deployment", "services"}
		if strings.Join(expectedLines, "\n") != strings.Join(resources, "\n") {
			t.Fatalf("expected:\n%s, actual:\n%s", expectedLines, resources)
		}
	})

	t.Run("filters out events from api resources", func(t *testing.T) {
		cmd := &mockCmd{
			output: []string{"services\nevents\nevents.events.k8s.io\ndeployment\n"},
		}
		f := newFixture(cmd)
		resources, err := f.kubectl.ListApiResources(context.Background(), true)
		if err != nil {
			t.Fatal(err)
		}
		expectedLines := []string{"deployment", "services"}
		if strings.Join(expectedLines, "\n") != strings.Join(resources, "\n") {
			t.Fatalf("expected:\n%s, actual:\n%s", expectedLines, resources)
		}
	})

	t.Run("returns an error if exec returns an error", func(t *testing.T) {
		cmd := &mockCmd{err: errors.New("this is an error")}
		f := newFixture(cmd)
		_, err := f.kubectl.ListApiResources(context.Background(), true)
		assert.NotNil(t, err)
		assert.Equals(t, "this is an error", err.Error())
	})
}

func TestKubectl_GetNamespacedResource(t *testing.T) {
	cmd := &mockCmd{output: []string{"resourceB\nresourceA\n"}}
	f := newFixture(cmd)
	actualResources, err := f.kubectl.GetNamespacedResources(context.Background(), "<ns>", "<res>")
	assert.Nil(t, err)
	assert.SliceEquals(t, []string{"resourceA", "resourceB"}, actualResources)
	assert.Equals(t, "kubectl", f.actualName)
	expectedArgs := []string{"--namespace=<ns>", "get", "--show-kind", "--ignore-not-found", "-o", "name", "<res>"}
	assert.SliceEquals(t, expectedArgs, f.actualArgs)
}

func TestKubectl_GetResources(t *testing.T) {
	t.Parallel()
	t.Run("works", func(t *testing.T) {
		cmd := &mockCmd{output: []string{"foo\nbar\nbaz\n"}}
		f := newFixture(cmd)
		actualResources, err := f.kubectl.GetResources(context.Background(), "deployment")
		assert.Nil(t, err)
		assert.SliceEquals(t, []string{"bar", "baz", "foo"}, actualResources)
		expectedArgs := []string{"get", "--show-kind", "--ignore-not-found", "-o", "name", "deployment"}
		assert.Equals(t, "kubectl", f.actualName)
		assert.SliceEquals(t, expectedArgs, f.actualArgs)

	})

	t.Run("returns stderr when there's an error", func(t *testing.T) {
		cmd := &mockCmd{output: []string{"foo\nbar\nbaz\n"}}
		cmd.err = &exec.ExitError{Stderr: []byte("some error")}
		f := newFixture(cmd)
		_, err := f.kubectl.GetResources(context.Background(), "deploy")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "some error")
	})
}
