//文胸商品条目爬虫,不使用代理爬取
package main

import (
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	crawler "github.com/nladuo/go-webcrawler"
	"github.com/nladuo/go-webcrawler/model"
	"strconv"
	"time"
)

type Bra struct {
	ID           int    `sql:"AUTO_INCREMENT"`
	ItemId       string `sql:"unique"`
	SellerId     string
	Title        string
	CommentCount string
}

const (
	thread_num int    = 100
	PARSE_ITEM string = "解析商品信息"
)

var (
	mDb      *gorm.DB
	mCrawler *crawler.Crawler
)

//mysql配置
const (
	DB_USER   = "root"
	DB_PASSWD = "root"
	DB_HOST   = "localhost"
	DB_PORT   = "3306"
	DBNAME    = "taobao"
)

func main() {
	db, err := gorm.Open("mysql", DB_USER+":"+DB_PASSWD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DBNAME+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}
	mDb = db
	//创建表用来储存文胸商品的基本信息
	if !db.HasTable(&Bra{}) {
		db.CreateTable(&Bra{})
	}
	//创建一个本地爬虫，用sql数据库存储任务队列，这里使用mysql
	mCrawler = crawler.NewLocalSqlCrawler(db, thread_num)
	addBaseTasks()
	//设置解析器
	itemParser := model.Parser{Identifier: PARSE_ITEM, Parse: ParseItem}
	//添加解析器
	mCrawler.AddParser(itemParser)
	//设置超时
	mCrawler.SetProxyTimeOut(10 * time.Second)
	//开始运行
	fmt.Println("Crawler Starting....")
	mCrawler.Run()
}

func addBaseTasks() {
	for i := 1; i <= 100; i++ {

		baseTask := model.Task{
			Identifier: PARSE_ITEM,
			Url:        "http://s.m.taobao.com/search?q=文胸&m=api4h5&page=" + strconv.FormatInt(int64(i), 10)}
		mCrawler.AddBaseTask(baseTask)
	}

	for i := 1; i <= 100; i++ {

		baseTask := model.Task{
			Identifier: PARSE_ITEM,
			Url:        "http://s.m.taobao.com/search?q=胸罩&m=api4h5&page=" + strconv.FormatInt(int64(i), 10)}
		mCrawler.AddBaseTask(baseTask)
	}

}

//解析文胸商品条目
func ParseItem(res model.Result, processor model.Processor) {
	fmt.Println(string(res.Response.Body))
	bras := parse_bras(res.Response.Body)
	if len(bras) == 0 {
		if checkItemAntiSpider(res.Response.Body) {
			//重新把task加入队列
			task := res.GetInitialTask()
			fmt.Println("被反爬虫了或者出现错误， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}
	for i := 0; i < len(bras); i++ {
		bra := &bras[i]
		mDb.Create(bra) // add bra to db
	}
}

//检查
func checkItemAntiSpider(body []byte) bool {
	tag := "url"
	js, err := simplejson.NewJson(body)
	if err != nil {
		return true
	}
	_, ok := js.CheckGet(tag)
	if ok {
		return true
	}
	return false
}

func parse_bras(body []byte) (bras []Bra) {

	bras = []Bra{}

	js, err := simplejson.NewJson(body)
	if err != nil {
		return
	}
	items, err := js.Get("listItem").Array()
	if err != nil {
		return
	}
	if len(items) == 0 {
		return
	}
	for i := range items {
		item := items[i].(map[string]interface{})
		if item["userId"] == nil {
			continue
		}
		bra := Bra{
			ItemId:       item["item_id"].(string),
			SellerId:     item["userId"].(string),
			Title:        item["title"].(string),
			CommentCount: item["commentCount"].(string)}
		fmt.Println(bra.ItemId, bra.Title)
		bras = append(bras, bra)
	}
	return
}
