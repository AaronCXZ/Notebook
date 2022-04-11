package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Hostbit(t *testing.T) {
	fmt.Println(hostbit)
}

func Test_LRU(t *testing.T) {
	lru := Constructor(2)

	lru.Put(1, 1)
	lru.Put(2, 2)
	assert.Equal(t, lru.Get(1), 1)
	lru.Put(3, 3)
	assert.Equal(t, lru.Get(2), -1)
	lru.Put(4, 4)
	assert.Equal(t, lru.Get(1), -1)
	assert.Equal(t, lru.Get(3), 3)
	assert.Equal(t, lru.Get(4), 4)
}

func Test_LRUPutGet(t *testing.T) {
	lru := Constructor(1)

	lru.Put(1, 2)
	assert.Equal(t, lru.Get(1), 2)
}

func Test_LRUPutGetPutGetGet(t *testing.T) {
	lru := Constructor(1)

	lru.Put(2, 1)
	assert.Equal(t, lru.Get(2), 1)
	lru.Put(3, 2)
	assert.Equal(t, lru.Get(2), -1)
	assert.Equal(t, lru.Get(3), 2)
}

func Test_LRUPPGPPG(t *testing.T) {
	lru := Constructor(2)

	lru.Put(2, 1)
	lru.Put(2, 2)
	assert.Equal(t, lru.Get(2), 2)
	lru.Put(1, 4)
	lru.Put(4, 1)
	assert.Equal(t, lru.Get(2), -1)
	assert.Equal(t, lru.Get(3), -1)
}

func Test_LRUPPGPPG2(t *testing.T) {
	lru := Constructor(2)

	lru.Put(2, 1)
	lru.Put(2, 2)
	assert.Equal(t, lru.Get(2), 2)
	lru.Put(1, 1)
	lru.Put(4, 1)
	assert.Equal(t, lru.Get(2), -1)
	assert.Equal(t, lru.Get(3), -1)

}
