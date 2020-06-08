package gocache

import (
	"fmt"
	"log"
	"sync"
)

// Group是GoKV最核心的数据结构，复杂与用户交互，并且控制缓存值与获取的流程
/**
                是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/

// 一个Group是一个缓存的命名空间，每一个Group有一个唯一名称，例如缓存学生成绩命名为scores
type Group struct {
	name      string // 缓存的名称
	getter    Getter // 缓存未命中时获取数据的回调
	mainCache cache  // 并发缓存主体
}

// 返回key对应的缓存
type Getter interface {
	Get(key string) ([]byte, error)
}

// 实现Getter接口Get方法
type GetterFunc func(key string) ([]byte, error)

// 利用回调函数，当缓存不存在时，调用回调函数得到源数据
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 实例话Group，并将Group实例存储在groups里
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 如果当前缓存中有此数据，则直接返回，如果没有，通过回调函数加载数据并存入当前缓存中
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GoKV] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
