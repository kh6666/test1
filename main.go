package main

import (
	"container/list"
	"fmt"
	"sync"
)

// entry 用于在双向链表中存储键值对，方便在淘汰时通过链表节点找到对应的 key
type entry struct {
	key   int
	value int
}

type LRUCache struct {
	mu       sync.Mutex
	capacity int
	cache    map[int]*list.Element // 哈希表：key -> 链表节点
	list     *list.List            // 双向链表：最新的在 Front，最久的在 Back
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

func (this *LRUCache) Get(key int) int {
	this.mu.Lock()
	defer this.mu.Unlock()
	// 1. 如果 key 不存在
	elem, exist := this.cache[key]
	if !exist {
		return -1
	}
	// 2. 如果 key 存在，将节点移动到链表头部（最新使用）
	this.list.MoveToFront(elem)
	return elem.Value.(*entry).value
}

func (this *LRUCache) Put(key int, value int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	// 1. 如果 key 已经存在，更新其 value 并移到头部
	if elem, exist := this.cache[key]; exist {
		this.list.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	// 2. 如果 key 不存在，创建新节点并插入到头部
	elem := this.list.PushFront(&entry{key: key, value: value})
	this.cache[key] = elem

	// 3. 如果超过容量，淘汰最久未使用的节点（链表尾部）
	if this.list.Len() > this.capacity {
		backElem := this.list.Back()
		if backElem != nil {
			this.list.Remove(backElem)
			kv := backElem.Value.(*entry)
			delete(this.cache, kv.key)
		}
	}
}

func main() {
	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Println(cache.Get(1)) // 返回 1

	cache.Put(3, 3)           // 该操作会导致密钥 2 作废
	fmt.Println(cache.Get(2)) // 返回 -1 (未找到)

	cache.Put(4, 4)           // 该操作会导致密钥 1 作废
	fmt.Println(cache.Get(1)) // 返回 -1 (未找到)
	fmt.Println(cache.Get(3)) // 返回 3
	fmt.Println(cache.Get(4)) // 返回 4
}
