package tool

import (
	"log"
	"os"
	"regexp"

	"net/url"

	yaml "gopkg.in/yaml.v3"
)

// 配置文件
var Conf Config

type checklist struct {
	Business []map[string]string `yaml:"businness"`
	Harbor   []map[string]string `yaml:"harbor"`
	Nfs      []map[string]string `yaml:"nfs"`
	Mysql    []map[string]string `yaml:"mysql"`
	Redis    []map[string]string `yaml:"redis"`
	Pgs      []map[string]string `yaml:"pgs"`
}

type Config struct {
	Prometheus map[string]string `yaml:"prometheus"`
	CheckList  checklist         `yaml:"checklist"`
	Controller map[string]string `yaml:"controller"`
}

func InitConfig() *Config {
	if !checkFile("./xunjian.yaml") {
		log.Fatal("no config file")
	}
	if f, err := os.Open("./xunjian.yaml"); err != nil {
		log.Fatal(err)
	} else {
		yaml.NewDecoder(f).Decode(&Conf)
		return &Conf
	}
	return &Conf
}

func UrlEncodeCode(str string) string {
	return url.QueryEscape(str)
}

func UrlDecode(str string) string {
	res, err := url.QueryUnescape(str)
	if err != nil {
		return ""
	} else {
		return res
	}
}

func checkFile(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 检查IP是否合法
func CheckIP(ip string) bool {
	//if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", value); !m {
	//	return false
	//}
	re, _ := regexp.Compile(`^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}$`)
	matchd := re.MatchString(ip)
	if matchd {
		return true
	} else {
		return false
	}
}

// 检查二个字符串是否匹配
func CompareString(content, pattern string) bool {
	//fmt.Println(content, pattern)
	re, _ := regexp.Compile(pattern)
	matchd := re.MatchString(content)
	if matchd {
		return true
	} else {
		return false
	}
}
