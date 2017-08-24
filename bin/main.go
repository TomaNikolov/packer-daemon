package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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

	daemonConfig, err := readConfig(*configPath)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	q := sqs.New(sess, &aws.Config{
		Region:      aws.String(daemonConfig.AwsRegion),
		Credentials: credentials.NewStaticCredentials(daemonConfig.AwsPublicKey, daemonConfig.AwsPriveteKey, ""),
	})

	worker.Start(q, daemonConfig)
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
