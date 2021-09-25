package snly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"warm/other"
)

type cjson struct {
	Success string `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		Arr struct {
			Action string `json:"action"`
			Type   string `json:"type"`
			Id     string `json:"id"`
			Paged  int    `json:"paged"`
		} `json:"arr"`
		Html   string `json:"html"`
		Starus int    `json:"starus"`
	} `json:"data"`
}

func Start() {

	id := []int{
		12, 13, 15,
	}

	for _, v := range id {
		go getList(v)
	}
}

// 获取列表
func getList(id int) {
	log.Println("少女领域开始采集")
	len1 := 1
	for i := 1; i <= 10000; i++ {
		res := post(id, i, 1)
		var rjson cjson
		json.Unmarshal([]byte(res), &rjson)
		if rjson.Success != "ok" {
			goto End
		} else {
			html := rjson.Data.Html

			zz := `<a class="meta-title" href="https://snly.vip/(.*?)/(.*?)/([0-9a-z]+)/">(.*?)</a>`

			compile := regexp.MustCompile(zz)
			submatch := compile.FindAllSubmatch([]byte(html), -1)
			for _, matches := range submatch {
				row := detail(string(matches[1]), string(matches[2]), string(matches[3]), string(matches[4]))
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
	log.Println("少女领域采集完毕")
	return

}

func detail(sub, class, id, title string) bool {
	rs, _ := other.GetId(title)

	url := "https://snly.vip/" + sub + "/" + class + "/" + id + "/"
	sqlUrl := url

	if rs {
		log.Println("少女领域 文章：" + title + " 已存在，跳过执行")
		return true
	} else {
		str, keyword, created_at := content(url, sub, class, id, title)
		if str != "" && keyword != "" && title != "" {
			other.StartIn(title, str, keyword, created_at, sqlUrl)
		} else {
			log.Println("少女领域插入失败:" + url)
		}

		return false
	}

}

func content(url, sub, class, id, title string) (string, string, string) {
	rs := other.HttpGet(url, 1)
	//log.Println("少女领域采集，当前地址：" + url)
	if strings.Contains(rs, "未找到页面") {
		return "", "", ""
	}

	keywords := ""
	wss := `<meta name="keywords" content="(.*?)".*?>`
	compile1 := regexp.MustCompile(wss)
	submatch1 := compile1.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch1 {
		keywords += string(matches[1])
	}

	zz := `<img src="(.*?)".*?class="aligncenter size-full.*?alt="(.*?)".*?>`
	compile := regexp.MustCompile(zz)
	submatch := compile.FindAllSubmatch([]byte(rs), -1)
	str := ""
	for _, matches := range submatch {
		alt := other.RegexpF(string(matches[2]))
		img := other.DownPic(id, "snly", class, title, string(matches[1]))
		str += "![" + alt + "](" + img + ")\n"
	}

	t1 := ""
	zz1 := `<span class="image-info-time"><i class="iconfont">.*?</i>(.*?)</span>`
	compile2 := regexp.MustCompile(zz1)
	submatch2 := compile2.FindAllSubmatch([]byte(rs), -1)
	for _, matches := range submatch2 {
		t1 = string(matches[1])
	}
	countSplit := strings.Split(t1, ".")
	if len(countSplit[1]) == 1 {
		countSplit[1] = "0" + countSplit[1]
	}
	if len(countSplit[2]) == 1 {
		countSplit[2] = "0" + countSplit[2]
	}
	t1 = strings.Replace(strings.Trim(fmt.Sprint(countSplit), "[]"), " ", "-", -1)
	stamp, _ := time.ParseInLocation("2006-01-02", t1, time.Local)
	created_at := strconv.FormatFloat(float64(stamp.Unix()), 'f', -1, 64)
	return str, keywords, created_at

}

func post(id, page, tp int) (res string) {
	url := "https://snly.vip/wp-admin/admin-ajax.php"
	param := "paged=" + strconv.Itoa(page) + "&action=postlist_newajax&id=" + strconv.Itoa(id) + "&type=cat"
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, strings.NewReader(param))
	other.CheckErr(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", url)
	response, err := client.Do(request)
	other.CheckErr(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	other.CheckErr(err)
	if tp == 1 {
		return string(body)
	} else {
		return other.ConvertToString(string(body), "gbk", "utf-8")
	}
}
