package metrics

import (
	"database/sql"
	"fmt"
	"sync"
	"xunjian-prometheus-server/tool"

	_ "github.com/lib/pq"
)

type PgsMetric struct {
	TypeName          string
	IP                string
	Username          string
	Password          string
	Port              string
	Status            string
	CurrentConnection int
	MAXConnection     int
}

func initPgsList() []PgsMetric {
	//fmt.Println(tool.Conf.CheckList.Mysql)
	PgsList := make([]PgsMetric, 0)
	for _, Client := range tool.Conf.CheckList.Pgs {
		var tmp PgsMetric
		//fmt.Println(Client["ip"], Client["password"], Client["username"], Client["port"])
		if !tool.CheckIP(Client["ip"]) || Client["password"] == "" || Client["username"] == "" || Client["port"] == "" {
			continue
		}
		tmp.TypeName = Client["type_name"]
		tmp.IP = Client["ip"]
		tmp.Password = Client["password"]
		tmp.Username = Client["username"]
		tmp.Port = Client["port"]
		PgsList = append(PgsList, tmp)
	}
	//fmt.Println("1", PgsList)
	return PgsList
}

// 数据库连接测试
func connectPgs(client *PgsMetric) {
	// 如果你的表里有应用到datetime字段，记得要加上parseTime=True，不然解析不了这个类型
	//dataSource := "postgres://postgres:" + client.Password + "@" + client.IP + ":" + client.Port + "/postgres?sslmode=disable&timeout=1s"
	dataSource := "host=" + client.IP + " port=" + client.Port + " user=" + client.Username + " password=" + client.Password + " dbname=postgres sslmode=disable connect_timeout=2"
	db, _ := sql.Open("postgres", dataSource)
	if err := db.Ping(); err != nil {
		client.Status = "Failed"
		client.CurrentConnection = 0
		client.MAXConnection = 0
	} else {
		client.Status = "Ok"
		client.CurrentConnection = db.Stats().InUse
		client.MAXConnection = db.Stats().MaxOpenConnections
	}
	defer db.Close()
}
func TestPgs(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["pgs"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始Pgs检查 ----------")
	PgsList := initPgsList()
	for i := 0; i < len(PgsList); i++ {
		connectPgs(&PgsList[i])
	}
	fmt.Println("--------------- 结束Pgs检查 ----------")
	resultChan <- PgsList
	wg.Done()
}
