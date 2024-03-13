package pool

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPool_Run(t *testing.T) {
	t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")
	pool := NewPool()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	pool.Run(ctx)

	// 1
	key := pool.Key()
	pool.Existed(key)
	assert.Equal(t, OneTimeGeneratedNumber-1, len(pool.cache))
	assert.Equal(t, int64(1), pool.produced)

	// 100
	for i := 0; i < OneTimeGeneratedNumber-1; i++ {
		pool.Key()
	}
	assert.Equal(t, 0, len(pool.cache))
	assert.Equal(t, int64(100), pool.produced)
	// 101
	pool.Key()
	assert.Equal(t, int64(101), pool.produced)

	assert.Equal(t, OneTimeGeneratedNumber-1, len(pool.cache))
}

func TestPool_Concurrent(t *testing.T) {
	t.Setenv("MONGODB_URI", "mongodb://127.0.0.1:27017/?retryWrites=false")

	pool := NewPool()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	pool.Run(ctx)

	var wg sync.WaitGroup
	const N = 150
	var lock sync.Mutex
	var cache = make(map[string]bool)
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer lock.Unlock()
			lock.Lock()
			cache[pool.Key()] = true
		}(i)
	}
	wg.Wait()
	assert.Equal(t, int64(N), pool.produced)
	pool.Key()
	assert.Equal(t, int64(N+1), pool.produced)

	assert.Equal(t, N, len(cache))
	for key := range cache {
		assert.True(t, pool.Existed(key))
	}

	time.Sleep(time.Second)
}
