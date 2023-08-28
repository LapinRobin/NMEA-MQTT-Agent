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
	case "$INLWY":
		return parseINLWY(line)
	case "$INMWV":
		return parseINMWV(line)
	case "$INMTW":
		return parseINMTW(line)
	case "$INMTA":
		return parseINMTA(line)
	case "$INRSA":
		return parseINRSA(line)
	case "$INMMB":
		return parseINMMB(line)
	case "$INVPW":
		return parseINVPW(line)
	case "$INHVD":
		return parseINHVD(line)
	case "$INVHW":
		return parseINVHW(line)
	case "$INMWD":
		return parseINMWD(line)
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
	data := make(map[string]string)

	// Remove the checksum
	splitLine := strings.Split(line, "*")

	// If there's anything other than 2 parts, the line was malformed
	if len(splitLine) != 2 {
		return data, false
	}

	// Split the remaining line on commas
	fields := strings.Split(splitLine[0], ",")
	for i := 1; i+3 < len(fields); i += 4 {
		status := fields[i]
		value := fields[i+1]
		// unit := fields[i+2]  // Uncomment if you need to use the unit
		label := fields[i+3]

		if i == 1 && status == "G" && value != "" {
			data[label] = value
		}

		if status == "A" && value != "" {
			data[label] = value
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

func parseINLWY(line string) (map[string]string, bool) {
	data := map[string]string{
		"angle": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["angle"] = fields[2]

	return data, true
}

func parseINMWV(line string) (map[string]string, bool) {
	data := map[string]string{
		"angle": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["angle"] = fields[1]

	return data, true
}

func parseINMTW(line string) (map[string]string, bool) {
	data := map[string]string{
		"temperature": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["temperature"] = fields[1]

	return data, true
}

func parseINMTA(line string) (map[string]string, bool) {
	data := map[string]string{
		"temperature": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["temperature"] = fields[1]

	return data, true
}

func parseINRSA(line string) (map[string]string, bool) {
	data := map[string]string{
		"angle": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["angle"] = fields[1]

	return data, true
}

func parseINMMB(line string) (map[string]string, bool) {
	data := map[string]string{
		"pressure": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["pressure"] = fields[3]

	return data, true
}

// $INVPW,-0.0,N,-0.0,M*55
func parseINVPW(line string) (map[string]string, bool) {
	data := map[string]string{
		"N": "0",
		"M": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["N"] = fields[1]
	data["M"] = fields[3]

	return data, true
}

func parseINHVD(line string) (map[string]string, bool) {
	data := map[string]string{
		"hvd": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["hvd"] = fields[1]

	return data, true
}

func parseINVHW(line string) (map[string]string, bool) {
	data := map[string]string{
		"angle":  "0",
		"speedN": "0",
		"speedK": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}

	fields := strings.Split(splitLine[0], ",")
	data["angle"] = fields[1]
	data["speedN"] = fields[5]
	data["speedK"] = fields[7]

	return data, true
}

func parseINMWD(line string) (map[string]string, bool) {
	data := map[string]string{
		"T": "0",
		"N": "0",
		"M": "0",
	}

	splitLine := strings.Split(line, "*")
	if len(splitLine) != 2 {
		return data, false
	}
	fields := strings.Split(splitLine[0], ",")
	data["T"] = fields[1]
	data["N"] = fields[5]
	data["M"] = fields[7]

	return data, true
}
