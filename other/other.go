package other

import (
	"bufio"
	"errors"
	"github.com/aWildProgrammer/fconf"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func UnixToDate(str, dateVal string) (asd string) {

	strArr1T, err := strconv.ParseInt(str, 10, 64)
	CheckErr(err)
	t := time.Unix(strArr1T, 0)

	return t.Format("200601021504")

}

func CheckErr(err error) {
	if err != nil {
		log.Println("exec failed:", err)
	}
}

// http get请求
func HttpGet(url string, tp int) (res string) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	CheckErr(err)
	//request.Header.Set("Authorization", "APPCODE "+config.GetConf("token"))
	response, err := client.Do(request)
	CheckErr(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	CheckErr(err)
	if tp == 1 {
		return string(body)
	} else {
		return ConvertToString(string(body), "gbk", "utf-8")
	}
}

func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	CheckErr(err)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		CheckErr(err)
		for _, addr := range addrs {
			ip := GetIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func GetIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func DownPic(id, paths, page, title, imgUrl string) string {
	if GetConf("file.Path") == "1" {
		return imgUrl
	}

	imgPath := GetConf("file.Path") + paths + "/" + id + "/"
	err := CreateMutiDir(imgPath)
	CheckErr(err)
	f, err := os.Create(imgPath + title + ".txt")
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
	} else {
		_, err = f.Write([]byte(""))
		//CheckErr(err)
		if err != nil {
			log.Println(err)
			return ""
		}
	}

	fileName := path.Base(imgUrl)
	res, err := http.Get(imgUrl)
	if err != nil {
		return ""
		log.Println(err)
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	pv := ""
	if page == "" {
		pv = ""
	} else {
		pv = page + "_"
	}
	if !IsFileExist(imgPath + pv + fileName) {
		file, err := os.Create(imgPath + pv + fileName)
		if err != nil {
			log.Println("exec failed:", err)
		}
		// 获得文件的writer对象
		writer := bufio.NewWriter(file)

		io.Copy(writer, reader)
	}
	return GetConf("file.Url") + "=" + paths + "/" + id + "/" + pv + fileName
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

func CreateMutiDir(filePath string) error {
	if !IsFileExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			log.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}

func GetConf(name string) string {
	c, err := fconf.NewFileConf("./set.ini")
	if err != nil {
		log.Println(err)
	}
	return c.String(name)
}
