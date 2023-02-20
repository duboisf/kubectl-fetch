package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"sync"

	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
)

// ProgressDisplayer is an interface for terminal.UI
type ProgressDisplayer interface {
	SetTotalKinds(int) chan<- *terminal.GetResourcesUpdate
}

// KubeClient is an interface for kubectl.Kubectl
type KubeClient interface {
	ListApiResources(ctx context.Context, namespaced bool) ([]string, error)
	GetResources(ctx context.Context, kind string) ([]string, error)
}

type Plugin struct {
	kubeClient KubeClient
	options    *Options
	ui        ProgressDisplayer
}

type fetchResult struct {
	results []string
	err     error
}

type getResourcesResult struct {
	kind      string
	resources []string
	err       error
}

// NewPlugin returns a new Plugin ready to be used
func NewPlugin(kubeClient KubeClient, options *Options, tui ProgressDisplayer) (*Plugin, error) {
	return &Plugin{
		kubeClient: kubeClient,
		options:    options,
		ui:        tui,
	}, nil
}

func (p *Plugin) Fetch(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := sync.WaitGroup{}
	kinds, err := p.kubeClient.ListApiResources(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("could not get namespaced API resources:\n%w", err)
	}
	if p.options.Pattern != nil {
		kinds = filterKinds(kinds, p.options.Pattern)
	}
	totalKinds := len(kinds)
	getResourcesUpdates := p.ui.SetTotalKinds(totalKinds)

	maxParallel := make(chan struct{}, p.options.MaxInFlight)
	getResourcesResults := make(chan *getResourcesResult, totalKinds)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, kind := range kinds {
			maxParallel <- struct{}{}
			kind := kind
			wg.Add(1)
			go func() {
				defer wg.Done()
				resources, err := p.kubeClient.GetResources(ctx, kind)
				getResourcesResults <- &getResourcesResult{
					kind:      kind,
					resources: resources,
					err:       err,
				}
				getResourcesUpdates <- &terminal.GetResourcesUpdate{
					Kind:      kind,
					Resources: len(resources),
				}
			}()
		}
	}()

	go func() {
		defer close(getResourcesResults)
		wg.Wait()
	}()

	var allResources []string
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case results, more := <-getResourcesResults:
			if !more {
				sort.Strings(allResources)
				close(getResourcesUpdates)
				return allResources, nil
			}
			if results.err != nil {
				cancel()
				return nil, results.err
			}
			allResources = append(allResources, results.resources...)
			<-maxParallel
		}
	}
}

func filterKinds(kinds []string, pattern *regexp.Regexp) []string {
	var filtered []string
	for _, kind := range kinds {
		if pattern.MatchString(kind) {
			filtered = append(filtered, kind)
		}
	}
	return filtered
}
