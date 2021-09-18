package other

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

var Db *sql.DB

type BaseInfo struct {
	Id        int
	IDCardNo  string
	RH        string
	BloodType string
	CustomerName string
}

func conn() {
	DB, _ := sql.Open("mysql", GetConf("mysql.User")+":"+GetConf("mysql.Pass")+"@tcp("+GetConf("mysql.Localhost")+":"+GetConf("mysql.Port")+")/"+GetConf("mysql.Database"))
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail", err)
		return
	}
	Db = DB
}

// 判断标题是否存在
func GetId(title string) (bool,int) {

	conn()

	type Check struct {
		Cid int
	}

	sql := "SELECT cid FROM `typecho_contents` WHERE title LIKE '%" + title + "%' limit 1"

	var ck Check

	rows, err := Db.Query(sql)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&ck.Cid)
	}
	rows.Close()

	if ck.Cid == 0 {
		return false,0
	} else {
		return true,ck.Cid
	}
}

// 判断标签是否存在
func GetMetaId(name string) (bool,int) {

	conn()

	type Check struct {
		Mid int
	}

	sql := "SELECT mid FROM `typecho_metas` WHERE name LIKE '%" + name + "%' limit 1"

	var ck Check

	rows, err := Db.Query(sql)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		rows.Scan(&ck.Mid)
	}
	rows.Close()

	if ck.Mid == 0 {
		return false, 0
	} else {
		return true, ck.Mid
	}
}



func StartIn(title, content, keyword, t, url  string) {
	conn()
	timea := strconv.Itoa(int(time.Now().Unix()))
	sql := fmt.Sprintf(
		"INSERT INTO `typecho_contents` (title, created, modified, text, type, status, authorId, allowComment, link) VALUES ('%s','%s','%s','%s','%s','%s', '%d', '%d')",
		title,
		t,
		timea,
		"<!--markdown-->" + content,
		"post",
		"publish",
		1,
		1,
		url,
	)
	_, err := Db.Exec(sql)
	if err != nil {
		log.Println("exec failed:", err, ", sql:", sql)
	}

	log.Println("文章：" + title + " 采集成功")

	Meta(keyword)
	ContentJoinMeta(title, keyword)
}

// 标签插入
func Meta(keyword string) {
	conn()
	a := strings.Split(keyword, ",")
	for _, v := range a{
		rs, _ := GetMetaId(v)
		if rs {
			continue
		}
		sql := "INSERT INTO `typecho_metas` (name,slug,type) VALUES ('" + v + "','" + v + "','tag')"
		stmt, err := Db.Prepare(sql)
		if err != nil {
			fmt.Println("Prepare fail:", err)
		}
		//设置参数以及执行sql语句
		res, err := stmt.Exec()
		if err != nil {
			log.Println("Exec fail:", err, res)
		}
	}
}

func ContentJoinMeta(name, keyword string) {
	conn()
	_, cid := GetId(name)

	a := strings.Split(keyword, ",")
	for _, v := range a{
		_, mid := GetMetaId(v)
		MetaAdd(mid)
		sql := fmt.Sprintf(
			"INSERT INTO `typecho_relationships` (mid, cid) VALUES ('%d','%d')",
			mid,
			cid,
		)
		_, err := Db.Exec(sql)
		if err != nil {
			log.Println("exec failed:", err, ", sql:", sql)
		}
	}
}

func MetaAdd(mid int) {
	conn()
	sql := fmt.Sprintf(
		"UPDATE `typecho_metas` set count=count+1 WHERE mid = %d", mid, )
	_, err := Db.Exec(sql)
	if err != nil {
		log.Println("exec failed:", err, ", sql:", sql)
	}

}