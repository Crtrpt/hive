package api

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml/v2"
	"xorm.io/xorm"
)

type Config struct {
	Source string
	Redis  string
}

var engine *xorm.Engine

var ctx = context.Background()
var rdb *redis.Client

func Init() {

	content, err := ioutil.ReadFile("app.toml") // the file is inside the local directory
	if err != nil {
		fmt.Println("Err")
	}

	var cfg Config

	err = toml.Unmarshal(content, &cfg)

	if err != nil {
		fmt.Println("toml error")
	}

	engine, err = xorm.NewEngine("mysql", cfg.Source)

	if err != nil {
		fmt.Println("db error")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	result, err := rdb.Info(ctx, "Server").Result()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("redis server \r\n %s\r\n", result)

	results, err := engine.Query("select version() as v")

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("db version %s \r\n", results[0]["v"])

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
