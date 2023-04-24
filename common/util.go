package common

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"google.golang.org/protobuf/proto"
)

func GenerateID() string {
	id := uuid.New()
	return id.String()
}

func First(items ...string) string {
	for _, s := range items {
		if s != "" {
			return s
		}
	}
	return ""
}

func FirstInt32(items ...int32) int32 {
	for _, s := range items {
		if s != -999 {
			return s
		}
	}
	return -999
}

func NewCache[T any]() Cache[T] {
	return &CacheWrapper[T]{
		c: cache.New(time.Minute*15, 30*time.Minute),
	}
}

type Cache[T any] interface {
	Get(key string) (T, bool)
	Set(key string, value T)
	IsEmpty() bool
	All() []T
}

type CacheWrapper[T any] struct {
	c *cache.Cache
}

func (c *CacheWrapper[T]) Get(key string) (T, bool) {
	item, found := c.c.Get(key)
	if found {
		return item.(T), found
	}
	var zero T
	return zero, false
}

func (c *CacheWrapper[T]) Set(key string, value T) {
	c.c.SetDefault(key, value)
}

func (c *CacheWrapper[T]) IsEmpty() bool {
	return c.c.ItemCount() == 0
}

func (c *CacheWrapper[T]) All() []T {
	items := c.c.Items()
	var rtn []T
	for _, k := range items {
		if !k.Expired() {
			item := k.Object.(T)
			rtn = append(rtn, item)
		}
	}
	return rtn
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func Clone[T proto.Message](obj T) T {
	return proto.Clone(obj).(T)
}

func Random64() uint64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0
	}
	return binary.LittleEndian.Uint64(b[:])
}

func Random1(n int32) int32 {
	num := Random64()
	val2 := (float64(num) / float64(math.MaxUint64) * float64(n)) + 1
	if val2 > float64(n) {
		return n
	}
	val3 := math.Round(val2)
	return int32(val3)
}
