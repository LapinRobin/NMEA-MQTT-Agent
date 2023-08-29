# Configurable NMEA parser and MQTT publisher written in Go

## Overview

This is a configurable NMEA parser and MQTT publisher written in Go. 

## Configuration

There are two config files: `udp_config.txt` and `mqtt_config.txt`. The `udp_config.txt` file contains the port to listen to on UDP, the interval of sending signals to MQTT, the list of NMEA sentences to parse and their structure in JSON format. The `mqtt_config.txt` file contains the MQTT broker address, port, clientID, username, and password.

## Compatibility

This program is built for Windows 10 and above only and has not been tested on other platforms.

## Usage

With the config files correctly configured, use the command `go build` to build the executable. Then run the executable either by clicking on it or by using the command `.\main.exe` in the Command Prompt or PowerShell.

To directly compile and run the code, use the command `go run .` .

## Caveats

To be able to stop the mosquitto service when the program terminates, run the program as administrator. 