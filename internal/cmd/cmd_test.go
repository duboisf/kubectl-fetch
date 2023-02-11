package cmd_test

import (
	"strings"
	"testing"

	"github.com/bitfield/script"
	"github.com/duboisf/kubectl-getall/internal/cmd"
)

func TestPlugin_GetAll(t *testing.T) {
	t.Parallel()

	t.Run("returns the list of all resources in a namespace", func(t *testing.T) {
		var stderrWriter strings.Builder
		calls := 0
		execOutput := [][]string{
			{"services", "deployment"}, // returns the list unsorted
			{"service bar", "service baz"},
			{"deployment foo"},
		}
		exec := func(_ string) *script.Pipe {
			if calls >= len(execOutput) {
				t.Fatalf("only expected %d calls of exec function", len(execOutput))
				return nil
			}
			pipe := script.Echo(strings.Join(execOutput[calls], "\n"))
			calls++
			return pipe
		}
		plugin := cmd.GetAllPlugin{
			Exec:   exec,
			Stderr: &stderrWriter,
		}
		resources, err := plugin.GetAll()
		if err != nil {
			t.Fatal(err)
		}
		expectedLines := []string{"deployment foo", "service bar", "service baz"}
		if strings.Join(expectedLines, "\n") != strings.Join(resources, "\n") {
			t.Fatalf("expected:\n%s, actual:\n%s", expectedLines, resources)
		}
		expectedStderr := strings.Join([]string{
			"Getting all kubernetes api resources... found 2.",
			"Getting resources (1/2)",
			"deployment\r          \033[FGetting resources (2/2)",
			"services", // must sort the resource kinds
			"\rWaiting for results (1/2)\rWaiting for results (2/2)",
			"", // extra newline
		}, "\n")
		actualStderr := stderrWriter.String()
		if expectedStderr != actualStderr {
			t.Fatalf("unexpected stderr output\nexpected:\n%s\nactual:\n%s\n", expectedStderr, actualStderr)
		}
	})

	t.Run("returns an error if exec returns an error", func(t *testing.T) {
		var stderrWriter strings.Builder
		exec := func(cmdLine string) *script.Pipe {
			return script.Exec("thisisanerror")
		}
		plugin := cmd.GetAllPlugin{
			Exec:   exec,
			Stderr: &stderrWriter,
		}
		_, err := plugin.GetAll()
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		const expectedStderr = `could not get api resources:
exec: "thisisanerror": executable file not found in $PATH`
		if expectedStderr != err.Error() {
			t.Fatalf("unexpected stderr output\nexpected:\n%s\nactual:\n%s", expectedStderr, err.Error())
		}
	})
}
