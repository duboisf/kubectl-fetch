package cmd_test

import (
	"testing"

	"github.com/duboisf/kubectl-fetch/internal/cmd"
	"github.com/duboisf/kubectl-fetch/internal/pkg/testing/assert"
)

func TestGetOptions(t *testing.T) {
	t.Parallel()
	t.Run("max in flight requests", func(t *testing.T) {
		opts, err := cmd.GetOptions([]string{"-p", "15"})
		assert.Nil(t, err)
		assert.Equals(t, 15, opts.MaxInFlight)
	})

	t.Run("only takes 1 optional argument", func(t *testing.T) {
		_, err := cmd.GetOptions([]string{"hi", "there"})
		assert.NotNil(t, err)
	})

	t.Run("returns an error if it cannot create a regex from the argument", func(t *testing.T) {
		_, err := cmd.GetOptions([]string{"(hi"})
		assert.NotNil(t, err)
	})
}
