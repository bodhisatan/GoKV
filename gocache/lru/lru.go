package lru

// double linked-list
import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access
type Cache struct {
	maxBytes   int64                    // 允许使用的最大内存
	nbytes     int64                    // 当前已使用的内存
	linkedList *list.List               // 双向链表
	cache      map[string]*list.Element // 键是字符串，值是双向链表中对应节点的指针
	// callback func, optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		linkedList: list.New(),
		cache:      make(map[string]*list.Element),
		OnEvicted:  onEvicted,
	}
}

// Get
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// move node to the end of the queue
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// Remove
func (c *Cache) RemoveOldest() {
	ele := c.linkedList.Back() // 获取最后一个元素，即最老的元素
	if ele != nil {
		c.linkedList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}

}

// Add/Modify
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.linkedList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.linkedList.Len()
}
