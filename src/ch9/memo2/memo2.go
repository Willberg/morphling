// 通过监控goroutine实现并发非阻塞缓存
package memo2

// 用于记忆的函数类型
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// request 是请求消息，key 通过Func来调用
type request struct {
	key string
	// 客户端需要单个result
	response chan<- result
}

type Memo struct {
	requests chan request
}

type entry struct {
	res result
	// res 准备好会被关闭，起到全局通知的作用
	ready chan struct{}
}

// New返回f的函数记忆， 客户端之后需要调用Close
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// 对这个key 的第一次访问(有点问题，可能会有重复, 多个goroutine重复赋值)
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			// 调用f(key)
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// 执行函数
	e.res.value, e.res.err = f(key)
	//通知数据已经准备好
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// 等待该数据准备完毕
	<-e.ready
	// 向客户端发送结果
	response <- e.res
}
