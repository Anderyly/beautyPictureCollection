package niutu114

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
		"ns":   "nvshen",
		"xqx":  "xiaoqingxin",
		"sx":   "suxiong",
		"tyjr": "tongyanjuru",
		"90h":  "90hou",
		"jp":   "jiepai",
		"yw":   "youwu",
		"yj":   "yujie",
		"br":   "baoru",
		"tnl":  "tunvlang",
		"ny":   "neiyiyouhuo",
		"qz":   "qizhimeinv",
		"jmh":  "jiemeihua",
		"qc":   "qingchunmeinv",
		"jp1":  "jipinmeinv",
		"mx":   "meixiong",
		"xxn":  "xiaoxiannv",
		"nm":   "nenmo",
		"swmt": "siwameitui",
		"jr":   "juru",
		"bjn":  "bijini",
		"zf":   "zhifuyouhuo",
		"qt":   "qiaotun",
		"ss":   "shishen",
	}

	for k, v := range list {
		go getList(k, v)
	}

}

// 获取列表
func getList(paths, dz string) {
	log.Println("牛图开始采集，当前分类：" + paths)
	len1 := 1
	for i := 1; i <= 10000; i++ {
		url := "http://www.niutu114.com/tag/" + dz + "/?page=" + strconv.Itoa(i)
		rs := other.HttpGet(url, 1)

		if !strings.Contains(rs, "am-fl") {
			goto End
		} else {
			zz := `<a href="http://www.niutu114.com/meinv/([0-9a-z]+)/([0-9a-z]+)/([0-9a-z]+).html" target="_blank">(.*?)</a>`

			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(rs), -1)
			for _, matches := range submatch {
				sub := string(matches[1])   // 分类
				date := string(matches[2])  // 日期
				id := string(matches[3])    // id
				title := string(matches[4]) // 标题
				row := detail(sub, date, id, other.TrimHtml(title))
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
	log.Println("牛图采集完毕，当前分类：" + paths)

	return

}

func detail(sub, date, id, title string) bool {
	rs, _ := other.GetId(title)

	url := "http://www.niutu114.com/meinv/" + sub + "/" + date + "/" + id
	sqlUrl := url
	//log.Println("牛图采集，当前地址：" + url)

	if rs {
		log.Println("牛图 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		// 第一次curl
		str := ""
		st, ids, keyword, created_at := content(url+".html", sub, date, "1", title, id)
		str += st
		if ids != false {
			for i := 2; i <= 10000; i++ {
				st, ids, _, _ := content(url+"_"+strconv.Itoa(i)+".html", sub, date, strconv.Itoa(i), title, id)
				str += st
				if ids == false {
					goto End
				}
			}
		}
	End:
		if str != "" && keyword != "" && title != "" {
			other.StartIn(other.TrimHtml(title), str, keyword, created_at, sqlUrl)
		} else {
			log.Println("牛图网插入失败:" + url)
		}

		return false
	}

}

func content(url, sub, date, page, title, id string) (string, bool, string, string) {
	//log.Println("牛图采集，当前地址：" + url)
	rs := other.HttpGet(url, 1)
	if strings.Contains(rs, "Page Not Found") {
		return "", false, "", ""
	}

	// 正则图片
	zz := `<a href="(.*?)"><img alt="(.*?)" src="(.*?)" /></a>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	str := ""
	for _, matches := range submatch {
		alt := string(matches[2])
		img := other.DownPic(id, "nt/"+sub+"/"+date, page, title, string(matches[3]))
		str += "![" + alt + "](" + img + ")\n"
	}

	// 正则分类
	keywords := ""
	wss := `<meta name="keywords" content="(.*?)".*?>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1])
	}

	keywords = strings.Replace(keywords, "，", ",", -1)

	t1 := ""
	zz1 := `<span class=\"title-time\">(.*?)<\/span>`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)

	return str, true, keywords, created_at

}
