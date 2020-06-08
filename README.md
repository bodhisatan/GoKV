# GoKV
a distributed kv cache like groupcache written by golang

## 项目结构
```
GokV/
  |--lru/
      |--lru.go  // lru 缓存淘汰策略
  |--byteview.go // 缓存值的抽象与封装
  |--cache.go    // 并发控制
  |--gocache.go // 负责与外部交互，控制缓存存储和获取的主流程
  |--http.go     // 提供被其他节点访问的能力(基于http)
```
