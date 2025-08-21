package pokecache
import(
	"time"
	"sync"
)

type Cache struct{
	CacheMap	map[string]cacheEntry
	Mu			*sync.RWMutex
}

type cacheEntry struct{
	createdAt 	time.Time
	val			[]byte
}

func NewCache (interval time.Duration) (*Cache){
	cache := &Cache{}
	cache.CacheMap = make(map[string]cacheEntry)
	cache.Mu = &sync.RWMutex{}
	go cache.ReapLoop(interval)

	return cache
}

func (c *Cache) Add (key string, val []byte){
	c.Mu.Lock()
	defer c.Mu.Unlock()
	
	var entry cacheEntry
	entry.createdAt = time.Now()
	entry.val = val
	c.CacheMap[key] = entry
}

func (c *Cache) Get (key string) (val []byte, found bool){
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	if entry, ok := c.CacheMap[key]; ok{
		return entry.val, ok
	}
	return []byte{}, false
}

func (c *Cache) ReapLoop(interval time.Duration){
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		currentTime := <-ticker.C
		c.Mu.Lock()
		for key, val := range c.CacheMap{
			if currentTime.After(val.createdAt.Add(interval)){
				delete(c.CacheMap, key)
			}
		}
		c.Mu.Unlock()
	}
}
