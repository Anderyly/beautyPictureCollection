package shejumm

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"warm/other"
)

func Start() {
	getList()
}

// 获取列表
func getList() {
	log.Println("射菊MM开始采集")
	len1 := 1
	for i := 1; i <= 10000; i++ {
		url := "https://www.shejumm.com/?page=" + strconv.Itoa(i)
		rs := other.HttpGet(url, 1)
		if !strings.Contains(rs, "?page") {
			goto End
		} else {

			zz := `<div class="item-title"><a href="/article/([0-9a-z]+)/">(.*?)</a></div>`

			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				row := detail(string(matches[1]), string(matches[2]))
				if row {
					len1 += 1
					if len1 == 3 {
						goto End
					}
					break
				}
			}
		}
	}
End:
	log.Println("射菊MM采集完毕")
	return

}

func detail(id, title string) bool {
	rs, _ := other.GetId(title)

	url := "https://www.shejumm.com/article/" + id + "/"
	sqlUrl := url

	if rs {
		log.Println("射菊MM 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		str, keyword, created_at := content(url, id, title)
		if str != "" && keyword != "" && title != "" {
			other.StartIn(title, str, keyword, created_at, sqlUrl)
		} else {
			log.Println("射菊MM插入失败:" + url)
		}

		return false
	}

}

func content(url, id, title string) (string, string, string) {
	rs := other.HttpGet(url, 1)
	//log.Println("仙女图采集，当前地址：" + url)
	if strings.Contains(rs, "DoesNotExist") {
		return "", "", ""
	}

	keywords := ""
	wss := `<meta name="keywords" content="(.*?)".*?>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1])
	}

	zz := `<img width="100%" height="100%" class="lazy mz-mg-100" src=".*?" data-original="(.*?)" alt="(.*?)"/>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	str := ""
	for _, matches := range submatch {
		alt := string(matches[2])
		img := other.DownPic(id, "sjmm", "", title, string(matches[1]))
		str += "![" + alt + "](" + img + ")\n"
	}

	t1 := ""
	zz1 := `<time>(.*?)</time>`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	stamp, _ := time.ParseInLocation("2006-01-02", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)
	return str, keywords, created_at

}
