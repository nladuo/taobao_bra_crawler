//文胸商品条目爬虫
package main

import (
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	crawler "github.com/nladuo/go-webcrawler"
	"github.com/nladuo/go-webcrawler/model"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	thread_num    int    = 100
	GET_PROXY_URL        = "http://www.89ip.cn/api/?tqsl=10&cf=1" //一个免费的代理服务器
	PARSE_ITEM    string = "解析商品信息"
)

var (
	mDb             *gorm.DB
	mCrawler        *crawler.Crawler
	mProxyGenerator *MyProxyGenerator
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
	mDb = &db
	//创建表用来储存文胸商品的基本信息
	if !db.HasTable(&Bra{}) {
		db.CreateTable(&Bra{})
	}
	//创建一个本地爬虫，用sql数据库存储任务队列，这里使用mysql
	mCrawler = crawler.NewLocalSqlCrawler(&db, thread_num)
	addBaseTasks()
	//设置解析器
	itemParser := model.Parser{Identifier: PARSE_ITEM, Parse: ParseItem}
	//添加解析器
	mCrawler.AddParser(itemParser)
	//设置代理生成器
	mProxyGenerator = NewMyProxyGenerator()
	mCrawler.SetProxyGenerator(mProxyGenerator)
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
			//换代理
			mProxyGenerator.ChangeProxy(res.UsedProxy)
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

//代理ip生成器
type MyProxyGenerator struct {
	proxy_list [10]model.Proxy
	index      int
	lock       *sync.Mutex
}

func NewMyProxyGenerator() *MyProxyGenerator {
	var generator MyProxyGenerator
	generator.lock = &sync.Mutex{}
	generator.index = 10
	return &generator
}

func (this *MyProxyGenerator) ChangeProxy(proxy model.Proxy) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.index >= 10 {
		this.updateProxyList()
	}

	if this.proxy_list[this.index].IP == proxy.IP &&
		this.proxy_list[this.index].Port == proxy.Port {
		//change proxy
		this.index++
	}
}

//网络请求，获取代理ip
// 利用正则表达式匹配 ip:port 格式
func (this *MyProxyGenerator) updateProxyList() {
RETRY:
	resp, err := http.Get(GET_PROXY_URL)
	if err != nil {
		goto RETRY
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto RETRY
	}
	reg_ip_and_port := regexp.MustCompile(`[0-9]*\.[0-9]*\.[0-9]*\.[0-9]*\:[0-9]*`)
	ip_and_port_strs := reg_ip_and_port.FindAllString(string(data), -1)
	for k, v := range ip_and_port_strs {
		strs := strings.Split(v, ":")
		this.proxy_list[k] = model.Proxy{IP: strs[0], Port: strs[1]}
		fmt.Println("Get proxy from network, ip:", strs[0], "port:", strs[1])
	}
	this.index = 0
}

func (this *MyProxyGenerator) GetProxy() model.Proxy {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.index >= 10 {
		this.updateProxyList()
	}

	proxy := this.proxy_list[this.index]
	fmt.Println("Get IP:", proxy.IP, "  Port:", proxy.Port)
	return proxy
}
