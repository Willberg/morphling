// 计算输入的字符数
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		in := bufio.NewReader(f)
		for {
			r, n, err := in.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
				os.Exit(1)
			}
			if r == unicode.ReplacementChar && n == 1 {
				invalid++
				continue
			}
			counts[r]++
			utflen[n]++
		}
		fmt.Printf("rune\tcount\n")
		for c, n := range counts {
			fmt.Printf("%q\t%d\n", c, n)
		}
		fmt.Printf("len\tcount\n")
		for i, n := range utflen {
			if i > 0 {
				fmt.Printf("%d\t%d\n", i, n)
			}
		}
		if invalid > 0 {
			fmt.Printf("\n%d invalid utf-8 characters\n", invalid)
		}
	}
}
