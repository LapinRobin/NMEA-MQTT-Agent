package main

import "strings"

// skeleton function
func parseSentence(sentenceType string, line string) (map[string]string, bool) {
	switch sentenceType {
	case "$GPRMC":
		return parseGPRMC(line)
	case "$INXDR":
		return parseINXDR(line)
	case "$INDPT":
		return parseINDPT(line)
	case "$INHDT":
		return parseINHDT(line)
	default:
		// Unsupported sentence typeinterface{}{}, false
		return map[string]string{}, false
	}

}

// helper function to parse $GPRMC sentences
func parseGPRMC(line string) (map[string]string, bool) {
	isValidData := true
	linePtr := line
	token := ""
	field := 0
	data := make(map[string]string)

	for isValidData && field <= 9 {
		i := 0
		token = ""

		for i < len(linePtr) && linePtr[i] != ',' {
			token += string(linePtr[i])
			i++
		}

		switch field {
		case 1:
			if len(token) == 0 {
				data["time"] = "0"
			} else {
				data["time"] = token
			}
		case 2:
			if token != "A" {
				isValidData = false
			}
		case 3:
			if len(token) == 0 {
				data["latitude"] = "0"
			} else {
				data["latitude"] = token
			}
		case 4:
			data["northSouth"] = token
		case 5:
			if len(token) == 0 {
				data["longitude"] = "0"
			} else {
				data["longitude"] = token
			}
		case 6:
			data["eastWest"] = token
		case 7:
			if len(token) == 0 {
				data["speed"] = "0"
			} else {
				data["speed"] = token
			}
		case 8:
			if len(token) == 0 {
				data["angle"] = "0"
			} else {
				data["angle"] = token
			}
		case 9:
			if len(token) == 0 {
				data["date"] = "0"
			} else {
				data["date"] = token
			}
		}

		if i < len(linePtr) && linePtr[i] == ',' {
			linePtr = linePtr[i+1:] // Skip the comma
		} else {
			linePtr = linePtr[i:]
		}
		field++
	}
	return data, isValidData
}

// helper function to parse $INXDR sentences
func parseINXDR(line string) (map[string]string, bool) {
	// Initialize map
	data := map[string]string{
		"Trim":    "0",
		"FoilA":   "0",
		"CANT":    "0",
		"D0":      "0",
		"D0lee":   "0",
		"Runner":  "0",
		"J2":      "0",
		"J3":      "0",
		"RSA":     "0",
		"Baro":    "0",
		"FOILMIN": "0",
		"Leeway":  "0",
		"MastRot": "0",
	}

	// Remove the checksum
	splitLine := strings.Split(line, "*")

	// If there's anything other than 2 parts, the line was malformed
	if len(splitLine) != 2 {
		return data, false
	}

	// Split the remaining line on commas
	fields := strings.Split(splitLine[0], ",")

	for i := 1; i < len(fields)-1; i += 4 {
		if fields[i] == "A" {
			if _, ok := data[fields[i+3]]; ok {
				data[fields[i+3]] = fields[i+1]
			}
		}
	}

	return data, true
}

// helper function to parse $INDPT sentences
func parseINDPT(line string) (map[string]string, bool) {
	// Initialize map
	data := map[string]string{
		"depth": "0",
	}

	// Remove the checksum
	splitLine := strings.Split(line, "*")

	// If there's anything other than 2 parts, the line was malformed
	if len(splitLine) != 2 {
		return data, false
	}

	// Split the remaining line on commas
	fields := strings.Split(splitLine[0], ",")

	data["depth"] = fields[1]

	return data, true
}

// helper function to parse $INHDT sentences
func parseINHDT(line string) (map[string]string, bool) {
	data := map[string]string{
		"heading": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["heading"] = fields[1]

	return data, true
}
