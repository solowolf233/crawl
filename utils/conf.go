package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-xorm/xorm"
	"gopkg.in/olivere/elastic.v3"
)

//Config define config
var Config config

type config struct {
	SpiderNumber     int
	SpiderListenPort int
	Database         string
	Log              *log.Logger
	Engine           *xorm.Engine
	ElasticClient    *elastic.Client
}

//Log return logger
func Log() *log.Logger {
	return Config.Log
}

//Init utilsl
func Init() {
	initLog()
	initConfig()
	initDatabase()
	initElastic()
}

func initElastic() {
	Config.ElasticClient.CreateIndex("torrent").Do()
}

func initConfig() {
	type config struct {
		Database         string `json:"database"`
		Elastic          string `json:"elastic"`
		SpiderNumber     int    `json:"spider_number"`
		SpiderListenPort int    `json:"spider_listen_port"`
	}

	f, err := os.Open("config/crawl.conf")
	exit(err)
	b, err := ioutil.ReadAll(f)
	exit(err)
	var c config
	err = json.Unmarshal(b, &c)
	exit(err)

	Config.SpiderNumber = c.SpiderNumber
	Config.Database = c.Database
	Config.SpiderListenPort = c.SpiderListenPort
	client, err := elastic.NewClient(elastic.SetURL(c.Elastic))
	exit(err)
	Config.ElasticClient = client
}

func initLog() {
	Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func initDatabase() {
	engine, err := xorm.NewEngine("mysql", Config.Database)
	if err != nil {
		panic(err)
	}
	engine.SetMaxIdleConns(1000)
	engine.SetMaxOpenConns(1000)
	Config.Engine = engine
}

func exit(err error) {
	if err != nil {
		Log().Fatalln(err)
	}
}
