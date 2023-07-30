package metrics

import (
	"database/sql"
	"fmt"
	"sync"
	"xunjian-prometheus-server/tool"

	_ "github.com/go-sql-driver/mysql" //导入mysql包
)

type MySQLMetric struct {
	TypeName          string
	IP                string
	Username          string
	Password          string
	Port              string
	Status            string
	CurrentConnection int
	MAXConnection     int
}

func initMySQLList() []MySQLMetric {
	//fmt.Println(tool.Conf.CheckList.Mysql)
	MySQLList := make([]MySQLMetric, 0)
	for _, Client := range tool.Conf.CheckList.Mysql {
		var tmp MySQLMetric
		//fmt.Println(Client["ip"], Client["password"], Client["username"], Client["port"])
		if !tool.CheckIP(Client["ip"]) || Client["password"] == "" || Client["username"] == "" || Client["port"] == "" {
			continue
		}
		tmp.TypeName = Client["type_name"]
		tmp.IP = Client["ip"]
		tmp.Password = Client["password"]
		tmp.Username = Client["username"]
		tmp.Port = Client["port"]
		MySQLList = append(MySQLList, tmp)
	}
	//fmt.Println("1", MySQLList)
	return MySQLList
}

// 数据库连接测试
func connectMySQL(client *MySQLMetric) {
	// 如果你的表里有应用到datetime字段，记得要加上parseTime=True，不然解析不了这个类型
	dataSource := client.Username + ":" + client.Password + "@tcp(" + client.IP + ":" + client.Port + ")/mysql?charset=utf8&parseTime=True&timeout=1s"
	db, _ := sql.Open("mysql", dataSource)
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
func TestMySQL(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["mysql"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始MySQL检查 ----------")
	MySQLList := initMySQLList()
	for i := 0; i < len(MySQLList); i++ {
		connectMySQL(&MySQLList[i])
	}
	fmt.Println("--------------- 结束MySQL检查 ----------")
	resultChan <- MySQLList
	wg.Done()
}
