package yuacg

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
		"xz":      "%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"llt":     "%e6%b4%9b%e4%b8%bd%e5%a1%94lolita%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"ny":      "%e5%a5%b3%e4%bc%98%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"hf":      "%e6%b1%89%e6%9c%8d%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"zp":      "%e8%87%aa%e6%8b%8d%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"mz":      "%e6%bc%ab%e5%b1%95%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"sx":      "%e6%b0%b4%e4%b8%8b%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"xzxsp":   "%e5%86%99%e7%9c%9f%e5%b0%8f%e8%a7%86%e9%a2%91",
		"cosplay": "cosplay%e5%86%99%e7%9c%9f%e6%91%84%e5%bd%b1",
		"xzsp":    "%e5%86%99%e7%9c%9f%e8%a7%86%e9%a2%91",
	}

	for k, v := range list {
		go getList(k, v)
	}

}

// 获取列表
func getList(paths, dz string) {
	log.Println("雨溪萌域开始采集，当前分类：" + paths)
	len1 := 1
	for i := 1; i <= 10000; i++ {
		url := "http://www.yuacg.com/" + dz + "/page/" + strconv.Itoa(i) + "/"
		rs := other.HttpGet(url, 1)
		if strings.Contains(rs, "_404") {
			goto End
		} else {
			zz := `<a href="http://www.yuacg.com/(.*?)/" title="(.*?)" rel=".*?"><img class=".*?"`
			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			if len(submatch) < 1 {
				zz := `<a target="_blank" href="http://www.yuacg.com/(.*?)/" title="(.*?)" rel=".*?"><img class=".*?"`
				compile := regexp.MustCompile(zz)
				submatch = compile.FindAllSubmatch([]byte(rs), -1)
			}
			for _, matches := range submatch {
				row := detail(string(matches[1]), "", string(matches[2]), paths)
				//goto End
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
	log.Println("雨溪萌域采集完毕，当前分类：" + paths)
	return

}

func detail(id, img, title, paths string) bool {

	rs, _ := other.GetId(title)

	url := "http://www.yuacg.com/" + id + "/"
	sqlUrl := url
	//log.Println("雨溪萌域采集，当前地址：" + url)

	if rs {
		log.Println("雨溪萌域 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		// 第一次curl
		str, keyword, created_at := content(url, id, "", title, paths)

		if str != "" && keyword != "" && title != "" {
			other.StartIn(title, str, keyword, created_at, sqlUrl)
		} else {
			log.Println("雨溪萌域插入失败:" + url)
		}

		return false
	}

}

func content(url, id, page, title, paths string) (string, string, string) {
	rs := other.HttpGet(url, 1)
	//log.Println("雨溪萌域采集，当前地址：" + url)
	if strings.Contains(rs, "_404") {
		return "", "", ""
	}

	keywords := ""
	wss := `<meta name="keywords" content="(.*?)".*?>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1])
	}
	alt, str, img := "", "", ""
	zz := `<img src="(.*?)" alt="(.*?)" border=\".*?\" \/>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch {
		alt = string(matches[2])
		img = other.DownPic(id, "yucca/"+paths, page, title, string(matches[1]))
		str += "![" + alt + "](" + img + ")\n"
	}

	if str == "" {
		zz := `<img title=".*?" src="(.*?)" alt="(.*?)" />`
		compile := regexp.MustCompile(zz)
		submatch := compile.FindAllSubmatch([]byte(rs), -1)
		for _, matches := range submatch {
			alt := string(matches[2])
			img := other.DownPic(id, "yucca/"+paths, page, title, string(matches[1]))
			str += "![" + alt + "](" + img + ")\n"
		}
		if str == "" {
			zz := `<img src="(.*?)" alt="(.*?)" title=".*?" />`
			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				alt := string(matches[2])
				img := other.DownPic(id, "yucca/"+paths, page, title, string(matches[1]))
				str += "![" + alt + "](" + img + ")\n"
			}
		}

		if str == "" {
			zz := `<img src="(.*?)" alt="(.*?)" border=\".*?\">`
			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				alt = string(matches[2])
				img = other.DownPic(id, "yucca/"+paths, page, title, string(matches[1]))
				str += "![" + alt + "](" + img + ")\n"
			}
		}
	}

	t1 := ""
	zz1 := `<p class="data-label">最近更新</p><p class="info">(.*?)</p>`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	t1 = strings.Replace(t1, "年", "-", -1)
	t1 = strings.Replace(t1, "月", "-", -1)
	t1 = strings.Replace(t1, "日", "", -1)
	stamp, _ := time.ParseInLocation("2006-01-02", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)
	return str, keywords, created_at

}
