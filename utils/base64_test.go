package utils

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPool_base64(t *testing.T) {
    assert.Equal(t, "a", Base64Encode(0))
    assert.Equal(t, "b", Base64Encode(1))
    assert.Equal(t, "-", Base64Encode(63))
    assert.Equal(t, "ba", Base64Encode(64))
    assert.Equal(t, "baaa", Base64Encode(64*64*64))
    assert.Equal(t, "---", Base64Encode(64*64*64-1))

    assert.Equal(t, int64(0), Base64Decode("a"))
    assert.Equal(t, int64(1), Base64Decode("b"))
    assert.Equal(t, int64(63), Base64Decode("-"))
    assert.Equal(t, int64(64), Base64Decode("ba"))
    assert.Equal(t, int64(64*64*64), Base64Decode("baaa"))
    assert.Equal(t, int64(64*64*64-1), Base64Decode("---"))
}
