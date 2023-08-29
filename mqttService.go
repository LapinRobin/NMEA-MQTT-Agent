package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var isConnected = false
var offlineMessages []string
var mutex sync.Mutex

type MqttConfig struct {
	Broker   string
	ClientID string
	Password string
	Username string
}

func IsMosquittoRunning() bool {
	// Use tasklist to check for mosquitto.exe process
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq mosquitto.exe")
	output, err := cmd.Output() // This captures the standard output of the command

	if err != nil {
		log.Fatalf("Failed to fetch task list: %s", err)
		return false
	}

	// If the output contains 'mosquitto.exe', then Mosquitto is running
	return strings.Contains(string(output), "mosquitto.exe")
}

func StartMosquitto() *exec.Cmd {
	if IsMosquittoRunning() {
		fmt.Println("mosquitto.exe is already running.")
		return nil
	}

	cmd := exec.Command(`C:\Program Files\mosquitto\mosquitto.exe`)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Mosquitto: %s", err)
	}

	return cmd
}

func StopMosquitto() {

	cmd := exec.Command(`taskkill`, `/IM`, `mosquitto.exe`, `/F`)
	err := cmd.Run()
	// wait for 1 second
	cmd.Wait()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 128 {
			// Process killed
			log.Printf("Warning: Mosquitto process not found.")
		} else {
			// Other errors
			log.Printf("Error while attempting to terminate Mosquitto: %s", err)
		}
	} else {
		log.Println("Successfully terminated Mosquitto.")
	}
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

func CreateAndStartClient(config MqttConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	opts.SetPassword(config.Password)
	opts.SetUsername(config.Username)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error() // Return the error instead of logging it
	}

	return client, nil // return client and nil error on successful connection
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

func isBrokerAvailable(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", host+":"+port, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true

}

func CheckBrokerConnectionRegularly(brokerURL string, client mqtt.Client) {
	// Extract hostname and port from the broker URL
	urlParts := strings.Split(brokerURL, "://")
	hostParts := strings.Split(urlParts[1], ":")
	host := hostParts[0]
	port := hostParts[1]

	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds; adjust as needed
	// first check without delay
	if !isBrokerAvailable(host, port) {
		// Connection is lost
		isConnected = false
		fmt.Println("Broker is not available")
	} else { // If previously it was not connected, but now it is
		isConnected = true
		fmt.Println("Broker is now available")
	}
	for range ticker.C {
		if !isBrokerAvailable(host, port) {
			// Connection is lost
			isConnected = false
			fmt.Println("Broker is not available")
		} else if !isConnected { // If previously it was not connected, but now it is
			isConnected = true
			fmt.Println("Broker is now available")
		}
	}
}

func PublishMessage(client mqtt.Client, topic string, qos byte, retained bool, payload string) {

	if isConnected {
		// If connected, first send all buffered messages
		fmt.Println("Client is connected.")
		mutex.Lock()
		for _, msg := range offlineMessages {
			fmt.Println("Publishing buffer messages...")
			token := client.Publish(topic, qos, retained, msg)
			token.Wait()
		}
		offlineMessages = []string{} // Clear the buffer
		mutex.Unlock()

		// Publish current message
		token := client.Publish(topic, qos, retained, payload)
		token.Wait()
	} else {
		fmt.Println("Client is disconnected. Buffering message: ", payload)
		// If not connected, store the message in the buffer
		mutex.Lock()
		offlineMessages = append(offlineMessages, payload)
		mutex.Unlock()
	}
}
