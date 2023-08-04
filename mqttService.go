package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Broker   string
	ClientID string
	Password string
	Username string
}

func GetMqttConfig() MqttConfig {
	file, err := os.Open("mqtt_config.txt") // or path to your file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := MqttConfig{}
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "=")
		if len(split) == 2 {
			switch split[0] {
			case "broker":
				config.Broker = split[1]
			case "clientID":
				config.ClientID = split[1]
			case "password":
				config.Password = split[1]
			case "username":
				config.Username = split[1]
			default:
				// skip
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return config
}

func CreateAndStartClient(config MqttConfig) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	opts.SetPassword(config.Password)
	opts.SetUsername(config.Username)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return client
}

func GetMqttTopic() string {
	file, err := os.Open("mqtt_config.txt") // or path to your file
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	topic := ""
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, "=")
		if len(split) == 2 && split[0] == "topic" {
			topic = split[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return topic
}
