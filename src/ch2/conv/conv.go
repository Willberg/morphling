// 摄氏度与华氏度转换
package main

import (
	tempconv "ch2/tempconv"
	"fmt"
)

func main() {
	fmt.Println(tempconv.CToF(99))
	fmt.Println(tempconv.FToC(99))
}
