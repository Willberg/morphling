// 通过锁实现并发非阻塞缓存
package memo1

import (
	"sync"
)

// 用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res result
	// res 准备好会被关闭，起到全局通知的作用
	ready chan struct{}
}

type Memo struct {
	f Func
	// 保护cache
	mu    sync.Mutex
	cache map[string]*entry
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// 对key第一次访问， 这个goroutine 负责计算数据和广播数据
		//已准备完毕的信息
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		// 调用慢函数
		e.res.value, e.res.err = memo.f(key)

		// 广播数据已经准备好的消息
		close(e.ready)
	} else {
		// 对这个key的重复访问
		memo.mu.Unlock()

		// 等待数据完成
		<-e.ready
	}
	return e.res.value, e.res.err
}
