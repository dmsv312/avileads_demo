package utils

import (
	"sync"

	"github.com/beego/beego/v2/client/cache"
)

const memoryCacheConfig = `{"interval":360}`

var lock = &sync.Mutex{}

type memoryCache struct {
	cache.Cache
}

var singleMemoryCache *memoryCache

func CreateMemoryCache() *memoryCache {
	if singleMemoryCache == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleMemoryCache == nil {
			singleMemoryCache = &memoryCache{}
			singleMemoryCache.Cache, _ = cache.NewCache("memory", memoryCacheConfig)
		}
	}

	return singleMemoryCache
}
