package filter

import (
	"log"
	"strings"
	"time"
)

const (
	timeLayout string = "2006-01-02 15:04:05"
)

type FilterInterface interface {
	Compare(*string, *string, *string) bool
	SetTag(tags []string)
	SetKey(keys []string)
}

type Filter struct {
	lastTime    int64
	tags        []string
	contentKeys []string
}

func NewFilter() FilterInterface {
	return &Filter{
		lastTime: time.Now().Unix(),
	}
}

func (filter *Filter) Compare(tags, createTime, content *string) bool {
	if filter.compareTime(*createTime) && (filter.compareContent(content) || filter.compareTag(*tags)) {
		return true
	}
	return false
}
func (filter Filter) compareContent(content *string) bool {
	for _, key := range filter.contentKeys {
		if strings.Contains(*content, key) {
			return true
		}
	}
	return false
}

func (filter Filter) compareTag(tags string) bool {
	for _, tag := range filter.tags {
		if strings.Contains(tags, tag) {
			return true
		}
	}
	return false
}

func (filter *Filter) compareTime(createTime string) bool {
	parsedTime, err := time.Parse(timeLayout, createTime)
	if err != nil {
		log.Printf("Error parsing time: %s", err.Error())
		return false
	}

	if filter.lastTime < parsedTime.Unix() {
		filter.lastTime = parsedTime.Unix()
		return true
	}
	return false
}

func (filter *Filter) SetTag(tags []string) {

	filter.tags = tags
}

func (filter *Filter) SetKey(keys []string) {
	filter.contentKeys = keys
}
