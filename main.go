package main

import (
	"encoding/json"
	"fmt" // Import the mqttConfig package
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	config := GetMqttConfig()
	mqttClient := CreateAndStartClient(config)
	topic := GetMqttTopic()
	interval := getIntervalFromConfig()
	print("Fetching interval from config.txt\n")

	// print(interval)
	if interval == -1 {
		fmt.Fprintf(os.Stderr, "Could not read config.txt, defaulting to 10 seconds\n")
		interval = 10000
	}

	port := getPortFromConfig()
	print("Fetching port from config.txt\n")

	sentences := getSentencesFromConfig()
	print("Fetching sentences from config.txt\n")
	print("Sentences to parse from: \n")
	for _, sentence := range sentences {
		fmt.Println(sentence)
	}

	buffSize := 256
	var conn *net.UDPConn
	var file *os.File
	nextWriteTime := time.Now()

	addr := net.UDPAddr{
		Port: port,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Could not start server, error:", err)
		return
	}
	defer conn.Close()

	file, err = os.Create("GPS_data.txt")
	if err != nil {
		fmt.Println("Could not create file, error:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("UT;POS;BSP;SOG;COG;TWS;TWD;TWA;PRES\n")
	if err != nil {
		fmt.Println("Could not write to file, error:", err)
		return
	}

	fmt.Println("Listening on UDP port", port)
	fmt.Println("Sending data every", interval, "milliseconds")

	// store parsed data for each sentence
	parsedData := make(map[string]map[string]string)

	// GPRMC parsing
	for {
		buffer := make([]byte, buffSize)
		_, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalf("ReadFromUDP failed: %v", err)
		}

		lines := strings.Split(string(buffer), "\r\n")

		for _, line := range lines {
			for _, sentence := range sentences {
				if strings.HasPrefix(line, sentence) {
					data, isValidData := parseSentence(sentence, line)

					if isValidData {
						if parsedData[sentence] == nil {
							parsedData[sentence] = make(map[string]string)
						}
						// store parsed data for each sentence type
						// basically a map of maps
						// stores the latest data for each sentence type
						parsedData[sentence] = data
					}
				}
			}
		}

		// If it's time to send data and there's any data to send
		if time.Now().After(nextWriteTime) && len(parsedData) > 0 {
			jsonData, err := json.Marshal(parsedData)
			if err != nil {
				log.Fatalf("Failed to marshal JSON: %v", err)
			}

			print("publishing data...\n")
			// print JSON data to console
			// fmt.Println(string(jsonData))
			// Publish JSON data to an MQTT topic

			token := mqttClient.Publish(topic, 0, false, jsonData)
			token.Wait()

			nextWriteTime = nextWriteTime.Add(time.Duration(interval) * time.Millisecond)
			parsedData = make(map[string]map[string]string)

		}
	}

}
