package main

import (
	"spiderTool/common"
	"spiderTool/spider"
	"strconv"
)

func main() {
	ch := make(chan string, 10)
	cpuCh := make(chan common.CPU, 10)
	defer close(ch)

	go spider.ItemGet(ch, cpuCh)
	for i := 1; i < 4; i++ {
		// time.Sleep(2 * time.Second)
		url := "http://detail.zol.com.cn/cpu/" + strconv.Itoa(i) + ".html"
		spider.ListGet(url, ch)
	}
	ch <- "NULL"
}
