package mm29

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
		"bizhi":   "bizhi",
		"xiezhen": "xiezhen",
	}

	for k, v := range list {
		go getList(k, v)
	}

}

// 获取列表
func getList(paths, dz string) {
	log.Println("mm29开始采集，当前分类：" + paths)
	len1 := 1
	for i := 1; i <= 10000; i++ {
		url := ""
		if i == 1 {
			url = "https://www.mm29.com/" + dz
		} else {
			url = "https://www.mm29.com/" + dz + "/list_" + strconv.Itoa(i) + ".html"
		}
		//log.Println(url)

		rs := other.HttpGet(url, 1)
		if strings.Contains(rs, "Not Found") {
			goto End
		} else {

			zz := `<a href="/` + dz + `/([0-9a-z]+)\.html" title="(.*?)"><img.*?</a>`

			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				row := detail(string(matches[1]), string(matches[2]), dz, paths)
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
	log.Println("mm29采集完毕，当前分类：" + paths)
	return

}

func detail(id, title, dz, paths string) bool {
	rs, _ := other.GetId(title)

	url := "https://www.mm29.com/" + dz + "/" + id + ".html"
	sqlUrl := url
	if rs {
		log.Println("mm29 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		// 第一次curl
		str := ""
		st, ids, keyword, created_at := content(url, id, "1", title, paths)
		str += st
		if ids != "" {
			for i := 2; i <= 1000; i++ {
				u := "https://www.mm29.com/" + dz + "/" + id + "_" + strconv.Itoa(i) + ".html"
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
			other.StartIn(title, str, strings.TrimRight(keyword, ","), created_at, sqlUrl)
		} else {
			log.Println("mm29插入失败:" + url)
		}
		return false
	}

}

func content(url, id, page, title, paths string) (string, string, string, string) {
	rs := other.HttpGet(url, 1)
	//log.Println("mm29采集，当前地址：" + url)
	if strings.Contains(rs, "Not Found") {
		return "", "", "", ""
	}

	keywords := ""
	wss := `class='badge badge-primary'>(.*?)</a>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1]) + ","
	}

	zz := `<a href='(.*?)'><img src='(.*?)' alt='(.*?)'></a>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	str := ""
	ids := ""
	for _, matches := range submatch {
		ids = string(matches[1])
		alt := other.RegexpF(string(matches[3]))
		img := other.DownPic(id, "mm29/"+paths, page, title, string(matches[2]))
		str += "![" + alt + "](" + img + ")\n"
	}

	t1 := ""
	zz1 := `<span class="time">时间：(.*?)</span>`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)
	return str, ids, keywords, created_at

}
