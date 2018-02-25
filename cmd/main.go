package main

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"code.cloudfoundry.org/lager"
	"os"
	"gopkg.in/yaml.v2"
)

// Element1 and Element2 have to start with uppercase
// it means they are "exported" in golang, otherwise
// marshaller will omit them
type ConfigMapDataElements struct {
	Element1 string `yaml:"element1"`
	Element2 string `yaml:"element2"`
}

type ConfigMapData struct {
	ListOfElements []ConfigMapDataElements `yaml:"listOfElements,omitempty"`
}

func main(){
	logger := lager.NewLogger("configmap-example-app")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	configMaps := clientSet.CoreV1().ConfigMaps("default")
	cm, err := configMaps.Get("test-cm", meta_v1.GetOptions{})
	if err != nil {
		logger.Error("Failed to get cm", err)
		panic("I'm so panicked")
	}
	data := ConfigMapData{}
	yaml.Unmarshal([]byte(cm.Data["entry"]), &data)
	logger.Info("marshalled data:", lager.Data{"data": cm.Data["entry"]})
	logger.Info("unmarshalled data", lager.Data{"data": data,})
	// let's add another key data
	newElement := ConfigMapDataElements{
		Element1: "newly added Element1",
		Element2: "also quite new element2",
	}
	data.ListOfElements = append(data.ListOfElements, newElement)
	marshalled, _ := yaml.Marshal(&data)
	cm.Data["entry"] = string(marshalled[:])
	configMaps.Update(cm)

}
