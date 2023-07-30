package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GenerateConfig() *restclient.Config {
	// 加载配置文件，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}

func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", value), 64)
	return value
}

func matchString(pattern string, value string) bool {
	if m, _ := regexp.MatchString(pattern, value); !m {
		return false
	}
	return true
}
