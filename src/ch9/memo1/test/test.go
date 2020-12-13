// 测试
package main

import (
	"ch9/memo1"
	"ch9/memo2"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func incomingURLs() []string {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Print(err)
		return nil
	}

	urls := make([]string, 100)
	n := 0
	for _, link := range visit(nil, doc) {
		if n >= 100 {
			break
		}
		if strings.Contains(link, "https") {
			urls = append(urls, link)
			n++
		}
	}
	return urls
}

func memo1Test() {
	m := memo1.New(httpGetBody)
	var n sync.WaitGroup
	for _, url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			if value != nil {
				fmt.Printf("%s,%s,%d bytes\n", url, time.Since(start), len(value.([]byte)))
			}
			n.Done()
		}(url)
	}
	n.Wait()
}

func memo2Test() {
	m := memo2.New(httpGetBody)
	var n sync.WaitGroup
	for _, url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			if value != nil {
				fmt.Printf("%s,%s,%d bytes\n", url, time.Since(start), len(value.([]byte)))
			}
			n.Done()
		}(url)
	}
	n.Wait()
	m.Close()
}

func main() {
	//memo1Test()
	memo2Test()
}
