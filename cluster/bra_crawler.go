//集群爬取，需要配置zookeeper
package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/bitly/go-simplejson"
	"github.com/jinzhu/gorm"
	crawler "github.com/nladuo/go-webcrawler"
	"github.com/nladuo/go-webcrawler/model"
	"github.com/qiniu/iconv"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Bra struct {
	ID           int    `sql:"AUTO_INCREMENT"`
	ItemId       string `sql:"unique"`
	SellerId     string
	Title        string
	CommentCount string
}

type BraRate struct {
	ID          int `sql:"AUTO_INCREMENT"`
	SizeInfo    string
	RateContent string
}

const (
	thread_num         int    = 2000
	GET_PROXY_URL             = "http://www.66ip.cn/getzh.php?getzh=mmpvmxywnwomuvw&getnum=10&isp=0&anonymoustype=4&start=&ports=&export=&ipaddress=&area=0&proxytype=0&api=https"
	PARSE_ITEM         string = "解析商品信息"
	PARSE_BRA_RATE     string = "解析商品评论信息"
	PARSE_BRA_RATE_NUM string = "解析商品评论数量"
)

var (
	mDb      *gorm.DB
	mCrawler *crawler.Crawler
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "lack parameter")
		os.Exit(-1)
	}
	config, err := model.GetDistributedConfigFromPath(os.Args[1])
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("mysql", "root:root@/taobao?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	mDb = &db
	//创建表用来储存文胸商品的基本信息
	if !db.HasTable(&Bra{}) {
		db.CreateTable(&Bra{})
	}
	//创建表用来储存文胸商品评论内容
	if !db.HasTable(&BraRate{}) {
		db.CreateTable(&BraRate{})
	}
	defer db.Close()
	//创建一个分布式爬虫，用sql数据库存储任务队列，这里使用mysql
	mCrawler = crawler.NewDistributedSqlCrawler(&db, config)
	addBaseTasks()
	//设置解析器
	itemParser := model.Parser{Identifier: PARSE_ITEM, Parse: ParseItem}
	braRateParser := model.Parser{Identifier: PARSE_BRA_RATE, Parse: ParseBraRate}
	braRateNumParser := model.Parser{Identifier: PARSE_BRA_RATE_NUM, Parse: ParseBraRateNum}
	//添加解析器
	mCrawler.AddParser(itemParser)
	mCrawler.AddParser(braRateParser)
	mCrawler.AddParser(braRateNumParser)
	//设置代理ip产生器
	mCrawler.SetProxyGenerator(NewMyProxyGenerator())
	//开始运行
	mCrawler.Run()
}

func addBaseTasks() {
	for i := 1; i <= 100; i++ {

		baseTask := model.Task{
			Identifier: PARSE_ITEM,
			Url:        "http://s.m.taobao.com/search?q=文胸&m=api4h5&page=" + strconv.FormatInt(int64(i), 10)}
		mCrawler.AddBaseTask(baseTask)
	}

}

type BraRateUserData struct {
	ItemId   string
	SellerId string
	PageNum  int
}

