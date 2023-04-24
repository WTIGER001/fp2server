package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom1(t *testing.T) {
	ten := int32(10)
	twenty := int32(20)
	two := int32(2)
	hundred := int32(100)
	one := int32(1)

	rolls10 := make(map[int32]int32)
	for i := 0; i < 1000; i++ {
		v := Random1(ten)
		assert.GreaterOrEqual(t, ten, v)
		assert.LessOrEqual(t, one, v)
		rolls10[v] = rolls10[v] + 1
	}
	assert.Equal(t, len(rolls10), 10)

	rolls20 := make(map[int32]int32)
	for i := 0; i < 1000; i++ {
		v := Random1(twenty)
		assert.GreaterOrEqual(t, twenty, v)
		assert.LessOrEqual(t, one, v)
		rolls20[v] = rolls20[v] + 1
	}
	assert.Equal(t, len(rolls20), 20)

	for i := 0; i < 1000; i++ {
		v := Random1(two)
		assert.GreaterOrEqual(t, two, v)
		assert.LessOrEqual(t, one, v)
	}

	rolls100 := make(map[int32]int32)
	for i := 0; i < 1000; i++ {
		v := Random1(hundred)
		assert.GreaterOrEqual(t, hundred, v)
		assert.LessOrEqual(t, one, v)
		rolls100[v] = rolls100[v] + 1
	}
	assert.Equal(t, len(rolls100), 100)
}
