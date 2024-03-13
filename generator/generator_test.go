package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryGenerator_Generate(t *testing.T) {
	t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
	const N = 10000
	t.Run("test mongo generator", func(t *testing.T) {
		g, err := NewGenerator(false)
		assert.NoError(t, err)
		for i := 0; i < N; i++ {
			key, err := g.Generate()
			assert.NoError(t, err)
			assert.True(t, g.Existed(key))
		}
	})

	t.Run("test memory generator", func(t *testing.T) {
		g, err := NewGenerator(true)
		assert.NoError(t, err)
		for i := 0; i < N; i++ {
			key, err := g.Generate()
			assert.NoError(t, err)
			assert.True(t, g.Existed(key))
		}
	})
}
