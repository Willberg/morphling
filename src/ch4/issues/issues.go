// 将符合搜索条件的issue格式化输出
package main

import (
	"ch4/github"
	"fmt"
	"log"
	"os"
)

func main() {
	PrintTable()
}

func PrintTable() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}
