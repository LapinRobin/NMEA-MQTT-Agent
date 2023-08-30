package main

import (
	"fmt"
	"strings"
)

// skeleton function
func parseSentence(sentenceType string, line string, sentenceMap map[string]map[string]interface{}) (map[string]string, bool) {

	if sentenceType == "$INXDR" {
		return parseINXDR(line, sentenceMap)
	} else {
		return parseCommon(line, sentenceMap)
	}
}

func parseCommon(line string, sentenceMap map[string]map[string]interface{}) (map[string]string, bool) {
	data := make(map[string]string)

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	sentenceType := fields[0]

	if mapping, exists := sentenceMap[sentenceType]; exists {
		for pos, name := range mapping {
			posInt := 0
			fmt.Sscan(pos, &posInt) // Convert the position from string to integer
			if posInt > 0 && posInt < len(fields) {
				data[name.(string)] = fields[posInt]
			}
		}
		return data, true
	}

	return data, false
}

// Parse $INXDR sentences
func parseINXDR(line string, sentenceMap map[string]map[string]interface{}) (map[string]string, bool) {
	data := make(map[string]string)

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	sentenceType := fields[0] // Should be "$INXDR"

	// Find the last field name, such as "FOILMIN"
	lastName := fields[len(fields)-1]

	// Check if this field is in the sentenceMap
	if mapping, exists := sentenceMap[sentenceType][lastName]; exists {
		mappingTyped := mapping.(map[string]interface{})
		for pos, name := range mappingTyped {
			posInt := 0
			fmt.Sscan(pos, &posInt) // Convert the position from string to integer
			if posInt > 0 && posInt < len(fields) {
				data[name.(string)] = fields[posInt]
			}
		}
		return data, true
	}

	return data, false
}
