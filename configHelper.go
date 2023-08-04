package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getFromConfig(key string) (string, error) {
	data, err := os.ReadFile("config.txt")
	if err != nil {
		return "", fmt.Errorf("Could not open config.txt: %v", err)
	}

	config := strings.Split(string(data), "\n")
	for _, line := range config {
		parts := strings.Split(line, "=")
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("Could not find key '%s' in config.txt", key)
}

func getIntervalFromConfig() int {
	intervalStr, err := getFromConfig("interval")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return -1
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse interval from config.txt: %v\n", err)
		return -1
	}

	return interval
}

func getPortFromConfig() int {
	portStr, err := getFromConfig("port")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return -1
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse port from config.txt: %v\n", err)
		return -1
	}

	return port
}

func getSentencesFromConfig() []string {
	sentencesStr, err := getFromConfig("sentences")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}

	// Use regex to match parentheses and split by commas
	r := regexp.MustCompile(`\((.*?)\)`)
	matches := r.FindStringSubmatch(sentencesStr)
	if len(matches) < 2 {
		fmt.Fprintf(os.Stderr, "Could not parse sentences from config.txt\n")
		return nil
	}
	sentences := strings.Split(matches[1], ",")

	// Trim spaces
	for i, sentence := range sentences {
		sentences[i] = "$" + strings.TrimSpace(sentence)
	}

	return sentences
}
