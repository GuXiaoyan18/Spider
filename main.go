package main

import (
	"spiderTool/spider"
	"strconv"
	"time"
)

func main() {
	ch := make(chan string, 10)
	defer close(ch)

	go spider.SoundcardGet(ch)
	for i := 1; i < 3; i++ {
		time.Sleep(2 * time.Second)
		url := "http://detail.zol.com.cn/sound_card/" + strconv.Itoa(i) + ".html"
		spider.ListGet(url, ch)
	}
	ch <- "NULL"
}
