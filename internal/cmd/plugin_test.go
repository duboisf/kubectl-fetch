package cmd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/cmd"
	"github.com/duboisf/kubectl-fetch/internal/pkg/terminal"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

type mockUI struct {
	actualTotalKinds int
	updates          chan *terminal.GetResourcesUpdate
}

func (m *mockUI) SetTotalKinds(i int) chan<- *terminal.GetResourcesUpdate {
	m.actualTotalKinds = i
	return m.updates
}

type mockKubeClient struct {
	listApiResources struct {
		output []string
		err    error
	}
	getResources struct {
		output map[string][]string
		err    error
	}
}

func (m *mockKubeClient) ListApiResources(ctx context.Context, namespaced bool) ([]string, error) {
	return m.listApiResources.output, m.listApiResources.err
}

func (m *mockKubeClient) GetResources(ctx context.Context, kind string) ([]string, error) {
	return m.getResources.output[kind], m.getResources.err
}

func TestPlugin_Fetch(t *testing.T) {
	t.Parallel()

	t.Run("returns the list of all resources in a namespace", func(t *testing.T) {
		// When
		kubeClient := &mockKubeClient{}
		kubeClient.listApiResources.output = []string{
			// returns the list unsorted
			"service",
			"deployment",
			"configmap",
		}
		kubeClient.getResources.output = map[string][]string{
			"deployment": {"deployment/foo"},
			"service":    {"service/bar", "service/baz"},
		}
		opts, err := cmd.GetOptions([]string{"(deployment|service)"})
		assert.Nil(t, err)
		ui := &mockUI{}
		ui.updates = make(chan *terminal.GetResourcesUpdate, len(kubeClient.listApiResources.output))
		plugin, err := cmd.NewPlugin(kubeClient, opts, ui)
		assert.Nil(t, err)

		// When
		resources, err := plugin.Fetch(context.Background())

		// Then
		assert.Nil(t, err)
		expectedResources := []string{"deployment/foo", "service/bar", "service/baz"}
		assert.SliceEquals(t, expectedResources, resources)
	})

	t.Run("returns an error if there's an error getting the list of api resources", func(t *testing.T) {
		// Given
		kubeClient := &mockKubeClient{}
		kubeClient.listApiResources.err = fmt.Errorf("error getting api resources")
		opts, err := cmd.GetOptions([]string{"(deployment|service)"})
		assert.Nil(t, err)
		ui := &mockUI{}
		ui.updates = make(chan *terminal.GetResourcesUpdate, len(kubeClient.listApiResources.output))
		plugin, err := cmd.NewPlugin(kubeClient, opts, ui)
		assert.Nil(t, err)

		// When
		_, err = plugin.Fetch(context.Background())

		// Then
		assert.NotNil(t, err)
	})

	t.Run("returns an error if there's an error getting the list of resources for a kind", func(t *testing.T) {
		// Given
		kubeClient := &mockKubeClient{}
		kubeClient.listApiResources.output = []string{
			// returns the list unsorted
			"service",
			"deployment",
			"configmap",
		}
		kubeClient.getResources.err = fmt.Errorf("error getting resources")
		opts, err := cmd.GetOptions([]string{"(deployment|service)"})
		assert.Nil(t, err)
		ui := &mockUI{}
		ui.updates = make(chan *terminal.GetResourcesUpdate, len(kubeClient.listApiResources.output))
		plugin, err := cmd.NewPlugin(kubeClient, opts, ui)
		assert.Nil(t, err)

		// When
		_, err = plugin.Fetch(context.Background())

		// Then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error getting resources")
	})
}
