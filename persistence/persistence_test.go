package persistence

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPersistence_Current(t *testing.T) {
    t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
    Init()
    defer func() {
        Reset()
        Destruct()
    }()
    seq, err := Current()
    assert.NoError(t, err)
    assert.Equal(t, int64(0), seq)

    err = Save(10)
    assert.NoError(t, err)
    seq, err = Current()
    assert.NoError(t, err)
    assert.Equal(t, int64(10), seq)
}
