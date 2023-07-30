package metrics

import (
	"fmt"
	"sync"
	"time"
	"xunjian-prometheus-server/tool"

	"github.com/go-redis/redis"
)

type RedisMetric struct {
	TypeName string
	IP       string
	Port     string
	Password string
	Status   string
}

func initRedisList() []RedisMetric {
	//fmt.Println(tool.Conf.CheckList.Mysql)
	RedisList := make([]RedisMetric, 0)
	for _, Client := range tool.Conf.CheckList.Redis {
		var tmp RedisMetric
		if !tool.CheckIP(Client["ip"]) || Client["port"] == "" {
			continue
		}
		tmp.TypeName = Client["type_name"]
		tmp.IP = Client["ip"]
		tmp.Port = Client["port"]
		tmp.Password = Client["password"]
		RedisList = append(RedisList, tmp)
	}
	return RedisList
}

func connectRedis(client *RedisMetric) {
	// password为空也没关系
	cli := redis.NewClient(&redis.Options{
		Addr:        client.IP + ":" + client.Port,
		Password:    client.Password,
		DB:          0,
		DialTimeout: 2 * time.Second,
	})

	_, err := cli.Ping().Result()
	if err != nil {
		client.Status = "Failed"
	} else {
		client.Status = "Ok"
	}
}

func TestRedis(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["redis"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始Redis检查 ----------")
	RedisList := initRedisList()
	for i := 0; i < len(RedisList); i++ {
		connectRedis(&RedisList[i])
	}
	fmt.Println("--------------- 结束Redis检查 ----------")
	resultChan <- RedisList
	wg.Done()
}
