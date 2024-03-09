package sequencer

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPersistence_Current(t *testing.T) {
    t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")

    run := func(sequencer Sequencer) {
        seq, err := sequencer.Current()
        assert.NoError(t, err)
        assert.Equal(t, int64(0), seq)

        err = sequencer.Save(10)
        assert.NoError(t, err)

        seq, err = sequencer.Current()
        assert.NoError(t, err)
        assert.Equal(t, int64(10), seq)

        sequencer.Reset()
    }

    t.Run("test mongo sequencer", func(t *testing.T) {
        sequencer, err := NewMongoSequencer()
        assert.NoError(t, err)
        run(sequencer)
    })

    t.Run("test memory sequencer", func(t *testing.T) {
        sequencer, err := NewMemorySequencer()
        assert.NoError(t, err)
        run(sequencer)
    })
}
