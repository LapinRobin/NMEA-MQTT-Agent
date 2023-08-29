package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getFromConfig(key string) (string, error) {
	data, err := os.ReadFile("udp_config.txt")
	if err != nil {
		return "", fmt.Errorf("Could not open udp_config.txt: %v", err)
	}
	// only get the first three lines
	config := strings.Split(string(data), "\n")[:3]
	for _, line := range config {
		parts := strings.Split(line, "=")
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("Could not find key '%s' in udp_config.txt", key)
}

func GetIntervalFromConfig() int {
	intervalStr, err := getFromConfig("interval")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return -1
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse interval from udp_config.txt: %v\n", err)
		return -1
	}

	return interval
}

func GetPortFromConfig() int {
	portStr, err := getFromConfig("port")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return -1
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse port from udp_config.txt: %v\n", err)
		return -1
	}

	return port
}

func GetSentencesFromConfig() []string {
	sentencesStr, err := getFromConfig("sentences")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}

	// Use regex to match parentheses and split by commas
	r := regexp.MustCompile(`\((.*?)\)`)
	matches := r.FindStringSubmatch(sentencesStr)
	if len(matches) < 2 {
		fmt.Fprintf(os.Stderr, "Could not parse sentences from udp_config.txt\n")
		return nil
	}
	sentences := strings.Split(matches[1], ",")

	// Trim spaces
	for i, sentence := range sentences {
		sentences[i] = "$" + strings.TrimSpace(sentence)
	}

	return sentences
}

func GetMapFromConfig() (map[string]map[string]interface{}, error) {
	file, err := os.Open("udp_config.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialize variables
	var sb strings.Builder
	collectJSON := false

	// Read each line from the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "map:" {
			collectJSON = true
			continue
		}
		if collectJSON {
			sb.WriteString(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Decode JSON content into a Go map
	var data map[string]map[string]interface{}
	err = json.Unmarshal([]byte(sb.String()), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
