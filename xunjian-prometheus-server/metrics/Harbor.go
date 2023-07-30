package metrics

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"xunjian-prometheus-server/tool"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type HarborMetric struct {
	TypeName         string
	IP               string
	Port             string
	Username         string
	Password         string
	Status           string
	Image            string
	ImageStatus      string
	CpuUsePercent    float64
	MemoryUsePercent float64
	DiskUsePercent   float64
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
type ComponentMetric struct {
	Components []Component `json:"components"`
	Status     string      `json:"status"`
}

func harborRequest(method, url, username, password string) ([]byte, int) {
	//	var method string = "GET"
	//	var url string = "https://192.168.0.10:30008/api/v2.0/ping"
	//	var username string = "admin"
	//	var password string = "Harbor12345"
	// 设置超时
	client := http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0
	}
	req.SetBasicAuth(username, password)
	response, err := client.Do(req)
	if err != nil {
		return nil, 0
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	//fmt.Println(response.StatusCode)
	return body, response.StatusCode
}

func initHarborList() []HarborMetric {
	//fmt.Println(tool.Conf.CheckList.Mysql)
	HarborList := make([]HarborMetric, 0)
	for _, Client := range tool.Conf.CheckList.Harbor {
		var tmp HarborMetric
		if Client["ip"] == "" || Client["port"] == "" || Client["username"] == "" || Client["password"] == "" || Client["image"] == "" {
			continue
		}
		tmp.TypeName = Client["type_name"]
		tmp.IP = Client["ip"]
		tmp.Port = Client["port"]
		tmp.Username = Client["username"]
		tmp.Password = Client["password"]
		tmp.Image = Client["image"]
		HarborList = append(HarborList, tmp)
	}
	return HarborList
}

func initHealthRe(client *HarborMetric, re []byte) {
	//	println(string(re))
	var tmp ComponentMetric
	err := json.Unmarshal(re, &tmp)
	// 如果反序列化失败，则结果不正确，认为失败
	if err != nil {
		client.Status = "Failed"
		return
	}
	// 获取组件健康结果
	if tmp.Status != "healthy" {
		client.Status = "Failed"
		return
	} else {
		client.Status = "Ok"
		return
	}
}

// 下载镜像
func downloadImage(image_cordinate, username, password string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, image_cordinate, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
	// io.Copy(os.Stdout, out)
}

//func pullImage(client *HarborMetric) {
//	// 获取repositorys
//	repositories_re, status_code := harborRequest("GET", "https://"+client.IP+":"+client.Port+"/api/v2.0/repositories", client.Username, client.Password)
//	if len(repositories_re) == 0 || status_code != 200 {
//		client.TestImage = "Failed"
//		client.TestImageStatus = "Failed"
//		return
//	}
//	//	fmt.Println(string(repositories_re))
//	var repositories []map[string]interface{}
//	err := json.Unmarshal(repositories_re, &repositories)
//	if err != nil {
//		client.TestImage = err.Error()
//		client.TestImageStatus = err.Error()
//		return
//	}
//	if len(repositories) == 0 {
//		client.TestImage = "NOT Found"
//		client.TestImageStatus = "Unkonw"
//		return
//	}
//	//	fmt.Println(repositories[0])
//	//project_id := repositories[0]["id"]
//	image_name := repositories[0]["name"].(string) //c2cloud/node-exporter
//
//	// 根据project+repositry 获取镜像的tags
//	image_tags, status_code := harborRequest("GET", "https://"+client.IP+":"+client.Port+"/api/v2.0/projects/"+strings.Split(image_name, "/")[0]+"/repositories/"+strings.Split(image_name, "/")[1]+"/artifacts", client.Username, client.Password)
//	if len(image_tags) == 0 || status_code != 200 {
//		client.TestImage = image_name
//		client.TestImageStatus = "Failed"
//		return
//	}
//	var tags []map[string]interface{}
//	err = json.Unmarshal(image_tags, &tags)
//	if err != nil {
//		client.TestImage = image_name
//		client.TestImageStatus = err.Error()
//		return
//	}
//	// 这里是只获取第一个tag
//	image_tag := tags[0]["tags"].([]interface{})[0].(map[string]interface{})["name"].(string)
//
//	image_cordinate := client.IP + ":" + client.Port + "/" + strings.Split(image_name, "/")[0] + "/" + strings.Split(image_name, "/")[1] + ":" + image_tag
//
//	err = downloadImage(image_cordinate, client.Username, client.Password)
//	if err != nil {
//		client.TestImage = image_cordinate
//		client.TestImageStatus = "download image Failed"
//		return
//	} else {
//		client.TestImage = image_cordinate
//		client.TestImageStatus = "Ok"
//	}
//}

func getHarborMetric(HarborList []HarborMetric) {
	nodeList := getnodelist()
	for index, harbor := range HarborList {
		for _, k := range nodeList {
			if k.IP == harbor.IP {
				HarborList[index].CpuUsePercent = k.CpuUsePercent
				HarborList[index].MemoryUsePercent = k.MemoryUsePercent
				HarborList[index].DiskUsePercent = k.DiskUsePercent
			}
		}
	}
}

func TestHarbor(wg *sync.WaitGroup, resultChan chan interface{}) {
	if tool.Conf.Controller["harbor"] != "true" {
		wg.Done()
		return
	}
	fmt.Println("--------------- 开始Harbor检查 ----------")
	HarborList := initHarborList()
	for i := 0; i < len(HarborList); i++ {
		// harbor组件状况检查
		health_re, status_code := harborRequest("GET", "https://"+HarborList[i].IP+":"+HarborList[i].Port+"/api/v2.0/health", HarborList[i].Username, HarborList[i].Password)
		if len(health_re) == 0 || status_code != 200 {
			HarborList[i].Status = "Failed"
			HarborList[i].Image = "N/A"
			HarborList[i].ImageStatus = "Failed"
		} else {
			initHealthRe(&HarborList[i], health_re)
			err := downloadImage(HarborList[i].Image, HarborList[i].Username, HarborList[i].Password)
			if err != nil {
				HarborList[i].ImageStatus = err.Error()
			} else {
				HarborList[i].ImageStatus = "Ok"
			}
		}
	}
	getHarborMetric(HarborList)
	fmt.Println("--------------- 结束Harbor检查 ----------")
	resultChan <- HarborList
	wg.Done()
}
