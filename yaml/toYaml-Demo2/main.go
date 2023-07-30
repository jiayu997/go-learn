package main

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type FileSdconfig struct {
	Files []string `yaml:"files"`
}

type ScrapeConfigs struct {
	JobName       string            `yaml:"job_name"`
	StaticConfigs interface{}       `yaml:"static_configs"`
	BasicAuth     map[string]string `yaml:"basic_auth"`
	FileSdconfigs []*FileSdconfig   `yaml:"file_sd_configs"`
}

type PrometheusConfig struct {
	Global        interface{}     `yaml:"global"`
	Alerting      interface{}     `yaml:"alerting"`
	RuleFiles     []string        `yaml:"rule_files"`
	ScrapeConfigs []ScrapeConfigs `yaml:"scrape_configs"`
}

func newScrapFileConfig(file ...string) *FileSdconfig {
	return &FileSdconfig{file}

}

func newScrapConfig(job, username, password string) *ScrapeConfigs {
	paths := []string{
		fmt.Sprint("sd/file/%s/*.yaml", job),
		fmt.Sprintf("sd/file/%s/*.json", job),
	}
	return &ScrapeConfigs{
		JobName: "test",
		BasicAuth: map[string]string{
			username: username,
			password: password,
		},
		FileSdconfigs: []*FileSdconfig{newScrapFileConfig(paths...)},
	}
}

func main() {
	f, err := os.Open("prometheus.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var config PrometheusConfig

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", config)

	output, err := os.Create("prometheus-3.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	scrap := newScrapConfig("test", "test", "test")
	config.ScrapeConfigs = append(config.ScrapeConfigs, *scrap)
	encoder := yaml.NewEncoder(output)
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal(err)
	}
}
