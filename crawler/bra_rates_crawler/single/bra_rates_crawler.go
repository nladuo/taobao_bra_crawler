//文胸信息爬虫，单机版
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
	thread_num         int    = 500
	GET_PROXY_URL      string = "http://www.89ip.cn/api/?tqsl=10&cf=1"
	PARSE_BRA_RATE     string = "解析商品评论信息"
	PARSE_BRA_RATE_NUM string = "解析商品评论数量"
	PPROF_PORT                = "6060"
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
	//创建表用来储存文胸商品评论内容
	if !db.HasTable(&BraRate{}) {
		db.CreateTable(&BraRate{})
	}
	defer db.Close()
	//创建一个本地爬虫，用sql数据库存储任务队列，这里使用mysql
	mCrawler = crawler.NewLocalSqlCrawler(&db, thread_num)
	addBaseTasks()
	//设置解析器
	braRateParser := model.Parser{Identifier: PARSE_BRA_RATE, Parse: ParseBraRate}
	braRateNumParser := model.Parser{Identifier: PARSE_BRA_RATE_NUM, Parse: ParseBraRateNum}
	//添加解析器
	mCrawler.AddParser(braRateParser)
	mCrawler.AddParser(braRateNumParser)
	//设置代理生成器
	mProxyGenerator = NewMyProxyGenerator()
	mCrawler.SetProxyGenerator(mProxyGenerator)
	//使用net/pprof，查看状态
	mCrawler.SetPProfPort(PPROF_PORT)
	//开始运行
	mCrawler.Run()
}

func addBaseTasks() {
	var bras []Bra
	mDb.Find(&bras)
	for i := 0; i < len(bras); i++ {
		bra := &bras[i]
		task_url := "https://rate.tmall.com/list_detail_rate.htm?itemId=" + bra.ItemId + "&sellerId=" + bra.SellerId + "&currentPage=1&pageSize=1000000"
		braByte, err := json.Marshal(bra)
		if err != nil {
			continue
		}
		task := model.Task{
			Identifier: PARSE_BRA_RATE_NUM,
			Url:        task_url,
			UserData:   braByte}
		mCrawler.AddBaseTask(task)
	}
}

//解析文胸商品的评论
func ParseBraRate(res *model.Result, processor model.Processor) {
	fmt.Println("parse bra rate")
	//fmt.Println(string(res.Response.Body))
	bra_rates := parse_bra_rate(res.Response.Body)

	if len(bra_rates) == 0 {
		if checkItemRateAntiSpider(res.Response.Body) {
			//换代理
			mProxyGenerator.ChangeProxy(&res.UsedProxy)
			//重新把task加入队列
			task := *res.GetInitialTask()
			fmt.Println("被反爬虫了被反爬虫了或者出现错误， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}

	for i := 0; i < len(bra_rates); i++ {
		bra_rate := &bra_rates[i]
		mDb.Create(bra_rate)
	}
}

//解析文胸商品的页数
func ParseBraRateNum(res *model.Result, processor model.Processor) {
	//fmt.Println(string(res.Response.Body))
	fmt.Println("parse bra rate num")
	rate_num := parse_bra_rate_num(res.Response.Body)
	if rate_num == 0 {
		if checkItemRateAntiSpider(res.Response.Body) {
			task := *res.GetInitialTask()
			fmt.Println("被反爬虫了被反爬虫了或者出现错误， 重新加入task：", task.Url)
			processor.AddTask(task)
		}
		return
	}

	var bra Bra
	err := json.Unmarshal(res.UserData, &bra)
	if err != nil {
		return
	}
	//添加任务
	for i := 2; i <= rate_num; i++ {
		url := "https://rate.tmall.com/list_detail_rate.htm?itemId=" + bra.ItemId +
			"&sellerId=" + bra.SellerId + "&currentPage=" +
			strconv.FormatInt(int64(i), 10) + "&pageSize=1000000"
		task := model.Task{
			Identifier: PARSE_BRA_RATE,
			Url:        url}
		processor.AddTask(task)
	}
	//解析数据
	bra_rates := parse_bra_rate(res.Response.Body)
	for i := 0; i < len(bra_rates); i++ {
		bra_rate := &bra_rates[i]
		mDb.Create(bra_rate)
	}
}

func checkItemRateAntiSpider(data []byte) bool {
	data_str := "{" + string(data) + "}"
	_, err := simplejson.NewJson([]byte(data_str))
	if err != nil {
		return true
	}
	return false
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
		fmt.Println(err)
		fmt.Println("iconv.Open failed!")

		return bra_rates
	}
	data = cd.ConvString(data)
	//defer cd.Close() ////不使用defer语句，提升性能

	js, err := simplejson.NewJson([]byte(data))
	if err != nil {
		cd.Close()
		return
	}
	rates, err := js.Get("rateDetail").Get("rateList").Array()
	if err != nil {
		cd.Close()
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
	cd.Close()
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
	generator.updateProxyList()
	return &generator
}

func (this *MyProxyGenerator) ChangeProxy(proxy *model.Proxy) {
	this.lock.Lock()

	if this.index >= 10 {
		this.updateProxyList()
	}

	if this.proxy_list[this.index].IP == proxy.IP &&
		this.proxy_list[this.index].Port == proxy.Port {
		//change proxy
		this.index++
	}
	this.lock.Unlock()
}

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
	if len(ip_and_port_strs) != 10 {
		goto RETRY
	}
	for k, v := range ip_and_port_strs {
		strs := strings.Split(v, ":")
		this.proxy_list[k] = model.Proxy{IP: strs[0], Port: strs[1]}
		fmt.Println("Get proxy from network, ip:", strs[0], "port:", strs[1])
	}
	this.index = 0
}

func (this *MyProxyGenerator) GetProxy() model.Proxy {
	this.lock.Lock()
	if this.index >= 10 {
		this.updateProxyList()
	}

	proxy := this.proxy_list[this.index]
	fmt.Println("Get IP:", proxy.IP, "  Port:", proxy.Port)
	this.lock.Unlock()
	return proxy
}
