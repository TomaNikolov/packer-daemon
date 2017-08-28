package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/tomanikolov/packer-daemon/types"
	"github.com/tomanikolov/packer-daemon/utils"
	"github.com/tomanikolov/packer-daemon/worker"
)

func main() {

	userDir, err := utils.GetUserDir()
	if err != nil {
		log.Fatal(err)
	}

	configPath := flag.String("config", path.Join(userDir, ".packer-daemon-config.json"), "Path to config file")
	flag.Parse()
	fmt.Println("path: " + *configPath)

	config, err := readConfig(*configPath)
	worker.Start(config)
}

func readConfig(configPath string) (types.Config, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	config := types.Config{}
	err = json.Unmarshal(file, &config)

	return config, err
}
