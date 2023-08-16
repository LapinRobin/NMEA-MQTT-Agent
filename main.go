package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func isZeroValue(s string) bool {
	// Try to parse as float first (which also handles integers)
	if val, err := strconv.ParseFloat(s, 64); err == nil && val == 0 {
		return true
	}
	return false
}

func main() {

	// Set up signal catching
	signals := make(chan os.Signal, 1)
	// Notify signals (SIGINT = Ctrl+C, SIGTERM = Termination request)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Block until we receive a signal
		sig := <-signals
		fmt.Printf("\nReceived signal: %s. Stopping Mosquitto service...\n", sig)
		// wait for mosquitto to stop
		time.Sleep(500 * time.Millisecond)
		if IsMosquittoRunning() {
			fmt.Println("Mosquitto service is still running, shutting down...")
			StopMosquitto()
		} else {
			fmt.Println("Mosquitto service is off, terminating program...")
		}

		os.Exit(0)
	}()

	// startMosquitto()
	print("Starting Mosquitto service...\n")
	StartMosquitto()

	print("Fetching MQTT config from mqtt_config.txt...\n")
	config := GetMqttConfig()

	print("Creating and starting MQTT client...\n")
	mqttClient := CreateAndStartClient(config)

	print("Fetching MQTT topic from mqtt_config.txt...\n")
	topic := GetMqttTopic()

	interval := getIntervalFromConfig()
	print("Fetching interval from config.txt...\n")

	// print(interval)
	if interval == -1 {
		fmt.Fprintf(os.Stderr,
			"Could not read config.txt, defaulting to 10 seconds\n")
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

	fmt.Println("Listening on UDP port", port)
	fmt.Println("Sending data every", interval, "milliseconds")

	// print broker address and topic

	fmt.Println("Topic:", topic)
	fmt.Println("Broker:", config.Broker)

	// store parsed data for each sentence
	parsedData := make(map[string]map[string]string)

	// parsing
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
						for key, value := range data {
							// Update only if the value is not representing a zero
							if !isZeroValue(value) {
								parsedData[sentence][key] = value
							}
						}
					}
				}
			}
		}

		// If it's time to send data and if any data to send
		if time.Now().After(nextWriteTime) && len(parsedData) > 0 {
			/* jsonData, err := json.Marshal(parsedData)
			if err != nil {
				log.Fatalf("Failed to marshal JSON: %v", err)
			}

			print("publishing data...\n")

			// Publish JSON data to an MQTT topic
			token := mqttClient.Publish(topic, 0, false, jsonData)
			token.Wait()

			nextWriteTime = nextWriteTime.Add(time.Duration(interval) * time.Millisecond)
			parsedData = make(map[string]map[string]string) */

			// 1. Transform string values to numbers.
			transformedData := make(map[string]interface{})

			var datetime string
			for key, innerMap := range parsedData {
				transformedData[key] = make(map[string]interface{})

				for innerKey, value := range innerMap {
					switch innerKey {
					case "northSouth":
						if value == "N" {
							transformedData[key].(map[string]interface{})[innerKey] = 1
						} else if value == "S" {
							transformedData[key].(map[string]interface{})[innerKey] = 0
						}
					case "eastWest":
						if value == "E" {
							transformedData[key].(map[string]interface{})[innerKey] = 1
						} else if value == "W" {
							transformedData[key].(map[string]interface{})[innerKey] = 0
						}
					case "date", "time":
						// We'll handle these below after processing all the other fields.
					default:
						if numVal, err := strconv.ParseFloat(value, 64); err == nil {
							transformedData[key].(map[string]interface{})[innerKey] = numVal
						} else {
							transformedData[key].(map[string]interface{})[innerKey] = value
						}
					}
				}

				// Combine date and time to form Unix datetime.
				if date, ok := innerMap["date"]; ok {
					if t, ok := innerMap["time"]; ok {
						// Split the time string to separate seconds and milliseconds.
						timeParts := strings.Split(t, ".")
						// Format seconds and milliseconds.
						seconds := timeParts[0]
						// Construct the datetime string without timezone since the provided data does not have a timezone.
						dtStr := fmt.Sprintf("20%s-%s-%sT%s:%s:%s+00:00", date[4:6], date[2:4], date[0:2], seconds[0:2], seconds[2:4], seconds[4:6])
						// Validate the datetime string using your format.
						if _, err := time.Parse("2006-01-02T15:04:05Z07:00", dtStr); err == nil {
							datetime = dtStr
						} else {
							log.Printf("Error parsing datetime: %v", err)
						}
					}
				}

			}

			transformedData["datetime"] = datetime

			jsonData, err := json.Marshal(transformedData)
			if err != nil {
				log.Fatalf("Failed to marshal JSON: %v", err)
			}

			print("publishing data...\n")

			// Publish JSON data to an MQTT topic
			token := mqttClient.Publish(topic, 0, false, jsonData)
			token.Wait()

			nextWriteTime = nextWriteTime.Add(time.Duration(interval) * time.Millisecond)
			parsedData = make(map[string]map[string]string)

		}
	}

}
