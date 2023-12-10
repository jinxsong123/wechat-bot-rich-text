//Get content info from sina

package content

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"wechat-bot-rich-text/src/filter"

	"github.com/PuerkitoBio/goquery"
)

const (
	cookie             string = "UOR=www.baidu.com,tech.sina.com.cn,; SINAGLOBAL=114.84.181.236_1579684610.152568; UM_distinctid=16fcc8a8b704c8-0a1d2def9ca4c6-33365a06-15f900-16fcc8a8b718f1; lxlrttp=1578733570; gr_user_id=2736e487-ee25-4d52-a1eb-c232ac3d58d6; grwng_uid=d762fe92-912b-4ea8-9a24-127a43143ebf; __gads=ID=d79f786106eb99a1:T=1582016329:S=ALNI_MZoErH_0nNZiM3D4E36pqMrbHHOZA; Apache=114.84.181.236_1582267433.457262; ULV=1582626620968:6:4:1:114.84.181.236_1582267433.457262:1582164462661; ZHIBO-SINA-COM-CN=; SUB=_2AkMpBPEzf8NxqwJRmfoWz2_ga4R2zQzEieKfWADoJRMyHRl-yD92qm05tRB6AoTf3EaJ7Bg2UU4l1CDZXUBCzEuJv3mP; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WhqhhGsPWdPjar0R99pFT8s"
	referer_url        string = "http://finance.sina.com.cn/7x24/?tag=0"
	base_url           string = "http://zhibo.sina.com.cn/api/zhibo/feed?callback=jQuery0&page=0%22+%22&page_size=20&zhibo_id=152&tag_id=0&dire=f&dpc=1&pagesize=20&_=0%20Request%20Method:GET%27"
	content_key_length int    = 3
)

type contentInterface interface {
	Content() []*string
}

type content struct {
	filter filter.FilterInterface
}

func NewContent(filter filter.FilterInterface) contentInterface {
	return &content{
		filter: filter,
	}
}

func (content content) get() ([][content_key_length]string, error) {
	req, _ := http.NewRequest("GET", base_url, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Host", "zhibo.sina.com.cn")
	req.Header.Add("Referer", referer_url)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to get content, err: %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	body_str := string(body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body_str[12 : len(body_str)-14]))
	if err != nil {
		log.Printf("Failed to parse HTML: %s", err.Error())
		return nil, err
	}
	bodyContent := doc.Find("body").Text()

	var data map[string]map[string]interface{}
	err = json.Unmarshal([]byte(bodyContent), &data)
	if err != nil {
		log.Printf("Failed to parse JSON: %s\n", err.Error())
		return nil, err
	}
	result := [][content_key_length]string{}
	if list, ok := data["result"]["data"].(map[string]interface{})["feed"].(map[string]interface{})["list"].([]interface{}); ok {

		for _, value := range list {
			tags := value.(map[string]interface{})["tag"].([]interface{})
			tag_string := "[ "
			for _, tag := range tags {
				tag_ := tag.(map[string]interface{})
				tag_string += tag_["name"].(string)
				tag_string += " "
			}
			tag_string += " ]"

			result = append(result, [content_key_length]string{tag_string, value.(map[string]interface{})["create_time"].(string), value.(map[string]interface{})["rich_text"].(string)})
		}
	}
	return result, nil

}

func (content content) Content() []*string {
	return content.filte()
}

func (content *content) filte() []*string {
	result := []*string{}
	infos, err := content.get()
	if err != nil {
		log.Printf("Failed to get content, err %s", err.Error())
	}
	for index := len(infos) - 1; index >= 0; index-- {
		info := infos[index]
		if len(info) == content_key_length {
			if content.filter.Compare(&info[0], &info[1], &info[2]) {
				message := fmt.Sprintf("%s %s %s", info[0], info[1], info[2])
				result = append(result, &message)
			}
		}
	}
	return result
}
