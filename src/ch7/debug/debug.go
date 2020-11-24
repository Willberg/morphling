//测试接口值
package main

import (
	"bytes"
	"io"
)

const debug = false

func main() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer)
	}

	f(buf)
}

func f(out io.Writer) {
	// out 的接口值的动态类型是*bytes.Buffer, out!=nil结果是true,但是动态值是nil，调用Write会出错（当debug=false时，out此时是包含空指针的非空接口）
	if out != nil {
		out.Write([]byte("done\n"))
	}
}
