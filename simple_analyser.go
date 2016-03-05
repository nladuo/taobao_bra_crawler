//商品购买信息简单分析
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"log"
	"regexp"
)

//mysql配置
const (
	DB_USER   = "root"
	DB_PASSWD = "root"
	DB_HOST   = "localhost"
	DB_PORT   = "3306"
	DBNAME    = "taobao"
)

func isKeyExist(key string, m map[string]int64) bool {
	for k, _ := range m {
		if k == key {
			return true
		}
	}
	return false
}

func main() {

	//根据文胸大小，A、B、C、D...分类
	basicBras := map[string]int64{}
	basicBras["whole"] = 0

	//根据文胸的具体大小，70A、70B、75B...分类
	detailBras := map[string]int64{}
	detailBras["whole"] = 0

	//根据文胸颜色分类
	colorBras := map[string]int64{}
	colorBras["whole"] = 0

	//使用正则表达式匹配70A、70B....
	regSizeDetail := regexp.MustCompile(`[5-9][0-9][A-K]{1}`)
	//使用正则表达式匹配A、B、C、D....
	regSize := regexp.MustCompile(`[A-P]`)

	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DBNAME+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//whether the item has add
	rows, err := db.Query("SELECT * FROM `bra_rates`")

	if err != nil {
		panic(err.Error())
	}
	var id int
	var sizeInfo string
	var rateContent string
	for rows.Next() {
		err := rows.Scan(&id, &sizeInfo, &rateContent)
		if err != nil {
			log.Fatal(err)
		}
		s1 := regSizeDetail.FindAllString(sizeInfo, -1)
		if len(s1) != 0 {
			if !isKeyExist(s1[0], detailBras) {
				detailBras[s1[0]] = 1
			} else {
				detailBras[s1[0]]++
			}
			s2 := regSize.FindAllString(s1[0], -1)

			if !isKeyExist(s2[0], basicBras) {
				basicBras[s2[0]] = 1
			} else {
				basicBras[s2[0]]++
			}

			basicBras["whole"]++
			detailBras["whole"]++
		}

	}

	colors := []string{"红色", "橙色", "黄色", "绿色", "蓝色", "紫色", "黑色", "白色", "粉色"}

	for i := 0; i < len(colors); i++ {
		rows, err = db.Query("SELECT * FROM `bra_rates` where size_info like '%" + colors[i] + "%'")
		if err != nil {
			panic(err.Error())
		}
		colorBras[colors[i]] = 0
		for rows.Next() {
			colorBras[colors[i]]++
			colorBras["whole"]++
		}
	}

	result := map[string](interface{}){}

	result["basic"] = basicBras
	result["detail"] = detailBras
	result["color"] = colorBras

	resStr, _ := json.Marshal(result)
	fmt.Println(string(resStr))

}