func ParseItem(res *model.Result, processor model.Processor) {
	fmt.Println(string(res.Response.Body))
	bras := parse_bras(res.Response.Body)
	if len(bras) == 0 {
		if checkItemAntiSpider(res.Response.Body) {
			task := *res.GetInitialTask()
			fmt.Println("被反爬虫了， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}
	for i := 0; i < len(bras); i++ {
		bra := &bras[i]
		mDb.Create(bra) // add bra to db
		task_url := "http://rate.tmall.com/list_detail_rate.htm?itemId=" + bra.ItemId + "&sellerId=" + bra.SellerId + "&currentPage=10000&pageSize=1000000"
		braByte, err := json.Marshal(bra)
		if err != nil {
			continue
		}
		task := model.Task{
			Identifier: PARSE_BRA_RATE_NUM,
			Url:        task_url,
			UserData:   braByte}
		processor.AddTask(task)
	}
}

func ParseBraRate(res *model.Result, processor model.Processor) {
	fmt.Println(string(res.Response.Body))
	bra_rates := parse_bra_rate(res.Response.Body)

	if len(bra_rates) == 0 {
		if checkItemRateAntiSpider(res.Response.Body) {
			task := *res.GetInitialTask()
			fmt.Println("被反爬虫了， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}

	for i := 0; i < len(bra_rates); i++ {
		bra_rate := &bra_rates[i]
		mDb.Create(bra_rate)
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

func checkItemRateAntiSpider(data []byte) bool {
	data_str := "{" + string(data) + "}"
	_, err := simplejson.NewJson([]byte(data_str))
	if err != nil {
		return true
	}
	return false
}

func ParseBraRateNum(res *model.Result, processor model.Processor) {
	fmt.Println(string(res.Response.Body))

	rate_num := parse_bra_rate_num(res.Response.Body)
	if rate_num == 0 {
		if checkItemRateAntiSpider(res.Response.Body) {
			task := *res.GetInitialTask()
			fmt.Println("被反爬虫了， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}

	var bra Bra
	err := json.Unmarshal(res.UserData, &bra)
	if err != nil {
		return
	}
	for i := 0; i <= rate_num; i++ {
		url := "http://rate.tmall.com/list_detail_rate.htm?itemId=" + bra.ItemId +
			"&sellerId=" + bra.SellerId + "&currentPage=" +
			strconv.FormatInt(int64(i), 10) + "&pageSize=1000000"
		task := model.Task{
			Identifier: PARSE_BRA_RATE,
			Url:        url}
		processor.AddTask(task)
	}
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

func parse_bra_rate_num(body []byte) int {
	data := "{" + string(body) + "}"
	js, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return 0
	}
	numFloat, err := js.Get("rateDetail").Get("paginator").Get("lastPage").Float64()
	if err != nil {
		return 0
	}
	return int(numFloat)
}

func parse_bra_rate(body []byte) (bra_rates []BraRate) {
	bra_rates = []BraRate{}

	data := "{" + string(body) + "}"
	cd, err := iconv.Open("utf-8", "gbk") // convert gbk to utf-8
	if err != nil {
		fmt.Println("iconv.Open failed!")
		return bra_rates
	}
	defer cd.Close()
	data = cd.ConvString(data)

	js, err := simplejson.NewJson([]byte(data))
	if err != nil {
		return
	}
	rates, err := js.Get("rateDetail").Get("rateList").Array()
	if err != nil {
		return
	}
	for i := range rates {
		rate := rates[i].(map[string]interface{})
		bra_rate := BraRate{
			SizeInfo:    rate["auctionSku"].(string),
			RateContent: rate["rateContent"].(string)}
		fmt.Println(bra_rate.SizeInfo, "   ", bra_rate.RateContent)
		bra_rates = append(bra_rates, bra_rate)
	}
	return
}

//代理ip生成器
type MyProxyGenerator struct {
	proxy_list [10]model.Proxy
	index      int
	used_times int
	lock       *sync.Mutex
}

func NewMyProxyGenerator() *MyProxyGenerator {
	var generator MyProxyGenerator
	generator.used_times = 0
	generator.lock = &sync.Mutex{}
	generator.index = 10
	return &generator
}

func (this *MyProxyGenerator) GetProxy() model.Proxy {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.used_times == 0 {
		this.index++
		if this.index >= 10 {
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
				this.proxy_list[k] = model.Proxy{IP: strs[0], Port: strs[1], Type: model.TYPE_HTTP}
				fmt.Println("ip:", strs[0], "port:", strs[1])
			}
			this.index = 0
		}
		this.used_times = 10 //一个代理ip用10次
	}
	proxy := this.proxy_list[this.index]
	//fmt.Println("Get IP:", proxy.IP, "  Port:", proxy.Port)
	this.used_times--
	return proxy
}
