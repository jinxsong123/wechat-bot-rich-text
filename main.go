package main

import (
	"log"
	"time"
	"wechat-bot-rich-text/src/content"
	"wechat-bot-rich-text/src/filter"
	"wechat-bot-rich-text/src/sender"
)

const ()

func main() {
	tags := []string{"焦点", "市场", "公司", "央行", "A股", "宏观"}
	keys := []string{"中概", "北向", "ETF", "道指", "标普", "纳斯达克", "经济新闻", "日报", "财经", "澎湃", "Wind", "午间公告", "每经", "时报", "新浪"}
	groups := []string{"test"}
	friends := []string{"test"}
	filter := filter.NewFilter()
	filter.SetTag(tags)
	filter.SetKey(keys)

	content := content.NewContent(filter)

	wx, err := sender.NewWXsender(friends, groups)
	if err != nil {
		log.Fatalf("Failed to new wx sender, err: %s", err.Error())
	}
	for {
		contents := content.Content()
		for _, info := range contents {
			log.Printf("info = %s", *info)
		}
		wx.Send(contents)
		time.Sleep(15 * time.Second)
	}

}
