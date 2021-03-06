package gocache

// 缓存值的数据结构

// b存储真实的缓存值，选择byte类型是为了能支持任意数据类型，例如字符串、图片
type ByteView struct {
	b []byte
}

// 返回所占内存大小，在lru.Cache中，要求被缓存对象必须实现 Value 接口，即 Len() int 方法，返回其所占的内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// 将缓存值转化为字符串返回
func (v ByteView) String() string {
	return string(v.b)
}
