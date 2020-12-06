// 统计文件
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileInfo struct {
	name string
	size int64
}

func walkDir(dir, typeName string, n *sync.WaitGroup, fileInfos chan<- FileInfo) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subDir := filepath.Join(dir, entry.Name())
			go walkDir(subDir, typeName, n, fileInfos)
		} else {
			if strings.Contains(entry.Name(), typeName) {
				f := FileInfo{
					name: entry.Name(),
					size: entry.Size(),
				}
				fileInfos <- f
			}
		}
	}
}

// sema 是一个用于限制目录并发数的计数计量量
var sema = make(chan struct{}, 20)

func dirents(dir string) []os.FileInfo {
	// 获取令牌
	sema <- struct{}{}
	defer func() {
		//释放令牌
		<-sema
	}()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}

func printDiskFileInfo(info FileInfo) {
	size := 0.0
	format := "%s\t%.3fB\n"
	if info.size < 1e3 {
		size = float64(info.size)
	} else if info.size < 1e6 {
		format = "%s\t%.3fKB\n"
		size = float64(info.size) / 1e3
	} else if info.size < 1e9 {
		format = "%s\t%.3fMB\n"
		size = float64(info.size) / 1e6
	} else {
		format = "%s\t%.3fGB\n"
		size = float64(info.size) / 1e9
	}

	fmt.Printf(format, info.name, size)
}

func main() {
	var typeName = flag.String("t", "mp4", "only search special type file")
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileInfos := make(chan FileInfo)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, *typeName, &n, fileInfos)
	}

	go func() {
		n.Wait()
		close(fileInfos)
	}()

	for fileInfo := range fileInfos {
		printDiskFileInfo(fileInfo)
	}
}
