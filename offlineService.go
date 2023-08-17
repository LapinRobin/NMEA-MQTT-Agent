package main

import "fmt"

var offlineMessages = []string{}
var isConnected = false

func onConnectionLost(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err.Error())
	isConnected = false
}

func onConnect(client MQTT.Client) {
	fmt.Println("Connected")
	isConnected = true

	// Republish buffered messages
	for _, msg := range offlineMessages {
		publish(client, "your/topic", msg)
	}
	offlineMessages = [] // Clear the buffer after republishing
}

func publish(client MQTT.Client, topic, message string) {
	if isConnected {
		token := client.Publish(topic, 0, false, message)
		token.Wait()
	} else {
		fmt.Println("Offline, buffering message:", message)
		offlineMessages = append(offlineMessages, message)
		// Optionally, you can write this to a file for persistence
	}
}


