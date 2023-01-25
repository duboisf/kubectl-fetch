package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bitfield/script"
)

type KubectlExecution struct {
	pipe *script.Pipe
	resourceApiName string
}

func Main() error {
	fmt.Fprintf(os.Stderr, "Getting all kubernetes api resources...")
	resourceKinds, err := script.Exec("kubectl api-resources --verbs=list --namespaced -o name").
		RejectRegexp(regexp.MustCompile(`^events$`)).Slice()
	if err != nil {
		return fmt.Errorf("Could not get api resources: %w", err)
	}

	totalResources := len(resourceKinds)
	fmt.Fprintf(os.Stderr, " found %d resources.\n", totalResources)

	processedResources := 0

	var kubectlPipes []*script.Pipe
	for _, resourceKind := range resourceKinds {
		kubectlPipe := script.Exec("kubectl get --show-kind --ignore-not-found -o name " + resourceKind).Reject("Warning:")
		kubectlPipes = append(kubectlPipes, kubectlPipe)
		processedResources += 1
		fmt.Fprintf(os.Stderr, "\rGetting resources (%d/%d)", processedResources, totalResources)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Fprintln(os.Stderr, "")

	processedResources = 0

	var allResourcesFound []string
	for _, kubectlPipe := range kubectlPipes {
		resourcesFound, err := kubectlPipe.Slice()
		if err != nil {
			return fmt.Errorf("could not get resources: %s", strings.Join(resourcesFound, "\n") + err.Error())
		}
		processedResources += 1
		allResourcesFound = append(allResourcesFound, resourcesFound...)
		fmt.Fprintf(os.Stderr, "\rWaiting for results (%d/%d)", processedResources, totalResources)
	}

	if len(allResourcesFound) == 0 {
		fmt.Fprintf(os.Stderr, "\nNo resources found!\n")
	} else {
		fmt.Fprintf(os.Stdout, "\n%s\n", strings.Join(allResourcesFound, "\n"))
	}

	return nil
}

func main() {
	err := Main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s\n", err)
		os.Exit(1)
	}
}
