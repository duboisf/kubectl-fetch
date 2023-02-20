package kubectl

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Cmd is an interface for exec.Cmd to make unit testing easier.
type Cmd interface {
	Output() ([]byte, error)
}

type CommandContext[C Cmd] func(ctx context.Context, name string, args ...string) C

type Kubectl[C Cmd] struct {
	commandContext CommandContext[C]
}

func New[C Cmd](newCommandContext CommandContext[C]) *Kubectl[C] {
	return &Kubectl[C]{
		commandContext: newCommandContext,
	}
}

// ListApiResources returns a list of api resource names. If `namespaced` is
// true, then only resources that live in namespaces are returned, otherwise
// only resources that are global (non-namespaced) will be returned.
// Note: whatever the value of `namespaced`, the events* resource is filtered
// out from the results.
func (k *Kubectl[C]) ListApiResources(ctx context.Context, namespaced bool) ([]string, error) {
	// cmd := fmt.Sprintf("kubectl api-resources --verbs=list --namespaced=%t -o name", namespaced)
	namespacedString := strconv.FormatBool(namespaced)
	cmd := k.commandContext(ctx, "kubectl", "api-resources", "--verbs=list", "--namespaced="+string(namespacedString), "-o", "name")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resourceKinds := splitFilterAndSort(string(output))
	sort.Strings(resourceKinds)
	return resourceKinds, nil
}

// GetNamespacedResources returns the resouces in the given namespace.
func (k *Kubectl[C]) GetNamespacedResources(ctx context.Context, namespace, kind string) ([]string, error) {
	cmd := k.commandContext(ctx, "kubectl", "--namespace="+namespace, "get", "--show-kind", "--ignore-not-found", "-o", "name", kind)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("could not run kubectl:\n%s", exitErr.Stderr)
		}
		return nil, err
	}
	return splitFilterAndSort(string(output)), nil
}

// GetResources returns non-namespaced resources.
func (k *Kubectl[C]) GetResources(ctx context.Context, kind string) ([]string, error) {
	cmd := k.commandContext(ctx, "kubectl", "get", "--show-kind", "--ignore-not-found", "-o", "name", kind)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("could not run kubectl command:\n%s",
				strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}
	if len(output) == 0 {
		return nil, nil
	}
	return splitFilterAndSort(string(output)), nil
}

var eventsRegex = regexp.MustCompile(`^events(\.events\.k8s.io)?$`)

func splitFilterAndSort(output string) []string {
	var filteredLines []string
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if !eventsRegex.MatchString(line) {
			filteredLines = append(filteredLines, line)
		}
	}
	sort.Strings(filteredLines)
	return filteredLines
}
