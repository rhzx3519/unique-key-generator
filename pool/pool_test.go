package pool

import (
    "context"
    "github.com/stretchr/testify/assert"
    "rhzx3519/unique-key-generator/persistence"
    "sync"
    "testing"
)

func TestPool_GenerateKey(t *testing.T) {
    t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
    persistence.Init()
    defer func() {
        persistence.Reset()
        persistence.Destruct()
    }()

    pool := NewPool()
    ctx, cancel := context.WithCancel(context.TODO())
    defer cancel()

    pool.Run(ctx)
    for i := 0; i < 64*64; i++ {
        assert.Equal(t, base64Encode(int64(i)), pool.Key())
    }
}

func TestPool_base64(t *testing.T) {
    assert.Equal(t, "a", base64Encode(0))
    assert.Equal(t, "b", base64Encode(1))
    assert.Equal(t, "-", base64Encode(63))
    assert.Equal(t, "ba", base64Encode(64))
    assert.Equal(t, "baaa", base64Encode(64*64*64))
    assert.Equal(t, "---", base64Encode(64*64*64-1))

    assert.Equal(t, int64(0), base64Decode("a"))
    assert.Equal(t, int64(1), base64Decode("b"))
    assert.Equal(t, int64(63), base64Decode("-"))
    assert.Equal(t, int64(64), base64Decode("ba"))
    assert.Equal(t, int64(64*64*64), base64Decode("baaa"))
    assert.Equal(t, int64(64*64*64-1), base64Decode("---"))
}

func TestPool_Concurrent(t *testing.T) {
    t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
    persistence.Init()
    defer func() {
        persistence.Reset()
        persistence.Destruct()
    }()

    pool := NewPool()
    ctx, cancel := context.WithCancel(context.TODO())
    defer cancel()

    pool.Run(ctx)

    var wg sync.WaitGroup
    const N = 64*64 + 33
    for i := 0; i < N; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            pool.Key()
        }(i)
    }
    wg.Wait()
    assert.Equal(t, base64Encode(int64(N)), pool.Key())
}

func TestPool_IsExist(t *testing.T) {
    t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
    persistence.Init()
    defer func() {
        persistence.Reset()
        persistence.Destruct()
    }()

    pool := NewPool()
    ctx, cancel := context.WithCancel(context.TODO())
    defer cancel()

    pool.Run(ctx)

    var wg sync.WaitGroup
    const N = 64 * 64
    for i := 0; i < N; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            pool.Key()
        }(i)
    }
    wg.Wait()

    assert.True(t, pool.IsExist(base64Encode(64*64)))
    assert.False(t, pool.IsExist(base64Encode(64*64+100)))
}
