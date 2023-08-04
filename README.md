# Configurable NMEA parser and MQTT publisher written in Go

## Overview

This is a configurable NMEA parser and MQTT publisher written in Go. It is designed to be used with the [OpenPlotter Signal K Node Server]

## Configuration

There are two config files: `config.txt` and `mqtt_config.txt`. The `config.txt` file contains the port to listen to on UDP, the interval of sending signals to MQTT, and the list of NMEA sentences to parse. The `mqtt_config.txt` file contains the MQTT broker address, port, clientID, username, and password.

## Compatibility

This program is built for Windows 10 and above only and has not been tested on other platforms.

## Usage

With the config files correctly configured, run the command `go build` to build the executable. Then run the executable either by clicking on it or by running the command `.\main.exe` in the terminal.