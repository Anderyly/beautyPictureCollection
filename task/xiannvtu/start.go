package xiannvtu

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"warm/other"
)

func Start() {

	list := map[string]string{
		"dx": "daxiong-27",
		"hs": "heisi-22",
		"ny": "nayi-293",
		"zf": "zhifu-34",
		"kb": "kunbang-123",
		"cf": "changfa-8",
	}

	for k, v := range list {
		getList(k, v)
	}

}

// 获取列表
func getList(paths, dz string) {
	log.Println("仙女图开始采集，当前分类：" + paths)
	len1 := 1
	for i := 1; i <= 10000; i++ {
		url := "https://www.xiannvtu.com/t/" + dz + "-" + strconv.Itoa(i) + ".html"
		rs := other.HttpGet(url, 2)
		if strings.Contains(rs, "首页") {
			goto End
		} else {
			zz := `<a target="_blank" href="/v/([0-9a-z]+)\.html" class="imageLink image"><img src="(.*?)" alt="(.*?)".*?</a>`

			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				row := detail(string(matches[1]), string(matches[2]), string(matches[3]), paths)
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
	log.Println("仙女图采集完毕，当前分类：" + paths)
	return

}

func detail(id, img, title, paths string) bool {
	rs, _ := other.GetId(title)

	url := "https://www.xiannvtu.com/v/" + id + ".html"
	sqlUrl := url

	if rs {
		log.Println("仙女图 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		// 第一次curl
		str := ""
		st, ids, keyword, created_at := content(url, id, "1", title, paths)
		str += st
		if ids != "" {
			for i := 2; i <= 1000; i++ {
				u := "https://www.xiannvtu.com/v/" + id + "_" + strconv.Itoa(i) + ".html"
				st, ids, _, _ := content(u, id, strconv.Itoa(i), title, paths)
				//fmt.Println("采集：" + u + "完毕")
				str += st
				if ids == "" {
					goto End
				}
			}
		}
	End:
		if str != "" && keyword != "" && title != "" {
			other.StartIn(title, str, keyword, created_at, sqlUrl)
		} else {
			log.Println("仙女图插入失败:" + url)
		}

		return false
	}

}

func content(url, id, page, title, paths string) (string, string, string, string) {
	rs := other.HttpGet(url, 2)
	//log.Println("仙女图采集，当前地址：" + url)
	if strings.Contains(rs, "404") {
		return "", "", "", ""
	}
	zz := `<a href='(.*?)'><img alt="(.*?)" src="(.*?)" .*?></a>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	str := ""
	ids := ""

	keywords := ""
	wss := `<meta name="keywords" content="(.*?)".*?>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1])
	}

	for _, matches := range submatch {
		ids = string(matches[1])
		alt := string(matches[2])
		img := other.DownPic(id, "xvt/"+paths, page, title, string(matches[3]))
		str += "![" + alt + "](" + img + ")\n"
	}

	t1 := ""
	zz1 := `\"pubDate\": \"(.*?)\",`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	t1 = strings.Replace(t1, "T", " ", -1)
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)
	return str, ids, keywords, created_at

}
