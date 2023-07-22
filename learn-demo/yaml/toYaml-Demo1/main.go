package main

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ContainerPorts struct {
	Name          string `yaml:"name"`
	ContainerPort string `yaml:"containerPort"`
}

type Dplables struct {
	Labels map[string]string `yaml:"labels"`
}

type Dptoleration struct {
	Key      string `yaml:"key"`
	Operator string `yaml:"operator"`
	Effect   string `yaml:"effect"`
	Test     string `yaml:"test,omitempty"`
}

type Container struct {
	Image          string           `yaml:"image"`
	Name           string           `yaml:"name"`
	Ports          []ContainerPorts `yaml:"ports"`
	ReadinessProbe interface{}      `yaml:"readinessProbe,omitempty"`
}

type DpSpec struct {
	NodeName   string         `yaml:"nodeName"`
	Toleration []Dptoleration `yaml:"tolerations"`
	Cotainers  []Container    `yaml:"containers"`
}

type Tpmetadata struct {
	Metadata Dplables `yaml:"metadata"`
	Spec     DpSpec   `yaml:"spec"`
}

type Deployment struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Spec       struct {
		Template Tpmetadata `yaml:"template"`
	} `yaml:"spec"`
}

func ReadYamlConfig(path string) {
	if f, err := os.Open(path); err != nil {
		log.Fatal(err)
	} else {
		var dp Deployment
		yaml.NewDecoder(f).Decode(&dp)
		dp.Spec.Template.Spec.Cotainers = append(dp.Spec.Template.Spec.Cotainers, Container{
			Image: "test",
			Name:  "test",
			Ports: []ContainerPorts{
				ContainerPorts{
					Name:          "test-1",
					ContainerPort: "9999",
				},
				ContainerPorts{
					Name:          "test-2",
					ContainerPort: "8888",
				},
			},
		})
		output, err := os.OpenFile("./deployment-backup", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 077)
		if err != nil {
			log.Fatal(err)
		}
		yaml.NewEncoder(output).Encode(&dp)
		defer output.Close()
	}
}

func main() {
	ReadYamlConfig("./deployment.yaml")
}
