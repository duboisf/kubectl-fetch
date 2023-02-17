package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bitfield/script"
)

// getPipeError helps returns the output of a pipe that has an error
func getPipeError(pipe *script.Pipe) error {
	// Need to discard the error on the pipe otherwise all
	// attemps to get the output are no-ops
	pipe.SetError(nil)
	output, _ := pipe.String()
	return fmt.Errorf("%s", output)
}

type Exec func(cmdLine string) *script.Pipe

type GetAllPlugin struct {
	Exec   Exec
	Stderr io.Writer
}

// New returns a new GetAllPlugin ready to be used
func New() *GetAllPlugin {
	return &GetAllPlugin{
		Exec:   script.Exec,
		Stderr: os.Stderr,
	}
}

func (g *GetAllPlugin) writeStderr(format string, a ...any) {
	fmt.Fprintf(g.Stderr, format, a...)
}

func (g *GetAllPlugin) GetAll() ([]string, error) {
	g.writeStderr("Getting all kubernetes api resources...")
	apiResourcesPipe := g.Exec("kubectl api-resources --verbs=list --namespaced -o name")
	eventsRegex := regexp.MustCompile(`^events(\.events\.k8s.io)?$`)
	resourceKinds, err := apiResourcesPipe.RejectRegexp(eventsRegex).Slice()
	if err != nil {
		g.writeStderr("\n")
		errorMsg := strings.Join(resourceKinds, "\n")
		return nil, fmt.Errorf("could not get api resources:\n%s", errorMsg)
	}

	sort.Strings(resourceKinds)
	totalResources := len(resourceKinds)
	g.writeStderr(" found %d.\n", totalResources)

	processedResources := 0
	var kubectlPipes []*script.Pipe
	clearConsoleLine := func() {
		g.writeStderr("\033[2K") // clear entire line
	}
	for _, resourceKind := range resourceKinds {
		kubectlGet := fmt.Sprintf("kubectl get --show-kind --ignore-not-found -o name %q", resourceKind)
		kubectlPipe := g.Exec(kubectlGet).Reject("Warning:")
		kubectlPipes = append(kubectlPipes, kubectlPipe)
		processedResources += 1
		if processedResources > 1 {
			clearConsoleLine()
			g.writeStderr("\033[F") // move cursor to start of previous line
		}
		g.writeStderr("Getting resources (%d/%d)\n", processedResources, totalResources)
		g.writeStderr("%s", resourceKind)
		time.Sleep(50 * time.Millisecond)
	}

	clearConsoleLine()

	processedResources = 0

	var allResourcesFound []string
	for _, kubectlPipe := range kubectlPipes {
		resourcesFound, err := kubectlPipe.Slice()
		if err != nil {
			return nil, fmt.Errorf("could not get resources: %w", getPipeError(kubectlPipe))
		}
		processedResources += 1
		allResourcesFound = append(allResourcesFound, resourcesFound...)
		g.writeStderr("\rWaiting for results (%d/%d)", processedResources, totalResources)
	}

	// Write a newline since we were not writing a newline for our progress message
	g.writeStderr("\n")

	sort.Strings(allResourcesFound)

	return allResourcesFound, nil
}
